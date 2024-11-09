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

// CollectGatewayHourlyStats collects gateway hourly stats at 5 minutes past each hour
// Gets data for previous 4 hours and stores in database
func CollectGatewayHourlyStats() error {
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
			err = common.DbUpsert[model.Ecgateway](context.Background(), common.GetDaprClient(), gateway, model.EcgatewayTableInfo.Name, gateway.ID)
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

		// Collect stats for previous 4 hours
		for i := 1; i <= 4; i++ {
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

				// Save to database using common.DbSave
				err = common.DbUpsert[model.Eco_gateway_1h](context.Background(), common.GetDaprClient(), hourlyStats, model.Eco_gateway_1hTableInfo.Name, hourlyStats.ID)
				if err != nil {
					return err
				}
			}
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
