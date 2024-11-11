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

type ParkDataGetter func(time.Time) ([]entity.LabelData, error)

// 通用的数据获取函数类型
type DataFetcher func(period string, queryTime time.Time, years, months, days, hours int, gatewayType ...int) ([]entity.LabelData, error)

// 通用的时间范围数据获取函数类型
type RangeDataFetcher func(period string, startTime, endTime time.Time, gatewayType ...int) ([]entity.LabelData, error)

func GetParkWaterConsumption(period string, queryTime time.Time) ([]entity.LabelData, error) {
	fmt.Printf("GetParkWaterConsumption: period=%s, queryTime=%v\n", period, queryTime)
	return getPeriodData(period, queryTime, getParkWaterDataWithTimeOffset)
}

func GetParkWaterConsumptionSubRange(period string, queryTime time.Time) ([]entity.LabelData, error) {
	fmt.Printf("GetParkWaterConsumptionSubRange: period=%s, queryTime=%v\n", period, queryTime)
	var startTime, endTime time.Time
	gatewayType := 0
	switch period {
	case PERIOD_DAY:
		endTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		startTime = endTime.AddDate(0, 0, -1)
		return getRangeData(PERIOD_HOUR, startTime, endTime, gatewayType, getParkWaterDataWithTimeRange)
	case PERIOD_MONTH:
		endTime = time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		startTime = endTime.AddDate(0, -1, 0)
		return getRangeData(PERIOD_DAY, startTime, endTime, gatewayType, getParkWaterDataWithTimeRange)
	case PERIOD_YEAR:
		endTime = time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		startTime = endTime.AddDate(-1, 0, 0)
		return getRangeData(PERIOD_MONTH, startTime, endTime, gatewayType, getParkWaterDataWithTimeRange)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

func GetParkCarbonEmissionSubRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	fmt.Printf("GetParkCarbonEmissionSubRange: period=%s, queryTime=%v, gatewayType=%d\n", period, queryTime, gatewayType)
	result, err := GetParkPowerConsumptionSubRange(period, queryTime, gatewayType)
	if err != nil {
		return nil, err
	}

	for i := range result {
		result[i].Value *= CARBON_FACTOR
	}
	return result, nil
}

func GetParkStandardCoalSubRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	fmt.Printf("GetParkStandardCoalSubRange: period=%s, queryTime=%v, gatewayType=%d\n", period, queryTime, gatewayType)
	result, err := GetParkPowerConsumptionSubRange(period, queryTime, gatewayType)
	if err != nil {
		return nil, err
	}

	for i := range result {
		result[i].Value *= COAL_FACTOR
	}
	return result, nil
}

func GetParkPowerConsumptionSubRange(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	fmt.Printf("GetParkPowerConsumptionSubRange: period=%s, queryTime=%v, gatewayType=%d\n", period, queryTime, gatewayType)

	var startTime, endTime time.Time
	switch period {
	case PERIOD_DAY:
		endTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		startTime = endTime.AddDate(0, 0, -1)
		return getRangeData(PERIOD_HOUR, startTime, endTime, gatewayType, getParkDataWithTimeRange)
	case PERIOD_MONTH:
		endTime = time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		startTime = endTime.AddDate(0, -1, 0)
		return getRangeData(PERIOD_DAY, startTime, endTime, gatewayType, getParkDataWithTimeRange)
	case PERIOD_YEAR:
		endTime = time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		startTime = endTime.AddDate(-1, 0, 0)
		return getRangeData(PERIOD_MONTH, startTime, endTime, gatewayType, getParkDataWithTimeRange)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

func GetParkPowerConsumption(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	fmt.Printf("GetParkPowerConsumption: period=%s, queryTime=%v, gatewayType=%d\n", period, queryTime, gatewayType)
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
		current, err := getter.current()
		if err != nil {
			return nil, err
		}

		hb, err := getter.hb()
		if err != nil {
			return nil, err
		}

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
		hb, err = fetcher(period, hbStartTime, hbEndTime, gatewayType)
		if err != nil {
			return nil, err
		}
	}
	if getTb {
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
		fmt.Printf("Query error: %v\n", err)
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

	for parkID, powerConsumption := range parkPowerMap {
		result = append(result, entity.LabelData{
			Id:    parkID,
			Label: parkID,
			Value: powerConsumption,
		})
	}

	return result, nil
}

func getParkWaterDataWithTimeOffset(period string, queryTime time.Time, years, months, days, hours int, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days).Add(time.Duration(hours) * time.Hour)

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
		fmt.Printf("Query error: %v\n", err)
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

	for parkID, waterConsumption := range parkWaterMap {
		result = append(result, entity.LabelData{
			Id:    parkID,
			Label: parkID,
			Value: waterConsumption,
		})
	}

	return result, nil
}

func getParkDataWithTimeRange(period string, startTime time.Time, endTime time.Time, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

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
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-02T15:00:00"), endTime.Format("2006-01-02T15:00:00"), whereClause),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T23:59:59"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-31T23:59:59"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-12-31T23:59:59"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
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

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_1h) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_1d) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_1m) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_1y) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	}

	// Convert map to slice and sort by time
	type keyValue struct {
		key   string
		value float64
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

	// Convert sorted data to result
	for _, kv := range sortedData {
		parts := strings.Split(kv.key, "_")
		parkID := parts[0]
		timeStr := parts[1]
		result = append(result, entity.LabelData{
			Id:    parkID,
			Label: timeStr,
			Value: kv.value,
		})
	}

	return result, nil
}

func getParkWaterDataWithTimeRange(period string, startTime time.Time, endTime time.Time, gatewayType ...int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

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
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-02T15:00:00"), endTime.Format("2006-01-02T15:00:00"), whereClause),
		)
	case PERIOD_DAY:
		tableName = model.Eco_park_water_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_park_water_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_park_water_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_park_water_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=$gte.%s&time=$lte.%s%s", startTime.Format("2006-01-01T00:00:00"), endTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
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

	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_park_water_1d) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_park_water_1d) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_park_water_1m) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_park_water_1y) {
			key := fmt.Sprintf("%s_%s", v.ParkID, time.Time(v.Time).Format(timeFormat))
			if len(gatewayType) > 0 && gatewayType[0] > 0 {
				parkPowerMap[key] += v.WaterConsumption
			} else {
				parkPowerMap[key] = v.WaterConsumption
			}
		}
	}

	// Convert map to slice and sort by time
	type keyValue struct {
		key   string
		value float64
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

	// Convert sorted data to result
	for _, kv := range sortedData {
		parts := strings.Split(kv.key, "_")
		parkID := parts[0]
		timeStr := parts[1]
		result = append(result, entity.LabelData{
			Id:    parkID,
			Label: timeStr,
			Value: kv.value,
		})
	}

	return result, nil
}
