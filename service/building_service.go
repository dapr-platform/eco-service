package service

import (
	"context"
	"eco-service/entity"
	"eco-service/model"
	"fmt"
	"time"

	"github.com/dapr-platform/common"
)

/*
同比:
- 天粒度: 与上个月同天的数据比较
- 月粒度: 与上年同月的数据比较
- 年粒度: 无同比,只有环比

环比:
- 天粒度: 与前一天的数据比较
- 月粒度: 与前一月的数据比较
- 年粒度: 与前一年数据比较
*/

var (
	buildingCacheMap = make(map[string]*model.Ecbuilding)
	floorCacheMap    = make(map[string]*model.Ecfloor)
)

type BuildingDataGetter func(time.Time) ([]entity.LabelData, error)
type BuildingRangeDataFetcher func(period string, startTime, endTime time.Time, buildingId string, gatewayType int) ([]entity.LabelData, error)
type FloorRangeDataFetcher func(period string, startTime, endTime time.Time, floorId string, gatewayType int) ([]entity.LabelData, error)

func GetBuildingFloorsPowerConsumption(buildingID string, period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	getters := map[string]BuildingDataGetter{
		PERIOD_DAY: func(t time.Time) ([]entity.LabelData, error) {
			return getBuildingFloorDataDay(buildingID, t, gatewayType)
		},
		PERIOD_MONTH: func(t time.Time) ([]entity.LabelData, error) {
			return getBuildingFloorDataMonth(buildingID, t, gatewayType)
		},
		PERIOD_YEAR: func(t time.Time) ([]entity.LabelData, error) {
			return getBuildingFloorDataYear(buildingID, t, gatewayType)
		},
	}

	if getter, ok := getters[period]; ok {
		return getter(queryTime)
	}

	return nil, fmt.Errorf("unsupported period: %s", period)
}

func GetBuildingsPowerConsumption(period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	getters := map[string]BuildingDataGetter{
		PERIOD_DAY:   func(t time.Time) ([]entity.LabelData, error) { return getBuildingDataDay(t, gatewayType) },
		PERIOD_MONTH: func(t time.Time) ([]entity.LabelData, error) { return getBuildingDataMonth(t, gatewayType) },
		PERIOD_YEAR:  func(t time.Time) ([]entity.LabelData, error) { return getBuildingDataYear(t, gatewayType) },
	}

	if getter, ok := getters[period]; ok {
		return getter(queryTime)
	}

	return nil, fmt.Errorf("unsupported period: %s", period)
}

func GetBuildingPowerConsumptionWithTimeRange(period string, queryTime time.Time,  buildingID string, gatewayType int) ([]entity.LabelData, error) {

	var startTime, endTime time.Time
	switch period {
	case PERIOD_DAY:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 0, 1)
		return getBuildingRangeData(PERIOD_HOUR, startTime, endTime, buildingID, gatewayType, getBuildingDataWithTimeRange)
	case PERIOD_MONTH:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 1, 0)
		return getBuildingRangeData(PERIOD_DAY, startTime, endTime, buildingID, gatewayType, getBuildingDataWithTimeRange)
	case PERIOD_YEAR:
		startTime = time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(1, 0, 0)
		return getBuildingRangeData(PERIOD_MONTH, startTime, endTime, buildingID, gatewayType, getBuildingDataWithTimeRange)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

func GetFloorPowerConsumptionWithTimeRange(period string, queryTime time.Time,  floorID string, gatewayType int) ([]entity.LabelData, error) {

	var startTime, endTime time.Time
	switch period {
	case PERIOD_DAY:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), queryTime.Day(), 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 0, 1)
		return getFloorRangeData(PERIOD_HOUR, startTime, endTime, floorID, gatewayType, getFloorDataWithTimeRange)
	case PERIOD_MONTH:
		startTime = time.Date(queryTime.Year(), queryTime.Month(), 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(0, 1, 0)
		return getFloorRangeData(PERIOD_DAY, startTime, endTime, floorID, gatewayType, getFloorDataWithTimeRange)
	case PERIOD_YEAR:
		startTime = time.Date(queryTime.Year(), 1, 1, 0, 0, 0, 0, queryTime.Location())
		endTime = startTime.AddDate(1, 0, 0)
		return getFloorRangeData(PERIOD_MONTH, startTime, endTime, floorID, gatewayType, getFloorDataWithTimeRange)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}


// 通用的时间范围数据获取函数
func getBuildingRangeData(period string, startTime, endTime time.Time, buildingId string, gatewayType int, fetcher BuildingRangeDataFetcher) ([]entity.LabelData, error) {
	common.Logger.Debugf("Getting range data: period=%s, startTime=%v, endTime=%v, gatewayType=%d", period, startTime, endTime, gatewayType)
	current, err := fetcher(period, startTime, endTime, buildingId, gatewayType)
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
		hb, err = fetcher(period, hbStartTime, hbEndTime, buildingId, gatewayType)
		if err != nil {
			return nil, err
		}
	}
	if getTb {
		common.Logger.Debugf("Getting TB data for range: startTime=%v, endTime=%v", tbStartTime, tbEndTime)
		tb, err = fetcher(period, tbStartTime, tbEndTime, buildingId, gatewayType)
		if err != nil {
			return nil, err
		}
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

// 通用的时间范围数据获取函数
func getFloorRangeData(period string, startTime, endTime time.Time, floorId string, gatewayType int, fetcher FloorRangeDataFetcher) ([]entity.LabelData, error) {
	common.Logger.Debugf("Getting range data: period=%s, startTime=%v, endTime=%v, gatewayType=%d", period, startTime, endTime, gatewayType)
	current, err := fetcher(period, startTime, endTime, floorId, gatewayType)
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
		hb, err = fetcher(period, hbStartTime, hbEndTime, floorId, gatewayType)
		if err != nil {
			return nil, err
		}
	}
	if getTb {
		common.Logger.Debugf("Getting TB data for range: startTime=%v, endTime=%v", tbStartTime, tbEndTime)
		tb, err = fetcher(period, tbStartTime, tbEndTime, floorId, gatewayType)
		if err != nil {
			return nil, err
		}
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingDataWithTimeRange(period string, startTime time.Time, endTime time.Time, buildingId string, gatewayType int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	common.Logger.Debugf("Getting park data with time range: period=%s, startTime=%v, endTime=%v", period, startTime, endTime)

	whereClause := ""
	if gatewayType > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType)
	}
	if buildingId != "" {
		whereClause += fmt.Sprintf("&building_id=%s", buildingId)
	}

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_building_1hTableInfo.Name
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T15:00:00"), endTime.Format("2006-01-02T15:00:00"), whereClause)
		data, err = common.DbQuery[model.Eco_building_1h](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	case PERIOD_DAY:
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause)
		tableName = model.Eco_building_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	case PERIOD_MONTH:
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause)
		tableName = model.Eco_building_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	case PERIOD_YEAR:
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause)
		tableName = model.Eco_building_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		common.Logger.Debugf("Query error: %v", err)
		return nil, err
	}

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
	calcTimeFormat := "2006-01-02T15:04:05"
	common.Logger.Debugf("getParkDataWithTimeRange data: %+v", data)
	keys := []string{}
	keyMap := make(map[string]bool)
	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_building_1h) {
			key := fmt.Sprintf("%s_%s", v.BuildingID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_building_1d) {
			key := fmt.Sprintf("%s_%s", v.BuildingID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_building_1m) {
			key := fmt.Sprintf("%s_%s", v.BuildingID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_building_1y) {
			key := fmt.Sprintf("%s_%s", v.BuildingID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	}

	var sortedData []keyValue
	for _, k := range keys {
		sortedData = append(sortedData, keyValue{k, parkPowerMap[k]})
	}

	result := fillSortedData(sortedData, period, startTime, endTime, calcTimeFormat, timeFormat)

	return result, nil
}
func getFloorDataWithTimeRange(period string, startTime time.Time, endTime time.Time, floorId string, gatewayType int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	common.Logger.Debugf("Getting park data with time range: period=%s, startTime=%v, endTime=%v", period, startTime, endTime)

	whereClause := ""
	if gatewayType > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType)
	}
	if floorId != "" {
		whereClause += fmt.Sprintf("&floor_id=%s", floorId)
	}

	switch period {
	case PERIOD_HOUR:
		tableName = model.Eco_floor_1hTableInfo.Name
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T15:00:00"), endTime.Format("2006-01-02T15:00:00"), whereClause)
		data, err = common.DbQuery[model.Eco_floor_1h](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	case PERIOD_DAY:
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause)
		tableName = model.Eco_floor_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	case PERIOD_MONTH:
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause)
		tableName = model.Eco_floor_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	case PERIOD_YEAR:
		param := fmt.Sprintf("time=$gte.%s&time=$lt.%s%s", startTime.Format("2006-01-02T00:00:00"), endTime.Format("2006-01-02T00:00:00"), whereClause)
		tableName = model.Eco_floor_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			param,
		)
		common.Logger.Debugf("table=%s, param=%s", tableName, param)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		common.Logger.Debugf("Query error: %v", err)
		return nil, err
	}

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
	calcTimeFormat := "2006-01-02T15:04:05"
	common.Logger.Debugf("getParkDataWithTimeRange data: %+v", data)
	keys := []string{}
	keyMap := make(map[string]bool)
	switch period {
	case PERIOD_HOUR:
		for _, v := range data.([]model.Eco_floor_1h) {
			key := fmt.Sprintf("%s_%s", v.FloorID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_floor_1d) {
			key := fmt.Sprintf("%s_%s", v.FloorID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_floor_1m) {
			key := fmt.Sprintf("%s_%s", v.FloorID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_floor_1y) {
			key := fmt.Sprintf("%s_%s", v.FloorID, time.Time(v.Time).Format(calcTimeFormat))
			if _, ok := keyMap[key]; !ok {
				keys = append(keys, key)
				keyMap[key] = true
			}
			if gatewayType == 0 {
				parkPowerMap[key] += v.PowerConsumption
			} else {
				parkPowerMap[key] = v.PowerConsumption
			}
		}
	}

	var sortedData []keyValue
	for _, k := range keys {
		sortedData = append(sortedData, keyValue{k, parkPowerMap[k]})
	}

	result := fillSortedData(sortedData, period, startTime, endTime, calcTimeFormat, timeFormat)

	return result, nil
}
func getBuildingDataWithTimeOffset(period string, queryTime time.Time, years, months, days int, gatewayType int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days)

	whereClause := ""
	if gatewayType > 0 {
		whereClause = fmt.Sprintf("&type=%d", gatewayType)
	}

	switch period {
	case PERIOD_DAY:
		tableName = model.Eco_building_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_building_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_building_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	// Get all buildings first
	buildings, err := common.DbQuery[model.Ecbuilding](
		context.Background(),
		common.GetDaprClient(),
		model.EcbuildingTableInfo.Name,
		"_order=index",
	)
	if err != nil {
		return nil, err
	}

	// Create map for power consumption data
	buildingPowerMap := make(map[string]float64)

	switch period {
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_building_1d) {
			if gatewayType == 0 {
				buildingPowerMap[v.BuildingID] += v.PowerConsumption
			} else {
				buildingPowerMap[v.BuildingID] = v.PowerConsumption
			}
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_building_1m) {
			if gatewayType == 0 {
				buildingPowerMap[v.BuildingID] += v.PowerConsumption
			} else {
				buildingPowerMap[v.BuildingID] = v.PowerConsumption
			}
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_building_1y) {
			if gatewayType == 0 {
				buildingPowerMap[v.BuildingID] += v.PowerConsumption
			} else {
				buildingPowerMap[v.BuildingID] = v.PowerConsumption
			}
		}
	}

	// Create result with all buildings
	result := make([]entity.LabelData, len(buildings))
	for i, building := range buildings {
		result[i] = entity.LabelData{
			Id:    building.ID,
			Label: building.BuildingName,
			Value: buildingPowerMap[building.ID], // Will be 0 if no data exists
		}
	}

	return result, nil
}

func getBuildingFloorDataWithTimeOffset(buildingID string, period string, queryTime time.Time, years, months, days int, gatewayType int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days)

	whereClause := fmt.Sprintf("&building_id=%s", buildingID)
	if gatewayType > 0 {
		whereClause += fmt.Sprintf("&type=%d", gatewayType)
	}

	switch period {
	case PERIOD_DAY:
		tableName = model.Eco_floor_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_floor_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_floor_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-02T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	// Get all floors for this building
	floors, err := common.DbQuery[model.Ecfloor](
		context.Background(),
		common.GetDaprClient(),
		model.EcfloorTableInfo.Name,
		fmt.Sprintf("building_id=%s&_order=index", buildingID),
	)
	if err != nil {
		return nil, err
	}

	// Create map for power consumption data
	floorPowerMap := make(map[string]float64)

	switch period {
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_floor_1d) {
			floorPowerMap[v.FloorID] = v.PowerConsumption
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_floor_1m) {
			floorPowerMap[v.FloorID] = v.PowerConsumption
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_floor_1y) {
			floorPowerMap[v.FloorID] = v.PowerConsumption
		}
	}

	// Create result with all floors
	result := make([]entity.LabelData, len(floors))
	for i, floor := range floors {
		result[i] = entity.LabelData{
			Id:    floor.ID,
			Label: floor.FloorName,
			Value: floorPowerMap[floor.ID], // Will be 0 if no data exists
		}
	}

	return result, nil
}

func getBuildingFloorDataDay(buildingID string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_DAY, queryTime, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一天)
	hb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_DAY, queryTime, 0, 0, -1, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(上月同天)
	tb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_DAY, queryTime, 0, -1, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingFloorDataMonth(buildingID string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_MONTH, queryTime, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(上月)
	hb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_MONTH, queryTime, 0, -1, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同月)
	tb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_MONTH, queryTime, -1, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingFloorDataYear(buildingID string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_YEAR, queryTime, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(去年)
	hb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_YEAR, queryTime, -1, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, nil)
	return current, nil
}

func getBuildingDataDay(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingDataWithTimeOffset(PERIOD_DAY, queryTime, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一天)
	hb, err := getBuildingDataWithTimeOffset(PERIOD_DAY, queryTime, 0, 0, -1, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(上月同天)
	tb, err := getBuildingDataWithTimeOffset(PERIOD_DAY, queryTime, 0, -1, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingDataMonth(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingDataWithTimeOffset(PERIOD_MONTH, queryTime, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(上月)
	hb, err := getBuildingDataWithTimeOffset(PERIOD_MONTH, queryTime, 0, -1, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同月)
	tb, err := getBuildingDataWithTimeOffset(PERIOD_MONTH, queryTime, -1, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingDataYear(queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingDataWithTimeOffset(PERIOD_YEAR, queryTime, 0, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(去年)
	hb, err := getBuildingDataWithTimeOffset(PERIOD_YEAR, queryTime, -1, 0, 0, gatewayType)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, nil)
	return current, nil
}

func calculateRatios(current, hb, tb []entity.LabelData) {
	for i := range current {
		// 取小数点后2位
		current[i].Value = float64(int(current[i].Value*100)) / 100

		// 计算环比
		if i < len(hb) && hb[i].Value != 0 {
			current[i].HB = float64(int(hb[i].Value*100)) / 100
			current[i].HBratio = float64(int((current[i].Value-hb[i].Value)/hb[i].Value*10000)) / 10000
		}

		// 计算同比
		if tb != nil && i < len(tb) && tb[i].Value != 0 {
			current[i].TB = float64(int(tb[i].Value*100)) / 100
			current[i].TBRatio = float64(int((current[i].Value-tb[i].Value)/tb[i].Value*10000)) / 10000
		}
	}
}
