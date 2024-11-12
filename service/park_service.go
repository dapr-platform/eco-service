package service

import (
	"context"
	"eco-service/entity"
	"eco-service/model"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dapr-platform/common"
)

const (
	CARBON_FACTOR = 0.61   // 碳排放系数
	COAL_FACTOR   = 0.1229 // 标准煤系数
)

type keyValue struct {
	key   string
	value float64
}

type ParkDataGetter func(time.Time) ([]entity.LabelData, error)

// 通用的数据获取函数类型
type DataFetcher func(period string, queryTime time.Time, years, months, days, hours int, gatewayType ...int) ([]entity.LabelData, error)

// 通用的时间范围数据获取函数类型
type RangeDataFetcher func(period string, startTime, endTime time.Time, gatewayType ...int) ([]entity.LabelData, error)

func GetParkWaterConsumption(period string, queryTime time.Time) ([]entity.LabelData, error) {
	common.Logger.Debugf("GetParkWaterConsumption: period=%s, queryTime=%v", period, queryTime)
	return getPeriodData(period, queryTime, getParkWaterDataWithTimeOffset)
}

func GetParkWaterConsumptionSubRange(period string, queryTime time.Time) ([]entity.LabelData, error) {
	common.Logger.Debugf("GetParkWaterConsumptionSubRange: period=%s, queryTime=%v", period, queryTime)
	var startTime, endTime time.Time
	gatewayType := 0
	switch period {
	case PERIOD_DAY:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 0, 1)
		return getRangeData(PERIOD_HOUR, startTime, endTime, gatewayType, getParkWaterDataWithTimeRange)
	case PERIOD_MONTH:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 1, 0)
		return getRangeData(PERIOD_DAY, startTime, endTime, gatewayType, getParkWaterDataWithTimeRange)
	case PERIOD_YEAR:
		startTime = time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(1, 0, 0)
		return getRangeData(PERIOD_MONTH, startTime, endTime, gatewayType, getParkWaterDataWithTimeRange)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

func GetParkCarbonEmissionSubRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	common.Logger.Debugf("GetParkCarbonEmissionSubRange: period=%s, queryTime=%v, gatewayType=%d", period, queryTime, gatewayType)
	result, err := GetParkPowerConsumptionSubRange(period, queryTime, gatewayType)
	if err != nil {
		return nil, err
	}

	for i := range result {
		result[i].Value = float64(int(result[i].Value*CARBON_FACTOR*100)) / 100
		result[i].HB = float64(int(result[i].HB*CARBON_FACTOR*100)) / 100
		result[i].TB = float64(int(result[i].TB*CARBON_FACTOR*100)) / 100
	}
	return result, nil
}

func GetParkStandardCoalSubRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	common.Logger.Debugf("GetParkStandardCoalSubRange: period=%s, queryTime=%v, gatewayType=%d", period, queryTime, gatewayType)
	result, err := GetParkPowerConsumptionSubRange(period, queryTime, gatewayType)
	if err != nil {
		return nil, err
	}

	for i := range result {
		result[i].Value = float64(int(result[i].Value*COAL_FACTOR*100)) / 100
		result[i].HB = float64(int(result[i].HB*COAL_FACTOR*100)) / 100
		result[i].TB = float64(int(result[i].TB*COAL_FACTOR*100)) / 100

	}
	return result, nil
}

func GetParkPowerConsumptionSubRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	common.Logger.Debugf("GetParkPowerConsumptionSubRange: period=%s, queryTime=%v, gatewayType=%d", period, queryTime, gatewayType)

	var startTime, endTime time.Time
	switch period {
	case PERIOD_DAY:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 0, 1)
		common.Logger.Debugf("PERIOD_DAY: startTime=%v, endTime=%v", startTime, endTime)
		return getRangeData(PERIOD_HOUR, startTime, endTime, gatewayType, getParkDataWithTimeRange)
	case PERIOD_MONTH:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 1, 0)
		common.Logger.Debugf("PERIOD_MONTH: startTime=%v, endTime=%v", startTime, endTime)
		return getRangeData(PERIOD_DAY, startTime, endTime, gatewayType, getParkDataWithTimeRange)
	case PERIOD_YEAR:
		startTime = time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(1, 0, 0)
		common.Logger.Debugf("PERIOD_YEAR: startTime=%v, endTime=%v", startTime, endTime)
		return getRangeData(PERIOD_MONTH, startTime, endTime, gatewayType, getParkDataWithTimeRange)
	default:
		common.Logger.Debugf("Unsupported period: %s", period)
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

func GetParkPowerConsumption(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	common.Logger.Debugf("GetParkPowerConsumption: period=%s, queryTime=%v, gatewayType=%d", period, queryTime, gatewayType)
	return getPeriodData(period, queryTime, getParkDataWithTimeOffset)
}

// 通用的周期数据获取函数
func getPeriodData(period string, queryTime time.Time, fetcher DataFetcher) ([]entity.LabelData, error) {
	getters := map[string]struct {
		current func() ([]entity.LabelData, error)
		hb      func() ([]entity.LabelData, error)
		tb      func() ([]entity.LabelData, error)
	}{
		PERIOD_HOUR: {
			current: func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, 0, 0) },
			hb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, 0, -1) },
			tb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, -1, 0) },
		},
		PERIOD_DAY: {
			current: func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, 0, 0) },
			hb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, -1, 0) },
			tb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, -1, 0, 0) },
		},
		PERIOD_MONTH: {
			current: func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, 0, 0) },
			hb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, -1, 0, 0) },
			tb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, -1, 0, 0, 0) },
		},
		PERIOD_YEAR: {
			current: func() ([]entity.LabelData, error) { return fetcher(period, queryTime, 0, 0, 0, 0) },
			hb:      func() ([]entity.LabelData, error) { return fetcher(period, queryTime, -1, 0, 0, 0) },
			tb:      func() ([]entity.LabelData, error) { return nil, nil },
		},
	}

	if getter, ok := getters[period]; ok {
		common.Logger.Debugf("Getting current period data for period: %s", period)
		current, err := getter.current()
		if err != nil {
			return nil, err
		}

		common.Logger.Debugf("Getting HB period data for period: %s", period)
		hb, err := getter.hb()
		if err != nil {
			return nil, err
		}

		common.Logger.Debugf("Getting TB period data for period: %s", period)
		tb, err := getter.tb()
		if err != nil {
			return nil, err
		}

		calculateRatios(current, hb, tb)
		return current, nil
	}

	return nil, fmt.Errorf("unsupported period: %s", period)
}

// 通用的时间范围数据获取函数
func getRangeData(period string, startTime, endTime time.Time, gatewayType int, fetcher RangeDataFetcher) ([]entity.LabelData, error) {
	common.Logger.Debugf("Getting range data: period=%s, startTime=%v, endTime=%v, gatewayType=%d", period, startTime, endTime, gatewayType)
	current, err := fetcher(period, startTime, endTime, gatewayType)
	if err != nil {
		return nil, err
	}

	var hbStartTime, hbEndTime, tbStartTime, tbEndTime time.Time
	var getHb, getTb bool

	switch period {
	case PERIOD_HOUR:
		hbStartTime = startTime.Add(-1 * time.Hour)
		hbEndTime = endTime.Add(-1 * time.Hour)
		tbStartTime = startTime.AddDate(0, 0, -1)
		tbEndTime = endTime.AddDate(0, 0, -1)
		getHb, getTb = true, true
	case PERIOD_DAY:
		hbStartTime = startTime.AddDate(0, 0, -1)
		hbEndTime = endTime.AddDate(0, 0, -1)
		tbStartTime = startTime.AddDate(0, -1, 0)
		tbEndTime = endTime.AddDate(0, -1, 0)
		getHb, getTb = true, true
	case PERIOD_MONTH:
		hbStartTime = startTime.AddDate(0, -1, 0)
		hbEndTime = endTime.AddDate(0, -1, 0)
		tbStartTime = startTime.AddDate(-1, 0, 0)
		tbEndTime = endTime.AddDate(-1, 0, 0)
		getHb, getTb = true, true
	case PERIOD_YEAR:
		tbStartTime = startTime.AddDate(-1, 0, 0)
		tbEndTime = endTime.AddDate(-1, 0, 0)
		getHb, getTb = false, true
	}

	var hb, tb []entity.LabelData
	if getHb {
		common.Logger.Debugf("Getting HB data for range: startTime=%v, endTime=%v", hbStartTime, hbEndTime)
		hb, err = fetcher(period, hbStartTime, hbEndTime, gatewayType)
		if err != nil {
			return nil, err
		}
	}
	if getTb {
		common.Logger.Debugf("Getting TB data for range: startTime=%v, endTime=%v", tbStartTime, tbEndTime)
		tb, err = fetcher(period, tbStartTime, tbEndTime, gatewayType)
		if err != nil {
			return nil, err
		}
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getParkDataWithTimeOffset(period string, queryTime time.Time, years, months, days, hours int, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days).Add(time.Duration(hours) * time.Hour)
	common.Logger.Debugf("Getting park data with time offset: period=%s, queryTime=%v, offsetTime=%v", period, queryTime, offsetTime)

	whereClause := ""
	if len(gatewayType) > 0 && gatewayType[0] > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType[0])
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
		common.Logger.Debugf("Query error: %v", err)
		return nil, err
	}

	result := make([]entity.LabelData, 0)
	parkPowerMap := make(map[string]float64)

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_1h) {
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_1d) {
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_1m) {
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_1y) {
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[v.ParkID] += v.PowerConsumption
			} else {
				parkPowerMap[v.ParkID] = v.PowerConsumption
			}
		}
	}

	// Get park info
	park, err := getParkInfo()
	if err != nil {
		return nil, err
	}

	// Return result with park ID even if no data
	result = append(result, entity.LabelData{
		Id:    park.ID,
		Label: park.ID,
		Value: parkPowerMap[park.ID],
	})

	return result, nil
}

func getParkWaterDataWithTimeOffset(period string, queryTime time.Time, years, months, days, hours int, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days).Add(time.Duration(hours) * time.Hour)
	common.Logger.Debugf("Getting water data with time offset: period=%s, queryTime=%v, offsetTime=%v", period, queryTime, offsetTime)

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_park_water_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006-01-02 15")),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_water_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006-01-02")),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_water_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006-01")),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_water_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006")),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		common.Logger.Debugf("Query error: %v", err)
		return nil, err
	}

	result := make([]entity.LabelData, 0)
	parkWaterMap := make(map[string]float64)

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_water_1d) {
			parkWaterMap[v.ParkID] = v.WaterConsumption
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_water_1d) {
			parkWaterMap[v.ParkID] = v.WaterConsumption
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_water_1m) {
			parkWaterMap[v.ParkID] = v.WaterConsumption
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_water_1y) {
			parkWaterMap[v.ParkID] = v.WaterConsumption
		}
	}

	// Get park info
	park, err := getParkInfo()
	if err != nil {
		return nil, err
	}

	// Return result with park ID even if no data
	result = append(result, entity.LabelData{
		Id:    park.ID,
		Label: park.ID,
		Value: parkWaterMap[park.ID],
	})

	return result, nil
}

func getParkDataWithTimeRange(period string, startTime time.Time, endTime time.Time, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	common.Logger.Debugf("Getting park data with time range: period=%s, startTime=%v, endTime=%v", period, startTime, endTime)

	whereClause := ""
	if len(gatewayType) > 0 && gatewayType[0] > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType[0])
	}

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_park_1hTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1h](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T15:00:00"), endTime.Format("2006-01-02T15:00:00"), whereClause),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		common.Logger.Debugf("Query error: %v", err)
		return nil, err
	}

	result := make([]entity.LabelData, 0)
	parkPowerMap := make(map[string]float64)
	var timeFormat string

	switch period {
	case PERIOD_HOUR:
		timeFormat = "15"
	case PERIOD_DAY:
		timeFormat = "02"
	case PERIOD_MONTH:
		timeFormat = "01"
	case PERIOD_YEAR:
		timeFormat = "2006"
	}
	calcTimeFormat := "2006-01-02_15:04:05"

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_1h) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_1d) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_1m) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_1y) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	}

	var sortedData []keyValue
	for k, v := range parkPowerMap {
		sortedData = append(sortedData, keyValue{k, v})
	}

	// Sort by time (second part of key after "_")
	sort.Slice(sortedData, func(i, j int) bool {
		time1 := strings.Split(sortedData[i].key, "_")[1]
		time2 := strings.Split(sortedData[j].key, "_")[1]
		return time1 < time2
	})
	result = fillSortedData(sortedData, period, startTime, endTime, calcTimeFormat, timeFormat)

	return result, nil
}

func getParkWaterDataWithTimeRange(period string, startTime time.Time, endTime time.Time, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	common.Logger.Debugf("Getting water data with time range: period=%s, startTime=%v, endTime=%v", period, startTime, endTime)

	whereClause := ""
	if len(gatewayType) > 0 && gatewayType[0] > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType[0])
	}

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_park_water_1hTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1h](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T15:00:00"), endTime.Format("2006-01-02T15:00:00"), whereClause),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_water_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_water_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_water_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		common.Logger.Debugf("Query error: %v", err)
		fmt.Printf("Query error: %v\n", err)
		return nil, err
	}

	result := make([]entity.LabelData, 0)
	parkPowerMap := make(map[string]float64)
	var timeFormat string

	switch period {
	case PERIOD_HOUR:
		timeFormat = "15:04"
	case PERIOD_DAY:
		timeFormat = "01-02"
	case PERIOD_MONTH:
		timeFormat = "2006-01"
	case PERIOD_YEAR:
		timeFormat = "2006"
	}
	calcTimeFormat := "2006-01-02_15:04:05"
	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_water_1d) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_water_1d) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_water_1m) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_water_1y) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(calcTimeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	}

	// Convert map to slice and sort by time

	var sortedData []keyValue
	for k, v := range parkPowerMap {
		sortedData = append(sortedData, keyValue{k, v})
	}

	// Sort by time (second part of key after "_")
	sort.Slice(sortedData, func(i, j int) bool {
		time1 := strings.Split(sortedData[i].key, "_")[1]
		time2 := strings.Split(sortedData[j].key, "_")[1]
		return time1 < time2
	})
	result = fillSortedData(sortedData, period, startTime, endTime, calcTimeFormat, timeFormat)

	return result, nil
}

func getParkInfo() (*model.Ecpark, error) {
	park, err := common.DbGetOne[model.Ecpark](
		context.Background(),
		common.GetDaprClient(),
		model.EcparkTableInfo.Name,
		"",
	)
	if err != nil {
		return nil, err
	}
	return park, nil
}

// fillSortedData fills in missing time points in sorted data with zero values
func fillSortedData(sortedData []keyValue, period string, startTime time.Time, endTime time.Time, calcTimeFormat string, timeFormat string) []entity.LabelData {
	// 创建一个map来去重
	uniqueLabels := make(map[string]entity.LabelData)

	// 创建map用于快速查找值
	valueMap := make(map[string]float64)
	var parkID string
	if len(sortedData) > 0 {
		for _, kv := range sortedData {
			parts := strings.Split(kv.key, "_")
			if parkID == "" {
				parkID = parts[0]
			}
			timeStr := parts[1]
			valueMap[timeStr] = kv.value
		}
	}

	// 如果没有找到parkID，使用默认值
	if parkID == "" {
		parkID = "default"
	}

	// 生成连续的时间点
	var step time.Duration
	switch period {
	case PERIOD_HOUR:
		step = time.Hour
	case PERIOD_DAY:
		step = 24 * time.Hour
	case PERIOD_MONTH:
		step = 30 * 24 * time.Hour
	case PERIOD_YEAR:
		step = 365 * 24 * time.Hour
	}

	// 使用map来确保label唯一性
	for t := startTime; !t.After(endTime); t = t.Add(step) {
		timeStr := t.Format(calcTimeFormat)
		label := t.Format(timeFormat)
		value := valueMap[timeStr]

		uniqueLabels[label] = entity.LabelData{
			Id:    parkID,
			Label: label,
			Value: value,
		}
	}

	// 转换为有序切片
	result := make([]entity.LabelData, 0, len(uniqueLabels))
	for _, data := range uniqueLabels {
		result = append(result, data)
	}

	// 按label排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result
}
