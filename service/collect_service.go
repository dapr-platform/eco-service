package service

import (
	"context"
	"crypto/md5"
	"eco-service/client"
	"eco-service/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/exp/rand"
)

var gatewayNeedRefreshContinuousAggregateMap = map[string]string{
	"f_eco_gateway_1d":  "day",
	"f_eco_gateway_1m":  "month",
	"f_eco_gateway_1y":  "year",
	"f_eco_floor_1d":    "day",
	"f_eco_floor_1m":    "month",
	"f_eco_floor_1y":    "year",
	"f_eco_building_1d": "day",
	"f_eco_building_1m": "month",
	"f_eco_building_1y": "year",
	"f_eco_park_1h":     "hour",
	"f_eco_park_1d":     "day",
	"f_eco_park_1m":     "month",
	"f_eco_park_1y":     "year",
}
var waterNeedRefreshContinuousAggregateMap = map[string]string{
	"f_eco_park_water_1m": "month",
	"f_eco_park_water_1y": "year",
}

func init() {
	// Start goroutine to collect stats every hour at 5 minutes past
	go func() {
		for {
			now := time.Now()
			// Calculate next run time (next hour at 5 minutes past)
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 5, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			common.Logger.Info("Starting scheduled data collection...")

			gateways, err := GetAllEcgateways()
			if err != nil {
				common.Logger.Errorf("Failed to get gateways: %v", err)
				continue
			}

			common.Logger.Infof("Found %d gateways to collect data from", len(gateways))

			if err := collectGatewaysHours(time.Now(), 4, gateways); err != nil {
				common.Logger.Errorf("Failed to collect gateway hourly stats: %v", err)
			}

			if err := refreshContinuousAggregate(time.Now().Add(-time.Hour*4), gatewayNeedRefreshContinuousAggregateMap); err != nil {
				common.Logger.Errorf("Failed to refresh continuous aggregates: %v", err)
			}

			common.Logger.Info("Scheduled data collection completed")
		}
	}()
	go func() {
		for {
			now := time.Now()
			timeHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
			demoWaterDataGenHourly(timeHour)
			time.Sleep(time.Hour)
		}
	}()
}

func CheckCollectPower(start, end string) ([]map[string]interface{}, error) {
	selectSql := "SELECT " +
		"DATE_TRUNC('day', time) as day," +
		"park_id," +
		"COUNT(*) as actual_records," +
		"24 as expected_records," +
		"(COUNT(*) * 100.0 / 24) as completeness_percentage"
	fromSql := " f_eco_gateway_1h " +
		"GROUP BY DATE_TRUNC('day', time), park_id" +
		"HAVING COUNT(*) < 24" +
		"ORDER BY day, park_id;"
	whereSql := "1=1"
	data, err := common.CustomSql[map[string]interface{}](context.Background(), common.GetDaprClient(), selectSql, fromSql, whereSql)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func ManuGenDemoWaterData(startDayStr ...string) {
	startTime := time.Now()
	if len(startDayStr) > 0 {
		var err error
		startTime, err = time.Parse("2006-01-02", startDayStr[0])
		if err != nil {
			common.Logger.Errorf("Failed to parse start day: %v", err)
			return
		}
	}
	common.Logger.Infof("Generating demo water data from %s", startTime.Format("2006-01-02"))
	endTime := startTime.Add(time.Hour * 25)
	for currentTime := startTime.Add(time.Hour); !currentTime.After(endTime); currentTime = currentTime.Add(time.Hour) {
		demoWaterDataGenHourly(currentTime)
	}
}
func demoWaterDataGenHourly(startTime time.Time) {
	common.Logger.Infof("Generating demo water data for %s", startTime.Format("2006-01-02 15:04"))
	// 生成上一个小时的时间
	lastHour := startTime.Add(-time.Hour)

	// 根据是否周末决定基准用水量
	baseWaterConsumption := 2000.0
	if lastHour.Weekday() == time.Saturday || lastHour.Weekday() == time.Sunday {
		baseWaterConsumption = 500.0
	}

	// 根据时段分配权重
	hourWeight := 1.0
	hour := lastHour.Hour()
	switch {
	case hour >= 6 && hour <= 9: // 早高峰
		hourWeight = 2.0
	case hour >= 11 && hour <= 13: // 午高峰
		hourWeight = 1.8
	case hour >= 17 && hour <= 19: // 晚高峰
		hourWeight = 1.5
	case hour >= 1 && hour <= 5: // 凌晨
		hourWeight = 0.3
	}

	// 生成随机数
	rand.Seed(uint64(time.Now().UnixNano()))
	// 计算这个小时的用水量：基准值/24*权重，上下浮动20%
	waterConsumption := (baseWaterConsumption/24.0)*hourWeight + (rand.Float64()-0.5)*baseWaterConsumption/24.0*0.4

	park, err := common.DbGetOne[model.Ecpark](context.Background(), common.GetDaprClient(), model.EcparkTableInfo.Name, "")
	if err != nil {
		common.Logger.Errorf("Failed to get park: %v", err)
		return
	}

	// 构造数据
	waterData := model.Eco_park_water_1h{
		ID:               fmt.Sprintf("%x", md5.Sum([]byte(park.ID+"_"+lastHour.Format("2006010215")))),
		Time:             common.LocalTime(lastHour),
		ParkID:           park.ID,
		WaterConsumption: waterConsumption,
	}

	// 插入数据
	err = common.DbUpsert(context.Background(), common.GetDaprClient(), waterData, model.Eco_park_water_1hTableInfo.Name, model.Eco_park_water_1h_FIELD_NAME_id)
	if err != nil {
		common.Logger.Errorf("Failed to insert water consumption data: %v", err)
		return
	}

	// 刷新连续聚合表
	err = refreshContinuousAggregate(lastHour, waterNeedRefreshContinuousAggregateMap)
	if err != nil {
		common.Logger.Errorf("Failed to refresh water continuous aggregates: %v", err)
		return
	}

	common.Logger.Infof("Generated water consumption data for %s: %.2f", lastHour.Format("2006-01-02 15:04"), waterConsumption)
}

// 手动收集指定日期的数据，开始时间结束时间格式为 2024-01-01
func ManuCollectGatewayHourlyStatsByDay(start, end string) error {
	common.Logger.Infof("Starting manual data collection from %s to %s", start, end)

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

	common.Logger.Infof("Found %d gateways to collect data from", len(gateways))

	if len(gateways) == 0 {
		return errors.New("No gateways found")
	}

	// Iterate through each day
	for currentDate := startTime; !currentDate.After(endTime); currentDate = currentDate.AddDate(0, 0, 1) {
		common.Logger.Infof("Collecting data for date: %s", currentDate.Format("2006-01-02"))

		if err := collectGatewaysFullDay(currentDate, gateways); err != nil {
			common.Logger.Errorf("Failed to collect stats for %s: %v",
				currentDate.Format("2006-01-02"), err)
			return err
		}

		if err := refreshContinuousAggregate(currentDate, gatewayNeedRefreshContinuousAggregateMap); err != nil {
			return err
		}

		common.Logger.Infof("Successfully collected and processed data for %s", currentDate.Format("2006-01-02"))
	}

	return nil
}
func ManuFillGatewayHourStats(month, value string) error {
	// Parse month string to time
	startTime, err := time.Parse("2006-01", month)
	if err != nil {
		return errors.Wrap(err, "Failed to parse month")
	}

	// Parse value to float64
	totalValue, err := cast.ToFloat64E(value)
	if err != nil {
		return errors.Wrap(err, "Failed to parse value")
	}

	// Get all gateways
	gateways, err := GetAllEcgateways()
	if err != nil {
		return errors.Wrap(err, "Failed to get gateways")
	}

	if len(gateways) == 0 {
		return errors.New("No gateways found")
	}

	// Calculate days in month
	endTime := startTime.AddDate(0, 1, 0)
	totalDays := int(endTime.Sub(startTime).Hours()) / 24

	// Calculate base value per day
	baseValuePerDay := totalValue / float64(totalDays)

	// Generate and save hourly stats for each day
	for currentDay := startTime; currentDay.Before(endTime); currentDay = currentDay.AddDate(0, 0, 1) {
		// Calculate daily value with random variation (-10% to +10%)
		randomFactor := 0.9 + (rand.Float64() * 0.2) // 0.9 to 1.1
		dailyValue := baseValuePerDay * randomFactor

		// Distribute daily value across hours with different patterns based on day type
		isWeekend := currentDay.Weekday() == time.Saturday || currentDay.Weekday() == time.Sunday
		hourlyDistribution := make([]float64, 24)

		for hour := 0; hour < 24; hour++ {
			if isWeekend {
				// Weekend pattern: more even distribution
				hourlyDistribution[hour] = 1.0
			} else {
				// Weekday pattern: peak during working hours
				switch {
				case hour < 6: // Night (0-5)
					hourlyDistribution[hour] = 0.3
				case hour < 9: // Morning ramp-up (6-8)
					hourlyDistribution[hour] = 0.8
				case hour < 18: // Working hours (9-17)
					hourlyDistribution[hour] = 1.5
				case hour < 22: // Evening (18-21)
					hourlyDistribution[hour] = 1.0
				default: // Late night (22-23)
					hourlyDistribution[hour] = 0.5
				}
			}
		}

		// Normalize distribution
		var sum float64
		for _, v := range hourlyDistribution {
			sum += v
		}
		for i := range hourlyDistribution {
			hourlyDistribution[i] = hourlyDistribution[i] / sum * dailyValue
		}

		// Save hourly stats
		for hour := 0; hour < 24; hour++ {
			currentTime := currentDay.Add(time.Duration(hour) * time.Hour)
			var hourlyStats []model.Eco_gateway_1h

			hourValue := hourlyDistribution[hour] / float64(len(gateways))
			for _, gateway := range gateways {
				// Add small random variation per gateway (-5% to +5%)
				gatewayFactor := 0.95 + (rand.Float64() * 0.1)
				stat := model.Eco_gateway_1h{
					ID:               gateway.ID + "_" + currentTime.Format("2006010215"),
					Time:             common.LocalTime(currentTime),
					GatewayID:        gateway.ID,
					FloorID:          gateway.FloorID,
					BuildingID:       gateway.BuildingID,
					Type:             gateway.Type,
					ParkID:           gateway.ParkID,
					PowerConsumption: hourValue * gatewayFactor,
				}
				hourlyStats = append(hourlyStats, stat)
			}

			if err := saveGatewayHourlyStats(hourlyStats); err != nil {
				return errors.Wrapf(err, "Failed to save hourly stats for time %s", currentTime.Format("2006-01-02 15:04"))
			}
		}
	}

	// Refresh continuous aggregates
	if err := refreshContinuousAggregate(startTime, gatewayNeedRefreshContinuousAggregateMap); err != nil {
		return errors.Wrap(err, "Failed to refresh continuous aggregates")
	}

	return nil
}

func DebugGetBoxHourStats(mac string, year string, month string, day string) (map[string]interface{}, error) {
	box, err := common.DbGetOne[model.Ecgateway](context.Background(), common.GetDaprClient(), model.EcgatewayTableInfo.Name, "mac_addr="+mac)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get gateway %s", mac)
	}
	reqBody := map[string]string{
		"projectCode": box.ProjectCode,
		"mac":         mac,
		"year":        year,
		"month":       month,
		"day":         day,
	}

	common.Logger.Infof("Requesting data for batch of %d gateways, date: %s", 1,
		time.Date(cast.ToInt(year), time.Month(cast.ToInt(month)), cast.ToInt(day), 0, 0, 0, 0, time.Local).Format("2006-01-02"))

	respBytes, err := client.GetBoxesHourStats(reqBody)
	if err != nil {
		common.Logger.Errorf("API request failed: %v", err)
		return nil, errors.Wrap(err, "Failed to get box hour stats")
	}

	var resp struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(respBytes, &resp); err != nil {
		common.Logger.Errorf("Failed to parse API response: %v", err)
		return nil, errors.Wrap(err, "Failed to unmarshal response")
	}

	if resp.Code != "0" {
		common.Logger.Errorf("API returned error code: %s, message: %s", resp.Code, resp.Message)
		return nil, fmt.Errorf("API error: %s", resp.Message)
	}

	common.Logger.Infof("Received data for %d gateways", len(resp.Data))
	return resp.Data, nil
}

func collectGatewaysFullDay(collectTime time.Time, gateways []model.Ecgateway) error {
	// Group gateways by project code
	projectGateways := make(map[string][]model.Ecgateway)
	for i := range gateways {
		projectCode := gateways[i].ProjectCode
		if len(projectCode) == 0 {
			var err error
			projectCode, err = client.GetBoxProjectCode(gateways[i].MacAddr)
			if err != nil {
				return errors.Wrapf(err, "Failed to get project code for gateway %s", gateways[i].ID)
			}
			gateways[i].ProjectCode = projectCode
			if err := common.DbUpsert[model.Ecgateway](context.Background(), common.GetDaprClient(), gateways[i], model.EcgatewayTableInfo.Name, model.Ecgateway_FIELD_NAME_id); err != nil {
				return errors.Wrapf(err, "Failed to update project code for gateway %s", gateways[i].ID)
			}
		}
		projectGateways[projectCode] = append(projectGateways[projectCode], gateways[i])
	}

	common.Logger.Infof("Grouped gateways into %d projects", len(projectGateways))

	// For each project, collect full day stats
	for projectCode, projectGateways := range projectGateways {
		common.Logger.Infof("Processing project %s with %d gateways", projectCode, len(projectGateways))

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

			common.Logger.Infof("Requesting data for batch of %d gateways, date: %s", len(gatewayBatch),
				collectTime.Format("2006-01-02"))

			respBytes, err := client.GetBoxesHourStats(reqBody)
			if err != nil {
				common.Logger.Errorf("API request failed: %v", err)
				return errors.Wrap(err, "Failed to get box hour stats")
			}

			var resp struct {
				Code    string                 `json:"code"`
				Message string                 `json:"message"`
				Data    map[string]interface{} `json:"data"`
			}

			if err := json.Unmarshal(respBytes, &resp); err != nil {
				common.Logger.Errorf("Failed to parse API response: %v", err)
				return errors.Wrap(err, "Failed to unmarshal response")
			}

			if resp.Code != "0" {
				common.Logger.Errorf("API returned error code: %s, message: %s", resp.Code, resp.Message)
				return fmt.Errorf("API error: %s", resp.Message)
			}

			common.Logger.Infof("Received data for %d gateways", len(resp.Data))

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
								ParkID:           gateway.ParkID,
								PowerConsumption: getTotalElectricity(stats),
							}
							hourlyStats = append(hourlyStats, stat)
						}
					}
				} else {
					common.Logger.Warnf("No data found for gateway %s (%s)", gateway.ID, gateway.MacAddr)
				}

				if len(hourlyStats) > 0 {
					common.Logger.Infof("Saving %d hourly stats for gateway %s", len(hourlyStats), gateway.ID)
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

	common.Logger.Infof("Starting collection for last %d hours from %s", hoursAgo, collectTime.Format("2006-01-02 15:04:05"))

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

	common.Logger.Infof("Grouped gateways into %d projects", len(projectGateways))

	// For each project, collect stats for specified hours
	for projectCode, projectGateways := range projectGateways {
		common.Logger.Infof("Processing project %s with %d gateways", projectCode, len(projectGateways))

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

				common.Logger.Infof("Requesting data for batch of %d gateways, hour: %s",
					len(gatewayBatch), hourTime.Format("2006-01-02 15:04:05"))

				respBytes, err := client.GetBoxesHourStats(reqBody)
				if err != nil {
					common.Logger.Errorf("API request failed: %v", err)
					return errors.Wrap(err, "Failed to get box hour stats")
				}

				var resp struct {
					Code    string                 `json:"code"`
					Message string                 `json:"message"`
					Data    map[string]interface{} `json:"data"`
				}

				if err := json.Unmarshal(respBytes, &resp); err != nil {
					common.Logger.Errorf("Failed to parse API response: %v", err)
					return errors.Wrap(err, "Failed to unmarshal response")
				}

				if resp.Code != "0" {
					common.Logger.Errorf("API returned error code: %s, message: %s", resp.Code, resp.Message)
					return fmt.Errorf("API error: %s", resp.Message)
				}

				common.Logger.Infof("Received data for %d gateways", len(resp.Data))

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
							ParkID:           gateway.ParkID,
							Type:             gateway.Type,
							PowerConsumption: getTotalElectricity(stats),
						}}

						common.Logger.Infof("Saving hourly stats for gateway %s, hour: %s",
							gateway.ID, hourTime.Format("2006-01-02 15:04:05"))

						if err := saveGatewayHourlyStats(hourlyStats); err != nil {
							return errors.Wrap(err, "Failed to save gateway hourly stats")
						}
					} else {
						common.Logger.Warnf("No data found for gateway %s (%s)", gateway.ID, gateway.MacAddr)
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
				common.Logger.Warnf("Invalid stat data format: %+v", statMap)
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

func refreshContinuousAggregate(collectTime time.Time, refreshDefineMap map[string]string) error {
	common.Logger.Infof("Starting continuous aggregate refresh for time: %s", collectTime.Format("2006-01-02 15:04:05"))

	for tableName, refreshType := range refreshDefineMap {
		startTime := collectTime
		endTime := startTime

		// 根据刷新类型设置不同的时间范围
		switch refreshType {
		case "hour": // 按小时刷新
			startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
			// 结束时间为第二天0点
			endTime = startTime.AddDate(0, 0, 1)

		case "day": // 按天刷新
			// 设置开始时间为当天0点
			startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
			// 结束时间为第二天0点
			endTime = startTime.AddDate(0, 0, 1)
		case "month": // 按月刷新
			// 设置开始时间为当月1号0点
			startTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, startTime.Location())
			// 结束时间为后两个月的1号0点,这样可以确保当月数据完整性
			endTime = startTime.AddDate(0, 2, 0)
		case "year": // 按年刷新
			// 设置开始时间为当年1月1日0点
			startTime = time.Date(startTime.Year(), 1, 1, 0, 0, 0, 0, startTime.Location())
			// 结束时间为后两年的1月1日0点,这样可以确保当年数据完整性
			endTime = startTime.AddDate(2, 0, 0)
		default:
			// 如果刷新类型不在预期内,返回错误
			return fmt.Errorf("Invalid refresh type: %s", refreshType)
		}

		common.Logger.Infof("Refreshing continuous aggregate for table %s from %s to %s",
			tableName, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

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
