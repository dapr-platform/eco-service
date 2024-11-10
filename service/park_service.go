package service

import (
	"context"
	"eco-service/entity"
	"eco-service/model"
	"fmt"
	"time"

	"github.com/dapr-platform/common"
)

const (
	CARBON_FACTOR = 0.61   // 碳排放系数
	COAL_FACTOR   = 0.1229 // 标准煤系数
)

type ParkDataGetter func(time.Time) ([]entity.LabelData, error)

func GetParkWaterConsumption(period string, queryTime time.Time) ([]entity.LabelData, error) {
	switch period {
	case PERIOD_DAY:
		return getParkDataDay(queryTime, 0)
	case PERIOD_MONTH:
		return getParkDataMonth(queryTime, 0)
	case PERIOD_YEAR:
		return getParkDataYear(queryTime, 0)
	}

}

func GetParkCarbonEmissionRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	result, err := GetParkPowerConsumptionRange(period, queryTime, gatewayType)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.Value *= CARBON_FACTOR
	}

	return result, nil
}
func GetParkStandardCoalRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	result, err := GetParkPowerConsumptionRange(period, queryTime, gatewayType)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.Value *= COAL_FACTOR
	}

	return result, nil
}

func GetParkPowerConsumptionRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	var result []entity.LabelData
	var err error

	switch period {
	case PERIOD_DAY:
		// 获取24小时数据
		startTime := time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		endTime := startTime.Add(24 * time.Hour)
		result, err = getParkDataHourRange(startTime, endTime, gatewayType)

	case PERIOD_MONTH:
		// 获取整月所有天数据
		startTime := time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		endTime := time.Date(queryTime.Year(), queryTime.Month()+1, 0, 23, 59, 59, 999999999, queryTime.Location())
		result, err = getParkDataDayRange(startTime, endTime, gatewayType)

	case PERIOD_YEAR:
		// 获取全年12月数据
		startTime := time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		endTime := time.Date(queryTime.Year()+1, 1, 1, 0, 0, 0, 0, queryTime.Location())
		result, err = getParkDataMonthRange(startTime, endTime, gatewayType)

	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	return result, err
}

func GetParkPowerConsumption(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	getters := map[string]ParkDataGetter{
		PERIOD_HOUR:  func(t time.Time) ([]entity.LabelData, error) { return getParkDataHour(t, gatewayType) },
		PERIOD_DAY:   func(t time.Time) ([]entity.LabelData, error) { return getParkDataDay(t, gatewayType) },
		PERIOD_MONTH: func(t time.Time) ([]entity.LabelData, error) { return getParkDataMonth(t, gatewayType) },
		PERIOD_YEAR:  func(t time.Time) ([]entity.LabelData, error) { return getParkDataYear(t, gatewayType) },
	}

	if getter, ok := getters[period]; ok {
		return getter(queryTime)
	}

	return nil, fmt.Errorf("unsupported period: %s", period)
}

func getParkDataWithTimeOffset(period string, queryTime time.Time, years, months, days, hours int, gatewayType int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days).Add(time.Duration(hours) * time.Hour)

	whereClause := ""
	if gatewayType > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType)
	}

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_park_1hTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1h](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02 15"), whereClause),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	result := make([]entity.LabelData, 0)
	parkPowerMap := make(map[string]float64)

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_1h) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_1d) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_1m) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_1y) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	}

	for parkID, powerConsumption := range parkPowerMap {
		result = append(result, entity.LabelData{
			Id:    parkID,
			Label: parkID,
			Value: powerConsumption,
		})
	}

	return result, nil
}

func getParkDataWithTimeRange(period string, startTime time.Time, endTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	whereClause := ""
	if gatewayType > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType)
	}

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_park_1hTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1h](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time$gte.%s&time$lte.%s%s", startTime.Format("2006-01-02T15"), endTime.Format("2006-01-02T15"), whereClause),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time$gte.%s&time$lte.%s%s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time$gte.%s&time$lte.%s%s", startTime.Format("2006-01"), endTime.Format("2006-01"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time$gte.%s&time$lte.%s%s", startTime.Format("2006"), endTime.Format("2006"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	result := make([]entity.LabelData, 0)
	parkPowerMap := make(map[string]float64)

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_1h) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_1d) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_1m) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_1y) {
			if gatewayType == 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	}

	for parkID, powerConsumption := range parkPowerMap {
		result = append(result, entity.LabelData{
			Id:    parkID,
			Label: parkID,
			Value: powerConsumption,
		})
	}

	return result, nil
}
func getParkDataHourRange(startTime time.Time, endTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeRange(PERIOD_HOUR, startTime, endTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一小时同时段)
	hbStartTime := startTime.Add(-1 * time.Hour)
	hbEndTime := endTime.Add(-1 * time.Hour)
	hb, err := getParkDataWithTimeRange(PERIOD_HOUR, hbStartTime, hbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(昨天同时段)
	tbStartTime := startTime.AddDate(0, 0, -1)
	tbEndTime := endTime.AddDate(0, 0, -1)
	tb, err := getParkDataWithTimeRange(PERIOD_HOUR, tbStartTime, tbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataDayRange(startTime time.Time, endTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeRange(PERIOD_DAY, startTime, endTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一天同时段)
	hbStartTime := startTime.AddDate(0, 0, -1)
	hbEndTime := endTime.AddDate(0, 0, -1)
	hb, err := getParkDataWithTimeRange(PERIOD_DAY, hbStartTime, hbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(上月同时段)
	tbStartTime := startTime.AddDate(0, -1, 0)
	tbEndTime := endTime.AddDate(0, -1, 0)
	tb, err := getParkDataWithTimeRange(PERIOD_DAY, tbStartTime, tbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataMonthRange(startTime time.Time, endTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeRange(PERIOD_MONTH, startTime, endTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(上月同时段)
	hbStartTime := startTime.AddDate(0, -1, 0)
	hbEndTime := endTime.AddDate(0, -1, 0)
	hb, err := getParkDataWithTimeRange(PERIOD_MONTH, hbStartTime, hbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同时段)
	tbStartTime := startTime.AddDate(-1, 0, 0)
	tbEndTime := endTime.AddDate(-1, 0, 0)
	tb, err := getParkDataWithTimeRange(PERIOD_MONTH, tbStartTime, tbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataYearRange(startTime time.Time, endTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeRange(PERIOD_YEAR, startTime, endTime, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同时段)
	tbStartTime := startTime.AddDate(-1, 0, 0)
	tbEndTime := endTime.AddDate(-1, 0, 0)
	tb, err := getParkDataWithTimeRange(PERIOD_YEAR, tbStartTime, tbEndTime, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, nil, tb)
	return current, nil
}

func getParkDataHour(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeOffset(PERIOD_HOUR, queryTime, 0, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一小时)
	hb, err := getParkDataWithTimeOffset(PERIOD_HOUR, queryTime, 0, 0, 0, -1, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(昨天同一小时)
	tb, err := getParkDataWithTimeOffset(PERIOD_HOUR, queryTime, 0, 0, -1, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataDay(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeOffset(PERIOD_DAY, queryTime, 0, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一天)
	hb, err := getParkDataWithTimeOffset(PERIOD_DAY, queryTime, 0, 0, -1, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(上月同天)
	tb, err := getParkDataWithTimeOffset(PERIOD_DAY, queryTime, 0, -1, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataMonth(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeOffset(PERIOD_MONTH, queryTime, 0, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(上月)
	hb, err := getParkDataWithTimeOffset(PERIOD_MONTH, queryTime, 0, -1, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同月)
	tb, err := getParkDataWithTimeOffset(PERIOD_MONTH, queryTime, -1, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataYear(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getParkDataWithTimeOffset(PERIOD_YEAR, queryTime, 0, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(去年)
	hb, err := getParkDataWithTimeOffset(PERIOD_YEAR, queryTime, -1, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, nil)
	return current, nil
}
