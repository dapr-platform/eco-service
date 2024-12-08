package service

import (
	"context"
	"crypto/md5"
	"eco-service/client"
	"eco-service/model"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
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
	"f_eco_floor_1h":    "hour",
	"f_eco_building_1h": "hour",
	"f_eco_building_1d": "day",
	"f_eco_building_1m": "month",
	"f_eco_building_1y": "year",
	"f_eco_park_1h":     "hour",
	"f_eco_park_1d":     "day",
	"f_eco_park_1m":     "month",
	"f_eco_park_1y":     "year",
}
var waterNeedRefreshContinuousAggregateMap = map[string]string{
	"f_eco_park_water_1h": "hour",
	"f_eco_park_water_1d": "day",
	"f_eco_park_water_1m": "month",
	"f_eco_park_water_1y": "year",
}
var COLLECT_TYPE_PLATFORM = 0
var COLLECT_TYPE_IOT = 1

// 初始化函数,启动定时任务收集数据
func init() {
	// Start goroutine to collect stats every hour at 5 minutes past
	go func() {
		for {
			now := time.Now()
			// Calculate next run time (next hour at 5 minutes past)
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 5, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			common.Logger.Info("Starting scheduled data collection...")

			gateways, err := GetAllEcgateways(COLLECT_TYPE_PLATFORM)
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

	// Start goroutine for daily full refresh at midnight
	go func() {
		for {
			now := time.Now()
			// Calculate next run time (midnight)
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			common.Logger.Info("Starting daily full continuous aggregate refresh...")

			if err := refreshContinuousAggregateFull(gatewayNeedRefreshContinuousAggregateMap); err != nil {
				common.Logger.Errorf("Failed to perform full refresh of continuous aggregates: %v", err)
			}

			common.Logger.Info("Daily full continuous aggregate refresh completed")
		}
	}()

	// Start goroutine to collect iot real data every hour at 59 minutes past
	go func() {
		for {
			now := time.Now()
			// Calculate next run time (59 minutes past the hour)
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 0, 0, now.Location())
			if next.Before(now) {
				next = next.Add(time.Hour)
			}
			time.Sleep(next.Sub(now))

			common.Logger.Info("Starting water meter data collection...")
			go func() {
				if err := CollectWaterMeterRealData(); err != nil {
					common.Logger.Errorf("Failed to collect water meter data: %v", err)
				}
			}()
			go func() {
				if err := CollectPowerRealData(); err != nil {
					common.Logger.Errorf("Failed to collect power meter data: %v", err)
				}
			}()
			common.Logger.Info("Iot real data collection completed")
		}
	}()
}

// 强制刷新所有连续聚合数据
func ForceRefreshContinuousAggregate() error {
	err := refreshContinuousAggregateFull(gatewayNeedRefreshContinuousAggregateMap)
	if err != nil {
		return err
	}
	err = refreshContinuousAggregateFull(waterNeedRefreshContinuousAggregateMap)
	if err != nil {
		return err
	}
	return nil
}

// 检查数据收集情况,返回数据缺失的记录
func CheckCollectData(start, end string, collectType int) ([]map[string]interface{}, error) {
	tablename := ""
	totalCount := 0
	if collectType == 0 {
		tablename = "f_eco_gateway_1h"
		gateways, err := GetAllEcgateways(COLLECT_TYPE_PLATFORM)
		if err != nil {
			return nil, err
		}
		totalCount = len(gateways) * 24
	} else {
		tablename = "f_eco_park_water_1h"
		totalCount = 24
	}
	selectSql := `DATE_TRUNC('day', time) as day,
		park_id,
		COUNT(*) as actual_records,
		` + strconv.Itoa(totalCount) + ` as expected_records,
		(COUNT(*) * 100.0 / ` + strconv.Itoa(totalCount) + `) as completeness_percentage `
	fromSql := tablename
	whereSql := "time<'" + end + "' and time>='" + start + "' "
	whereSql += `GROUP BY DATE_TRUNC('day', time), park_id
		HAVING COUNT(*) < ` + strconv.Itoa(totalCount) + `
		ORDER BY day, park_id`
	data, err := common.CustomSql[map[string]interface{}](context.Background(), common.GetDaprClient(), selectSql, fromSql, whereSql)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 收集电力网关实时数据，从IOT平台收集
func CollectPowerRealData() error {
	gateways, err := GetAllEcgateways(COLLECT_TYPE_IOT)
	if err != nil {
		return errors.Wrap(err, "Failed to get gateways")
	}
	common.Logger.Infof("Found %d gateways to collect data from", len(gateways))
	now := time.Now().Truncate(time.Hour)
	for _, gateway := range gateways {
		common.Logger.Infof("Processing gateway: %s (ID: %s)", gateway.CmCode, gateway.ID)
		resp, err := client.GetRealDataByCmCode[client.PowerMeterData](gateway.CmCode)
		if err != nil {
			common.Logger.Errorf("Failed to get real data for meter %s: %v", gateway.CmCode, err)
			continue
		}

		if resp.ResultCode != "0000" {
			common.Logger.Errorf("Error response for meter %s: %s", gateway.CmCode, resp.Message.Message)
			continue
		}

		// Calculate hourly usage by comparing with stored cumulative flow
		currentCumFlow := resp.QueryData.RtData.Energe
		hourlyUsage := currentCumFlow - gateway.RealDataValue
		if hourlyUsage < 0 {
			// 当网关超过量程时，currentCumFlow从0重新开始计算
			// 计算网关的最大量程(根据TotalValue的位数)
			maxValue := math.Pow10(len(strconv.FormatFloat(gateway.RealDataValue, 'f', -1, 64)))
			// 加上最大量程得到正确用量
			hourlyUsage = currentCumFlow + (maxValue - gateway.RealDataValue)
		}

		initial := false
		if gateway.RealDataValue == 0 {
			initial = true
			common.Logger.Infof("Initial data collection for meter %s, setting baseline value to %.2f", gateway.CmCode, currentCumFlow)
		}

		// Update meter's cumulative flow
		gateway.RealDataValue = currentCumFlow
		common.Logger.Debugf("Attempting to update meter %s with data: %+v", gateway.CmCode, gateway)

		if err := common.DbUpsert[model.Ecgateway](context.Background(), common.GetDaprClient(), gateway,
			model.EcgatewayTableInfo.Name, model.Ecgateway_FIELD_NAME_id); err != nil {
			common.Logger.Errorf("Failed to update meter %s: %v", gateway.CmCode, err)
			continue
		}
		common.Logger.Debugf("Updated cumulative flow for meter %s to %.2f", gateway.CmCode, currentCumFlow)

		if initial {
			continue
		}

		// Insert hourly usage into 1h table
		stat := model.Eco_gateway_1h{
			ID:               gateway.ID + "_" + now.Format("2006010215"),
			Time:             common.LocalTime(now),
			GatewayID:        gateway.ID,
			FloorID:          gateway.FloorID,
			BuildingID:       gateway.BuildingID,
			Type:             gateway.Type,
			Level:            gateway.Level,
			ParkID:           gateway.ParkID,
			PowerConsumption: hourlyUsage,
		}

		if err := common.DbUpsert[model.Eco_gateway_1h](context.Background(), common.GetDaprClient(), stat, model.Eco_gateway_1hTableInfo.Name, model.Eco_gateway_1h_FIELD_NAME_id); err != nil {
			common.Logger.Errorf("Failed to insert hourly data for meter %s: %v", gateway.CmCode, err)
			continue
		}

		common.Logger.Infof("Successfully collected data for meter %s: hourly usage %.2f", gateway.CmCode, hourlyUsage)

	}

	return nil
}

// 收集水表实时数据
func CollectWaterMeterRealData() error {
	common.Logger.Info("Starting water meter real-time data collection")

	waterMeters, err := GetAllWaterMeters()
	if err != nil {
		return errors.Wrap(err, "Failed to get water meters")
	}

	common.Logger.Infof("Found %d water meters to collect data from", len(waterMeters))

	// Get current time rounded to hour
	now := time.Now().Truncate(time.Hour)
	common.Logger.Infof("Collecting data for time: %s", now.Format("2006-01-02 15:04:05"))

	for _, meter := range waterMeters {
		common.Logger.Debugf("Processing water meter: %s (ID: %s)", meter.CmCode, meter.ID)

		// Get real-time data for each water meter
		resp, err := client.GetRealDataByCmCode[client.WaterMeterData](meter.CmCode)
		if err != nil {
			common.Logger.Errorf("Failed to get real data for meter %s: %v", meter.CmCode, err)
			continue
		}

		if resp.ResultCode != "0000" {
			common.Logger.Errorf("Error response for meter %s: %s", meter.CmCode, resp.Message.Message)
			continue
		}

		// Calculate hourly usage by comparing with stored cumulative flow
		currentCumFlow := resp.QueryData.RtData.CumFlow
		hourlyUsage := currentCumFlow - meter.TotalValue
		if hourlyUsage < 0 {
			// 当水表超过量程时，currentCumFlow从0重新开始计算
			// 计算水表的最大量程(根据TotalValue的位数)
			maxValue := math.Pow10(len(strconv.FormatFloat(meter.TotalValue, 'f', -1, 64)))
			// 加上最大量程得到正确用量
			hourlyUsage = currentCumFlow + (maxValue - meter.TotalValue)
		}

		initial := false
		if meter.TotalValue == 0 {
			initial = true
			common.Logger.Infof("Initial data collection for meter %s, setting baseline value to %.2f", meter.CmCode, currentCumFlow)
		}

		// Update meter's cumulative flow
		meter.TotalValue = currentCumFlow
		common.Logger.Debugf("Attempting to update meter %s with data: %+v", meter.CmCode, meter)

		if err := common.DbUpsert[model.Ecwater_meter](context.Background(), common.GetDaprClient(), meter,
			model.Ecwater_meterTableInfo.Name, model.Ecwater_meter_FIELD_NAME_id); err != nil {
			common.Logger.Errorf("Failed to update meter %s: %v", meter.CmCode, err)
			continue
		}
		common.Logger.Debugf("Updated cumulative flow for meter %s to %.2f", meter.CmCode, currentCumFlow)

		if initial {
			continue
		}

		// Insert hourly usage into 1h table
		hourData := model.Eco_water_meter_1h{
			ID:               fmt.Sprintf("%x", md5.Sum([]byte(meter.ID+"_"+now.Format("2006010215")))),
			Time:             common.LocalTime(now),
			WaterMeterID:     meter.ID,
			BuildingID:       meter.BuildingID,
			ParkID:           meter.ParkID,
			Type:             meter.Type,
			WaterConsumption: hourlyUsage,
		}

		if err := common.DbUpsert[model.Eco_water_meter_1h](context.Background(), common.GetDaprClient(), hourData, model.Eco_water_meter_1hTableInfo.Name, model.Eco_water_meter_1h_FIELD_NAME_id); err != nil {
			common.Logger.Errorf("Failed to insert hourly data for meter %s: %v", meter.CmCode, err)
			continue
		}

		common.Logger.Infof("Successfully collected data for meter %s: hourly usage %.2f", meter.CmCode, hourlyUsage)
	}
	if err := refreshContinuousAggregateFull(waterNeedRefreshContinuousAggregateMap); err != nil {
		return errors.Wrap(err, "Failed to refresh continuous aggregates")
	}
	common.Logger.Info("Completed water meter real-time data collection")
	return nil
}

// 手动收集指定日期的网关数据
func ManuCollectGatewayHourlyStatsByDay(start, end, macAddr string) error {
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

	gateways, err := GetAllEcgateways(COLLECT_TYPE_PLATFORM)
	if err != nil {
		return errors.Wrap(err, "Failed to get gateways")
	}
	if macAddr != "" {
		filteredGateways := make([]model.Ecgateway, 0)
		for _, gateway := range gateways {
			if gateway.MacAddr == macAddr {
				filteredGateways = append(filteredGateways, gateway)
			}
		}
		gateways = filteredGateways
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

func ManuFillPowerCollectIotData(cmCode, start, end, value string) error {
	qstr := model.Ecgateway_FIELD_NAME_cm_code + "=" + cmCode
	gateway, err := common.DbGetOne[model.Ecgateway](context.Background(), common.GetDaprClient(), model.EcgatewayTableInfo.Name, qstr)
	if err != nil {
		return errors.Wrap(err, "Failed to get gateway")
	}
	if gateway == nil {
		return errors.New("Gateway not found")
	}

	// Parse month string to time
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return errors.Wrap(err, "Failed to parse start")
	}

	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return errors.Wrap(err, "Failed to parse end")
	}

	// 验证开始时间不能大于结束时间
	if startTime.After(endTime) {
		return errors.New("Start time cannot be after end time")
	}

	// Parse value to float64
	totalValue, err := cast.ToFloat64E(value)
	if err != nil {
		return errors.Wrap(err, "Failed to parse value")
	}

	// 验证值必须大于0
	if totalValue <= 0 {
		return errors.New("Value must be greater than 0")
	}

	common.Logger.Infof("Parsed inputs - Start time: %s, Total value: %.2f", startTime.Format("2006-01"), totalValue)

	// Calculate days in month
	totalDays := int(endTime.Sub(startTime).Hours())/24 + 1 // 修复日期计算,加1包含结束日期

	common.Logger.Infof("Generating daily values for %d days", totalDays)

	// First generate daily values that sum to total
	dailyValues := make([]float64, totalDays)
	var dailyTotal float64

	// 使用固定种子以保证可重复性
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate random daily values
	for i := 0; i < totalDays; i++ {
		// Random factor between 0.8 and 1.2
		factor := 0.8 + (rand.Float64() * 0.4)
		dailyValues[i] = factor
		dailyTotal += factor
	}

	// Normalize daily values to sum to total
	for i := range dailyValues {
		dailyValues[i] = (dailyValues[i] / dailyTotal) * totalValue
	}

	// For each day in the month
	for day := 0; day < totalDays; day++ {
		currentDate := startTime.AddDate(0, 0, day)
		dailyValue := dailyValues[day]

		// Generate hourly values for this day
		hourlyValues := make([]float64, 24)
		var hourlyTotal float64

		// Generate random hourly values with peak/off-peak patterns
		isWeekend := currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday
		for hour := 0; hour < 24; hour++ {
			var baseFactor float64
			if isWeekend {
				switch {
				case hour >= 7 && hour < 10: // 早高峰
					baseFactor = 1.5 + (rand.Float64() * 0.3) // 1.5-1.8
				case hour >= 11 && hour < 14: // 午高峰
					baseFactor = 1.3 + (rand.Float64() * 0.3) // 1.3-1.6
				case hour >= 17 && hour < 20: // 晚高峰
					baseFactor = 1.4 + (rand.Float64() * 0.3) // 1.4-1.7
				case hour >= 23 || hour < 6: // 深夜
					baseFactor = 0.2 + (rand.Float64() * 0.2) // 0.2-0.4
				default: // 其他时段
					baseFactor = 0.8 + (rand.Float64() * 0.3) // 0.8-1.1
				}
			} else {
				switch {
				case hour >= 6 && hour < 9: // 早高峰
					baseFactor = 1.8 + (rand.Float64() * 0.4) // 1.8-2.2
				case hour >= 11 && hour < 14: // 午高峰
					baseFactor = 1.5 + (rand.Float64() * 0.3) // 1.5-1.8
				case hour >= 17 && hour < 20: // 晚高峰
					baseFactor = 1.6 + (rand.Float64() * 0.3) // 1.6-1.9
				case hour >= 23 || hour < 5: // 深夜
					baseFactor = 0.1 + (rand.Float64() * 0.2) // 0.1-0.3
				default: // 其他时段
					baseFactor = 0.7 + (rand.Float64() * 0.3) // 0.7-1.0
				}
			}
			hourlyValues[hour] = baseFactor
			hourlyTotal += baseFactor
		}

		// Normalize hourly values to sum to daily value
		for hour := range hourlyValues {
			hourlyValues[hour] = (hourlyValues[hour] / hourlyTotal) * dailyValue
		}

		// Insert hourly values into database
		for hour := 0; hour < 24; hour++ {
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), hour, 0, 0, 0, time.Local)

			stat := model.Eco_gateway_1h{
				ID:               gateway.ID + "_" + timestamp.Format("2006010215"),
				Time:             common.LocalTime(timestamp),
				GatewayID:        gateway.ID,
				FloorID:          gateway.FloorID,
				BuildingID:       gateway.BuildingID,
				Type:             gateway.Type,
				Level:            gateway.Level,
				ParkID:           gateway.ParkID,
				PowerConsumption: hourlyValues[hour],
			}

			err := common.DbUpsert(context.Background(), common.GetDaprClient(), stat, model.Eco_gateway_1hTableInfo.Name, model.Eco_gateway_1h_FIELD_NAME_id)
			if err != nil {
				return errors.Wrapf(err, "Failed to insert hour stats for %s", timestamp.Format("2006-01-02 15:04:05"))
			}
		}

		common.Logger.Infof("Successfully inserted hourly values for %s", currentDate.Format("2006-01-02"))
	}

	// Refresh continuous aggregates
	if err := refreshContinuousAggregate(startTime, gatewayNeedRefreshContinuousAggregateMap); err != nil {
		return errors.Wrap(err, "Failed to refresh continuous aggregates")
	}

	common.Logger.Info("Successfully completed ManuFillParkWaterHourStats")
	return nil
}

// 手动填充园区水表小时数据
func ManuFillParkWaterHourStats(cmCode, start, end, value string) error {
	common.Logger.Infof("Starting ManuFillParkWaterHourStats for start: %s, end: %s, value: %s", start, end, value)

	// Parse month string to time
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return errors.Wrap(err, "Failed to parse start")
	}

	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return errors.Wrap(err, "Failed to parse end")
	}

	// 验证开始时间不能大于结束时间
	if startTime.After(endTime) {
		return errors.New("Start time cannot be after end time")
	}

	// Parse value to float64
	totalValue, err := cast.ToFloat64E(value)
	if err != nil {
		return errors.Wrap(err, "Failed to parse value")
	}

	// 验证值必须大于0
	if totalValue <= 0 {
		return errors.New("Value must be greater than 0")
	}

	common.Logger.Infof("Parsed inputs - Start time: %s, Total value: %.2f", startTime.Format("2006-01"), totalValue)

	waterMeters, err := GetAllWaterMeters()
	if err != nil {
		return errors.Wrap(err, "Failed to get water meters")
	}
	if len(waterMeters) == 0 {
		return errors.New("No water meters found")
	}

	// 验证cmCode不能为空
	if cmCode == "" {
		return errors.New("CM code cannot be empty")
	}

	var waterMeter *model.Ecwater_meter
	for _, meter := range waterMeters {
		if meter.CmCode == cmCode {
			waterMeter = &meter
			break
		}
	}
	if waterMeter == nil {
		return errors.New("Water meter not found")
	}

	// Calculate days in month
	totalDays := int(endTime.Sub(startTime).Hours())/24 + 1 // 修复日期计算,加1包含结束日期

	common.Logger.Infof("Generating daily values for %d days", totalDays)

	// First generate daily values that sum to total
	dailyValues := make([]float64, totalDays)
	var dailyTotal float64

	// 使用固定种子以保证可重复性
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate random daily values
	for i := 0; i < totalDays; i++ {
		// Random factor between 0.8 and 1.2
		factor := 0.8 + (rand.Float64() * 0.4)
		dailyValues[i] = factor
		dailyTotal += factor
	}

	// Normalize daily values to sum to total
	for i := range dailyValues {
		dailyValues[i] = (dailyValues[i] / dailyTotal) * totalValue
	}

	// For each day in the month
	for day := 0; day < totalDays; day++ {
		currentDate := startTime.AddDate(0, 0, day)
		dailyValue := dailyValues[day]

		// Generate hourly values for this day
		hourlyValues := make([]float64, 24)
		var hourlyTotal float64

		// Generate random hourly values with peak/off-peak patterns
		isWeekend := currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday
		for hour := 0; hour < 24; hour++ {
			var baseFactor float64
			if isWeekend {
				switch {
				case hour >= 7 && hour < 10: // 早高峰
					baseFactor = 1.5 + (rand.Float64() * 0.3) // 1.5-1.8
				case hour >= 11 && hour < 14: // 午高峰
					baseFactor = 1.3 + (rand.Float64() * 0.3) // 1.3-1.6
				case hour >= 17 && hour < 20: // 晚高峰
					baseFactor = 1.4 + (rand.Float64() * 0.3) // 1.4-1.7
				case hour >= 23 || hour < 6: // 深夜
					baseFactor = 0.2 + (rand.Float64() * 0.2) // 0.2-0.4
				default: // 其他时段
					baseFactor = 0.8 + (rand.Float64() * 0.3) // 0.8-1.1
				}
			} else {
				switch {
				case hour >= 6 && hour < 9: // 早高峰
					baseFactor = 1.8 + (rand.Float64() * 0.4) // 1.8-2.2
				case hour >= 11 && hour < 14: // 午高峰
					baseFactor = 1.5 + (rand.Float64() * 0.3) // 1.5-1.8
				case hour >= 17 && hour < 20: // 晚高峰
					baseFactor = 1.6 + (rand.Float64() * 0.3) // 1.6-1.9
				case hour >= 23 || hour < 5: // 深夜
					baseFactor = 0.1 + (rand.Float64() * 0.2) // 0.1-0.3
				default: // 其他时段
					baseFactor = 0.7 + (rand.Float64() * 0.3) // 0.7-1.0
				}
			}
			hourlyValues[hour] = baseFactor
			hourlyTotal += baseFactor
		}

		// Normalize hourly values to sum to daily value
		for hour := range hourlyValues {
			hourlyValues[hour] = (hourlyValues[hour] / hourlyTotal) * dailyValue
		}

		// Insert hourly values into database
		for hour := 0; hour < 24; hour++ {
			timestamp := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), hour, 0, 0, 0, time.Local)

			waterData := model.Eco_water_meter_1h{
				ID:               fmt.Sprintf("%x", md5.Sum([]byte(waterMeter.ID+"_"+timestamp.Format("2006010215")))),
				Time:             common.LocalTime(timestamp),
				ParkID:           waterMeter.ParkID,
				WaterMeterID:     waterMeter.ID,
				BuildingID:       waterMeter.BuildingID,
				Type:             waterMeter.Type,
				WaterConsumption: hourlyValues[hour],
			}

			err := common.DbUpsert(context.Background(), common.GetDaprClient(), waterData, model.Eco_water_meter_1hTableInfo.Name, model.Eco_water_meter_1h_FIELD_NAME_id)
			if err != nil {
				return errors.Wrapf(err, "Failed to insert hour stats for %s", timestamp.Format("2006-01-02 15:04:05"))
			}
		}

		common.Logger.Infof("Successfully inserted hourly values for %s", currentDate.Format("2006-01-02"))
	}

	// Refresh continuous aggregates
	if err := refreshContinuousAggregate(startTime, waterNeedRefreshContinuousAggregateMap); err != nil {
		return errors.Wrap(err, "Failed to refresh continuous aggregates")
	}

	common.Logger.Info("Successfully completed ManuFillParkWaterHourStats")
	return nil
}

// 手动填充网关小时数据
func ManuFillGatewayHourStats(month, value string) error {
	common.Logger.Infof("Starting ManuFillGatewayHourStats for month: %s, value: %s", month, value)

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

	// Round totalValue to 2 decimal places
	totalValue = float64(int64(totalValue*100)) / 100

	common.Logger.Infof("Parsed inputs - Start time: %s, Total value: %.2f", startTime.Format("2006-01"), totalValue)

	// Get all gateways
	gateways, err := GetAllEcgateways(COLLECT_TYPE_PLATFORM)
	if err != nil {
		return errors.Wrap(err, "Failed to get gateways")
	}
	buildingGateways := []model.Ecgateway{}
	for _, gateway := range gateways {
		if gateway.Level == 1 {
			buildingGateways = append(buildingGateways, gateway)
		}
	}
	gateways = buildingGateways

	if len(gateways) == 0 {
		return errors.New("No gateways found")
	}

	common.Logger.Infof("Found %d gateways to distribute data to", len(gateways))

	// Calculate days in month
	endTime := startTime.AddDate(0, 1, 0)
	totalDays := int(endTime.Sub(startTime).Hours()) / 24

	common.Logger.Infof("Generating daily values for %d days", totalDays)

	// First generate daily values that sum to total
	dailyValues := make([]float64, totalDays)
	var dailyTotal float64

	// Generate random daily values
	for i := 0; i < totalDays; i++ {
		// Random factor between 0.8 and 1.2
		factor := 0.8 + (rand.Float64() * 0.4)
		dailyValues[i] = factor
		dailyTotal += factor
	}

	// Normalize daily values to sum to total and round to 2 decimal places
	var sumDailyValues float64
	for i := range dailyValues[:len(dailyValues)-1] {
		dailyValues[i] = float64(int64((dailyValues[i]/dailyTotal)*totalValue*100)) / 100
		sumDailyValues += dailyValues[i]
	}
	// Last day takes the remaining value to ensure exact total
	dailyValues[len(dailyValues)-1] = float64(int64((totalValue-sumDailyValues)*100)) / 100

	common.Logger.Infof("Generated daily values - First day: %.2f, Last day: %.2f", dailyValues[0], dailyValues[len(dailyValues)-1])

	dayIndex := 0
	for currentDay := startTime; currentDay.Before(endTime); currentDay = currentDay.AddDate(0, 0, 1) {
		dailyValue := dailyValues[dayIndex]
		common.Logger.Debugf("Processing day %s with value %.2f", currentDay.Format("2006-01-02"), dailyValue)

		// Generate hourly distribution for this day
		isWeekend := currentDay.Weekday() == time.Saturday || currentDay.Weekday() == time.Sunday
		hourlyValues := make([]float64, 24)
		var hourlyTotal float64

		// Generate random hourly values with peak/off-peak patterns
		for hour := 0; hour < 24; hour++ {
			var baseFactor float64
			if isWeekend {
				baseFactor = 0.8 + (rand.Float64() * 0.4) // 0.8-1.2 for weekends
			} else {
				switch {
				case hour < 6: // Night (0-5)
					baseFactor = 0.2 + (rand.Float64() * 0.2) // 0.2-0.4
				case hour < 9: // Morning ramp-up (6-8)
					baseFactor = 0.6 + (rand.Float64() * 0.4) // 0.6-1.0
				case hour < 18: // Working hours (9-17)
					baseFactor = 1.3 + (rand.Float64() * 0.4) // 1.3-1.7
				case hour < 22: // Evening (18-21)
					baseFactor = 0.8 + (rand.Float64() * 0.4) // 0.8-1.2
				default: // Late night (22-23)
					baseFactor = 0.3 + (rand.Float64() * 0.4) // 0.3-0.7
				}
			}
			hourlyValues[hour] = baseFactor
			hourlyTotal += baseFactor
		}

		// Normalize hourly values to sum to daily value and round to 2 decimal places
		var sumHourlyValues float64
		for hour := 0; hour < 23; hour++ {
			hourlyValues[hour] = float64(int64((hourlyValues[hour]/hourlyTotal)*dailyValue*100)) / 100
			sumHourlyValues += hourlyValues[hour]
		}
		// Last hour takes the remaining value to ensure exact daily total
		hourlyValues[23] = float64(int64((dailyValue-sumHourlyValues)*100)) / 100

		// Save stats for this hour
		for hour := 0; hour < 24; hour++ {
			var hourlyStats []model.Eco_gateway_1h
			currentTime := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), hour, 0, 0, 0, currentDay.Location())

			// Calculate value per gateway, ensuring total matches input
			valuePerGateway := float64(int64((hourlyValues[hour]/float64(len(gateways)))*100)) / 100
			remainingValue := hourlyValues[hour] - (valuePerGateway * float64(len(gateways)-1))

			common.Logger.Debugf("Hour %02d:00 - Total value: %.4f, Per gateway: %.4f, Last gateway: %.4f",
				hour, hourlyValues[hour], valuePerGateway, remainingValue)

			for i, gateway := range gateways {
				value := valuePerGateway
				if i == len(gateways)-1 {
					value = remainingValue // Last gateway gets remaining value
				}

				stat := model.Eco_gateway_1h{
					ID:               gateway.ID + "_" + currentTime.Format("2006010215"),
					Time:             common.LocalTime(currentTime),
					GatewayID:        gateway.ID,
					FloorID:          gateway.FloorID,
					BuildingID:       gateway.BuildingID,
					Type:             gateway.Type,
					Level:            gateway.Level,
					ParkID:           gateway.ParkID,
					PowerConsumption: value,
				}
				hourlyStats = append(hourlyStats, stat)
			}

			if err := saveGatewayHourlyStats(hourlyStats); err != nil {
				return errors.Wrapf(err, "Failed to save hourly stats for time %s", currentTime.Format("2006-01-02 15:04"))
			}
		}

		dayIndex++
	}

	common.Logger.Infof("Successfully generated and saved hourly stats for %d days", dayIndex)

	// Refresh continuous aggregates
	if err := refreshContinuousAggregate(startTime, gatewayNeedRefreshContinuousAggregateMap); err != nil {
		return errors.Wrap(err, "Failed to refresh continuous aggregates")
	}

	common.Logger.Info("Completed ManuFillGatewayHourStats successfully")
	return nil
}

// 调试获取网关小时数据
func DebugGetBoxHourStats(mac string, year string, month string, day string) (map[string]interface{}, error) {

	projectCode, err := client.GetBoxProjectCode(mac)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get project code for gateway %s", mac)
	}

	reqBody := map[string]string{
		"projectCode": projectCode,
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

// 调试获取网关小时数据
func DebugGetBoxMonthStats(mac string, year string, month string) (map[string]interface{}, error) {

	projectCode, err := client.GetBoxProjectCode(mac)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get project code for gateway %s", mac)
	}

	reqBody := map[string]string{
		"projectCode": projectCode,
		"mac":         mac,
		"year":        year,
		"month":       month,
	}

	common.Logger.Infof("Requesting data for batch of %d gateways, date: %s", 1,
		time.Date(cast.ToInt(year), time.Month(cast.ToInt(month)), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"))

	respBytes, err := client.GetBoxesMonthStats(reqBody)
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

// 收集网关全天数据
func collectGatewaysFullDay(collectTime time.Time, gateways []model.Ecgateway) error {
	// Group gateways by project code
	projectGateways := make(map[string][]model.Ecgateway)
	for i := range gateways {
		projectCode := gateways[i].ProjectCode
		if len(projectCode) == 0 {
			var err error
			if strings.HasSuffix(gateways[i].MacAddr, "_1") {
				projectCode, err = client.GetBoxProjectCode(strings.TrimSuffix(gateways[i].MacAddr, "_1"))
			} else if strings.HasSuffix(gateways[i].MacAddr, "_2") {
				projectCode, err = client.GetBoxProjectCode(strings.TrimSuffix(gateways[i].MacAddr, "_2"))
			} else {
				projectCode, err = client.GetBoxProjectCode(gateways[i].MacAddr)
			}
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
				if strings.HasSuffix(gateway.MacAddr, "_1") {
					macAddrs[j] = strings.TrimSuffix(gateway.MacAddr, "_1")
				} else if strings.HasSuffix(gateway.MacAddr, "_2") {
					macAddrs[j] = strings.TrimSuffix(gateway.MacAddr, "_2")
				} else {
					macAddrs[j] = gateway.MacAddr
				}
			}

			reqBody := map[string]string{
				"projectCode": projectCode,
				"mac":         strings.Join(macAddrs, ","),
				"year":        collectTime.Format("2006"),
				"month":       collectTime.Format("01"),
				"day":         collectTime.Format("02"),
			}

			common.Logger.Infof("Requesting data for batch of %d gateways, date: %s",
				len(gatewayBatch), collectTime.Format("2006-01-02"))

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
				var macAddr string
				addr := 0
				if strings.HasSuffix(gateway.MacAddr, "_1") {
					macAddr = strings.TrimSuffix(gateway.MacAddr, "_1")
					addr = 1
				} else if strings.HasSuffix(gateway.MacAddr, "_2") {
					macAddr = strings.TrimSuffix(gateway.MacAddr, "_2")
					addr = 2
				} else {
					macAddr = gateway.MacAddr
				}
				if gatewayData, ok := resp.Data[macAddr].(map[string]interface{}); ok {
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
								Level:            gateway.Level,
								ParkID:           gateway.ParkID,
								PowerConsumption: getTotalElectricity(stats, addr),
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
			if strings.HasSuffix(gateway.MacAddr, "_1") {
				projectCode, err = client.GetBoxProjectCode(strings.TrimSuffix(gateway.MacAddr, "_1"))
			} else if strings.HasSuffix(gateway.MacAddr, "_2") {
				projectCode, err = client.GetBoxProjectCode(strings.TrimSuffix(gateway.MacAddr, "_2"))
			} else {
				projectCode, err = client.GetBoxProjectCode(gateway.MacAddr)
			}
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
				if strings.HasSuffix(gateway.MacAddr, "_1") {
					macAddrs[j] = strings.TrimSuffix(gateway.MacAddr, "_1")
				} else if strings.HasSuffix(gateway.MacAddr, "_2") {
					macAddrs[j] = strings.TrimSuffix(gateway.MacAddr, "_2")
				} else {
					macAddrs[j] = gateway.MacAddr
				}
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
					var macAddr string
					addr := 0
					if strings.HasSuffix(gateway.MacAddr, "_1") {
						macAddr = strings.TrimSuffix(gateway.MacAddr, "_1")
						addr = 1
					} else if strings.HasSuffix(gateway.MacAddr, "_2") {
						macAddr = strings.TrimSuffix(gateway.MacAddr, "_2")
						addr = 2
					} else {
						macAddr = gateway.MacAddr
					}
					if statsArr, ok := resp.Data[macAddr].([]interface{}); ok {
						stats := processHourStats(statsArr)
						hourlyStats := []model.Eco_gateway_1h{{
							ID:               gateway.ID + "_" + hourTime.Format("2006010215"),
							Time:             common.LocalTime(hourTime),
							GatewayID:        gateway.ID,
							FloorID:          gateway.FloorID,
							BuildingID:       gateway.BuildingID,
							ParkID:           gateway.ParkID,
							Type:             gateway.Type,
							Level:            gateway.Level,
							PowerConsumption: getTotalElectricity(stats, addr),
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
}, addr int) float64 {
	for _, stat := range stats {
		if stat.Addr == addr {
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

func refreshContinuousAggregateFull(refreshDefineMap map[string]string) error {
	for tableName := range refreshDefineMap {
		if err := common.DbRefreshContinuousAggregateFull(context.Background(), common.GetDaprClient(), tableName); err != nil {
			return errors.Wrapf(err, "Failed to refresh continuous aggregate for table %s", tableName)
		}
		common.Logger.Infof("Refreshed continuous aggregate for table %s", tableName)
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

func GetAllEcgateways(collectType int) ([]model.Ecgateway, error) {
	qstr := model.Ecgateway_FIELD_NAME_collect_type + "=" + strconv.Itoa(collectType)
	datas, err := common.DbQuery[model.Ecgateway](context.Background(), common.GetDaprClient(), model.EcgatewayTableInfo.Name, qstr)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query gateways")
	}
	return datas, nil
}
func GetAllWaterMeters() ([]model.Ecwater_meter, error) {
	datas, err := common.DbQuery[model.Ecwater_meter](context.Background(), common.GetDaprClient(), model.Ecwater_meterTableInfo.Name, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query water meters")
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
