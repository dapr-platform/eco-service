package service

import (
	"context"
	"eco-service/client"
	"eco-service/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
)
var needRefreshContinuousAggregateMap = map[string]string{
	"f_eco_gateway_1d":  "day",
	"f_eco_gateway_1m":  "month",
	"f_eco_gateway_1y":  "year",
	"f_eco_floor_1d":    "day",
	"f_eco_floor_1m":    "month", 
	"f_eco_floor_1y":    "year",
	"f_eco_building_1d": "day",
	"f_eco_building_1m": "month",
	"f_eco_building_1y": "year",
}

func init() {
	// Start goroutine to collect stats every hour at 5 minutes past
	go func() {
		for {
			now := time.Now()
			// Calculate next run time (next hour at 5 minutes past)
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 5, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			if err := CollectGatewayHourlyStats(); err != nil {
				common.Logger.Errorf("Failed to collect gateway hourly stats: %v", err)
			}
		}
	}()
}

// CollectGatewayHourlyStats collects gateway hourly stats at 5 minutes past each hour
// Gets data for previous hours and stores in database
// hoursAgo specifies how many previous hours to collect, defaults to 4 if not specified
func CollectGatewayHourlyStats(hoursAgo ...int) error {
	common.Logger.Infof("Collecting gateway hourly stats")
	// Get hours to collect, default to 4 if not specified
	hours := 4
	if len(hoursAgo) > 0 && hoursAgo[0] > 0 {
		hours = hoursAgo[0]
	}

	// Get current time
	now := time.Now()

	// Get all gateways
	gateways, err := GetAllEcgateways()
	if err != nil {
		return err
	}

	// Group gateways by project code
	projectGateways := make(map[string][]model.Ecgateway)
	for _, gateway := range gateways {
		projectCode := gateway.ProjectCode
		if len(projectCode) == 0 {
			projectCode, err = client.GetBoxProjectCode(gateway.MacAddr)
			if err != nil {
				common.Logger.Errorf("Failed to get project code for gateway %s: %v", gateway.ID, err)
				return errors.Wrap(err, "Failed to get project code for gateway")
			}
			gateway.ProjectCode = projectCode
			err = common.DbUpsert[model.Ecgateway](context.Background(), common.GetDaprClient(), gateway, model.EcgatewayTableInfo.Name, model.Ecgateway_FIELD_NAME_id)
			if err != nil {
				common.Logger.Errorf("Failed to update project code for gateway %s: %v", gateway.ID, err)
				return errors.Wrap(err, "Failed to update project code for gateway")
			}
		}

		projectGateways[projectCode] = append(projectGateways[projectCode], gateway)
	}

	// For each project, collect stats for all gateways
	for projectCode, gateways := range projectGateways {
		// Get mac addresses for all gateways in this project
		macAddrs := make([]string, len(gateways))
		for i, gateway := range gateways {
			macAddrs[i] = gateway.MacAddr
		}

		// Collect stats for previous hours
		for i := 1; i <= hours; i++ {
			// Calculate hour timestamp
			hourTime := now.Add(time.Duration(-i) * time.Hour)
			hourTime = time.Date(hourTime.Year(), hourTime.Month(), hourTime.Day(),
				hourTime.Hour(), 0, 0, 0, hourTime.Location())

			// Call GET_BOXES_HOUR_STATS API with all mac addresses
			reqBody := map[string]string{
				"projectCode": projectCode,
				"mac":         strings.Join(macAddrs, ","),
				"year":        hourTime.Format("2006"),
				"month":       hourTime.Format("01"),
				"day":         hourTime.Format("02"),
				"hour":        hourTime.Format("15"),
			}

			respBytes, err := client.GetBoxesHourStats(reqBody)
			if err != nil {
				return err
			}

			var resp struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Data    map[string][]struct {
					Addr        int     `json:"addr"`
					Electricity float64 `json:"electricity"`
				} `json:"data"`
			}

			if err := json.Unmarshal(respBytes, &resp); err != nil {
				return err
			}

			if resp.Code != "0" {
				return fmt.Errorf("API error: %s", resp.Message)
			}

			// Process response for each gateway
			for _, gateway := range gateways {
				hourlyStats := model.Eco_gateway_1h{
					ID:               gateway.ID + "_" + hourTime.Format("2006010215"),
					Time:             common.LocalTime(hourTime),
					GatewayID:        gateway.ID,
					FloorID:          gateway.FloorID,
					BuildingID:       gateway.BuildingID,
					Type:             gateway.Type,
					PowerConsumption: 0,
				}

				// Get total electricity from addr 0
				if stats, ok := resp.Data[gateway.MacAddr]; ok {
					for _, stat := range stats {
						if stat.Addr == 0 {
							hourlyStats.PowerConsumption = stat.Electricity
							break
						}
					}
				}

				
			}
		}
		
	}

	// Refresh continuous aggregates based on type
	for tableName, refreshType := range needRefreshContinuousAggregateMap {
		startTime := time.Now()
		endTime := startTime
		
		switch refreshType {
		case "day":
			startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
			endTime = startTime.AddDate(0, 0, 1)
		case "month":
			startTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, startTime.Location())
			endTime = startTime.AddDate(0, 1, 0)
		case "year":
			startTime = time.Date(startTime.Year(), 1, 1, 0, 0, 0, 0, startTime.Location())
			endTime = startTime.AddDate(1, 0, 0)
		}

		if err := common.DbRefreshContinuousAggregate(context.Background(), common.GetDaprClient(), 
			tableName, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")); err != nil {
			return err
		}
	}

	return nil
}

func GetAllEcgateways() ([]model.Ecgateway, error) {
	datas, err := common.DbQuery[model.Ecgateway](context.Background(), common.GetDaprClient(), model.EcgatewayTableInfo.Name, "")
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func GetAllEcbuildings() ([]model.Ecbuilding, error) {
	datas, err := common.DbQuery[model.Ecbuilding](context.Background(), common.GetDaprClient(), model.EcbuildingTableInfo.Name, "")
	if err != nil {
		return nil, err
	}
	return datas, nil
}
