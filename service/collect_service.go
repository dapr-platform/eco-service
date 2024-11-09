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

			gateways, err := GetAllEcgateways()
			if err != nil {
				common.Logger.Errorf("Failed to get gateways: %v", err)
				continue
			}

			if err := collectGatewaysHours(time.Now(), 4, gateways); err != nil {
				common.Logger.Errorf("Failed to collect gateway hourly stats: %v", err)
			}

			if err := refreshContinuousAggregate(time.Now()); err != nil {
				common.Logger.Errorf("Failed to refresh continuous aggregates: %v", err)
			}
		}
	}()
}

// 手动收集指定日期的数据，开始时间结束时间格式为 2024-01-01
func ManuCollectGatewayHourlyStatsByDay(start, end string) error {
	if start == "" || end == "" {
		return errors.New("Start and end dates are required")
	}

	// Parse start and end dates
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return errors.Wrap(err, "Failed to parse start date")
	}

	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return errors.Wrap(err, "Failed to parse end date")
	}

	// Validate dates
	if endTime.Before(startTime) {
		return errors.New("End date must be after start date")
	}

	gateways, err := GetAllEcgateways()
	if err != nil {
		return errors.Wrap(err, "Failed to get gateways")
	}

	if len(gateways) == 0 {
		return errors.New("No gateways found")
	}

	// Iterate through each day
	for currentDate := startTime; !currentDate.After(endTime); currentDate = currentDate.AddDate(0, 0, 1) {
		if err := collectGatewaysFullDay(currentDate, gateways); err != nil {
			common.Logger.Errorf("Failed to collect stats for %s: %v",
				currentDate.Format("2006-01-02"), err)
			return err
		}

		if err := refreshContinuousAggregate(currentDate); err != nil {
			return err
		}
	}

	return nil
}

func collectGatewaysFullDay(collectTime time.Time, gateways []model.Ecgateway) error {
	// Group gateways by project code
	projectGateways := make(map[string][]model.Ecgateway)
	for _, gateway := range gateways {
		projectCode := gateway.ProjectCode
		if len(projectCode) == 0 {
			var err error
			projectCode, err = client.GetBoxProjectCode(gateway.MacAddr)
			if err != nil {
				return errors.Wrapf(err, "Failed to get project code for gateway %s", gateway.ID)
			}
			gateway.ProjectCode = projectCode
			if err := common.DbUpsert[model.Ecgateway](context.Background(), common.GetDaprClient(), gateway, model.EcgatewayTableInfo.Name, model.Ecgateway_FIELD_NAME_id); err != nil {
				return errors.Wrapf(err, "Failed to update project code for gateway %s", gateway.ID)
			}
		}
		projectGateways[projectCode] = append(projectGateways[projectCode], gateway)
	}

	// For each project, collect full day stats
	for projectCode, projectGateways := range projectGateways {
		// Process gateways in batches of 20
		for i := 0; i < len(projectGateways); i += 20 {
			end := i + 20
			if end > len(projectGateways) {
				end = len(projectGateways)
			}
			
			gatewayBatch := projectGateways[i:end]
			macAddrs := make([]string, len(gatewayBatch))
			for j, gateway := range gatewayBatch {
				macAddrs[j] = gateway.MacAddr
			}

			reqBody := map[string]string{
				"projectCode": projectCode,
				"mac":         strings.Join(macAddrs, ","),
				"year":        collectTime.Format("2006"),
				"month":       collectTime.Format("01"),
				"day":         collectTime.Format("02"),
			}

			respBytes, err := client.GetBoxesHourStats(reqBody)
			if err != nil {
				return errors.Wrap(err, "Failed to get box hour stats")
			}

			var resp struct {
				Code    string                 `json:"code"`
				Message string                 `json:"message"`
				Data    map[string]interface{} `json:"data"`
			}

			if err := json.Unmarshal(respBytes, &resp); err != nil {
				return errors.Wrap(err, "Failed to unmarshal response")
			}

			if resp.Code != "0" {
				return fmt.Errorf("API error: %s", resp.Message)
			}

			// Process response for each gateway and hour
			for _, gateway := range gatewayBatch {
				var hourlyStats []model.Eco_gateway_1h

				if gatewayData, ok := resp.Data[gateway.MacAddr].(map[string]interface{}); ok {
					for hour := 0; hour < 24; hour++ {
						hourStr := fmt.Sprintf("%02d", hour)
						if hourData, ok := gatewayData[hourStr].([]interface{}); ok {
							stats := processHourStats(hourData)
							hourTime := time.Date(collectTime.Year(), collectTime.Month(), collectTime.Day(), hour, 0, 0, 0, collectTime.Location())
							
							stat := model.Eco_gateway_1h{
								ID:               gateway.ID + "_" + hourTime.Format("2006010215"),
								Time:             common.LocalTime(hourTime),
								GatewayID:        gateway.ID,
								FloorID:          gateway.FloorID,
								BuildingID:       gateway.BuildingID,
								Type:             gateway.Type,
								PowerConsumption: getTotalElectricity(stats),
							}
							hourlyStats = append(hourlyStats, stat)
						}
					}
				}

				if len(hourlyStats) > 0 {
					if err := saveGatewayHourlyStats(hourlyStats); err != nil {
						return errors.Wrap(err, "Failed to save gateway hourly stats")
					}
				}
			}
		}
	}

	return nil
}

func collectGatewaysHours(collectTime time.Time, hoursAgo int, gateways []model.Ecgateway) error {
	if hoursAgo <= 0 {
		return errors.New("hoursAgo must be greater than 0")
	}

	// Group gateways by project code
	projectGateways := make(map[string][]model.Ecgateway)
	for _, gateway := range gateways {
		projectCode := gateway.ProjectCode
		if len(projectCode) == 0 {
			var err error
			projectCode, err = client.GetBoxProjectCode(gateway.MacAddr)
			if err != nil {
				return errors.Wrapf(err, "Failed to get project code for gateway %s", gateway.ID)
			}
			gateway.ProjectCode = projectCode
			if err := common.DbUpsert[model.Ecgateway](context.Background(), common.GetDaprClient(), gateway, model.EcgatewayTableInfo.Name, model.Ecgateway_FIELD_NAME_id); err != nil {
				return errors.Wrapf(err, "Failed to update project code for gateway %s", gateway.ID)
			}
		}
		projectGateways[projectCode] = append(projectGateways[projectCode], gateway)
	}

	// For each project, collect stats for specified hours
	for projectCode, projectGateways := range projectGateways {
		// Process gateways in batches of 20
		for i := 0; i < len(projectGateways); i += 20 {
			end := i + 20
			if end > len(projectGateways) {
				end = len(projectGateways)
			}
			
			gatewayBatch := projectGateways[i:end]
			macAddrs := make([]string, len(gatewayBatch))
			for j, gateway := range gatewayBatch {
				macAddrs[j] = gateway.MacAddr
			}

			for i := 1; i <= hoursAgo; i++ {
				hourTime := collectTime.Add(time.Duration(-i) * time.Hour)
				hourTime = time.Date(hourTime.Year(), hourTime.Month(), hourTime.Day(),
					hourTime.Hour(), 0, 0, 0, hourTime.Location())

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
					return errors.Wrap(err, "Failed to get box hour stats")
				}

				var resp struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Data    map[string]interface{} `json:"data"`
				}

				if err := json.Unmarshal(respBytes, &resp); err != nil {
					return errors.Wrap(err, "Failed to unmarshal response")
				}

				if resp.Code != "0" {
					return fmt.Errorf("API error: %s", resp.Message)
				}

				// Process response for each gateway
				for _, gateway := range gatewayBatch {
					if statsArr, ok := resp.Data[gateway.MacAddr].([]interface{}); ok {
						stats := processHourStats(statsArr)
						hourlyStats := []model.Eco_gateway_1h{{
							ID:               gateway.ID + "_" + hourTime.Format("2006010215"),
							Time:             common.LocalTime(hourTime),
							GatewayID:        gateway.ID,
							FloorID:          gateway.FloorID,
							BuildingID:       gateway.BuildingID,
							Type:             gateway.Type,
							PowerConsumption: getTotalElectricity(stats),
						}}

						if err := saveGatewayHourlyStats(hourlyStats); err != nil {
							return errors.Wrap(err, "Failed to save gateway hourly stats")
						}
					}
				}
			}
		}
	}

	return nil
}

func processHourStats(data []interface{}) []struct {
	Addr        int     `json:"addr"`
	Electricity float64 `json:"electricity"`
} {
	var stats []struct {
		Addr        int     `json:"addr"`
		Electricity float64 `json:"electricity"`
	}

	for _, stat := range data {
		if statMap, ok := stat.(map[string]interface{}); ok {
			addr, ok1 := statMap["addr"].(float64)
			electricity, ok2 := statMap["electricity"].(float64)
			if !ok1 || !ok2 {
				continue
			}
			stats = append(stats, struct {
				Addr        int     `json:"addr"`
				Electricity float64 `json:"electricity"`
			}{
				Addr:        int(addr),
				Electricity: electricity,
			})
		}
	}

	return stats
}

func getTotalElectricity(stats []struct {
	Addr        int     `json:"addr"`
	Electricity float64 `json:"electricity"`
}) float64 {
	for _, stat := range stats {
		if stat.Addr == 0 {
			return stat.Electricity
		}
	}
	return 0
}

func saveGatewayHourlyStats(stats []model.Eco_gateway_1h) error {
	if len(stats) == 0 {
		return nil
	}

	err := common.DbBatchUpsert(context.Background(), common.GetDaprClient(), stats, model.Eco_gateway_1hTableInfo.Name, model.Eco_gateway_1h_FIELD_NAME_id)
	if err != nil {
		return errors.Wrap(err, "Failed to batch upsert gateway hourly stats")
	}
	
	return nil
}

func refreshContinuousAggregate(collectTime time.Time) error {
	for tableName, refreshType := range needRefreshContinuousAggregateMap {
		startTime := collectTime
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
		default:
			return fmt.Errorf("Invalid refresh type: %s", refreshType)
		}

		if err := common.DbRefreshContinuousAggregate(context.Background(), common.GetDaprClient(),
			tableName, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")); err != nil {
			return errors.Wrapf(err, "Failed to refresh continuous aggregate for table %s", tableName)
		}
	}
	return nil
}

func GetAllEcgateways() ([]model.Ecgateway, error) {
	datas, err := common.DbQuery[model.Ecgateway](context.Background(), common.GetDaprClient(), model.EcgatewayTableInfo.Name, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query gateways")
	}
	return datas, nil
}

func GetAllEcbuildings() ([]model.Ecbuilding, error) {
	datas, err := common.DbQuery[model.Ecbuilding](context.Background(), common.GetDaprClient(), model.EcbuildingTableInfo.Name, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query buildings")
	}
	return datas, nil
}
