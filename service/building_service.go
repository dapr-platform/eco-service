package service

import (
	"context"
	"eco-service/entity"
	"eco-service/model"
	"fmt"
	"sort"
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

func GetBuildingFloorsPowerConsumption(buildingID string, period string, queryTime time.Time, gatewayType int) ([]entity.LabelData, error) {
	getters := map[string]BuildingDataGetter{
		PERIOD_DAY:   func(t time.Time) ([]entity.LabelData, error) { return getBuildingFloorDataDay(buildingID, t, gatewayType) },
		PERIOD_MONTH: func(t time.Time) ([]entity.LabelData, error) { return getBuildingFloorDataMonth(buildingID, t, gatewayType) },
		PERIOD_YEAR:  func(t time.Time) ([]entity.LabelData, error) { return getBuildingFloorDataYear(buildingID, t, gatewayType) },
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
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_building_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	result := make([]entity.LabelData, 0)
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

	type buildingData struct {
		labelData entity.LabelData
		index     int32
	}

	tempResult := make([]buildingData, 0)
	for buildingID, powerConsumption := range buildingPowerMap {
		building, err := getBuildingInfo(buildingID)
		if err != nil {
			return nil, err
		}
		tempResult = append(tempResult, buildingData{
			labelData: entity.LabelData{
				Id:    buildingID,
				Label: building.BuildingName,
				Value: powerConsumption,
			},
			index: building.Index,
		})
	}

	// Sort by building index
	sort.Slice(tempResult, func(i, j int) bool {
		return tempResult[i].index < tempResult[j].index
	})

	// Convert back to []entity.LabelData
	for _, v := range tempResult {
		result = append(result, v.labelData)
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
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_floor_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s%s", offsetTime.Format("2006-01-01T00:00:00"), whereClause),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	type floorData struct {
		labelData entity.LabelData
		index     int32
	}

	tempResult := make([]floorData, 0)

	switch period {
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_floor_1d) {
			floor, err := getFloorInfo(v.FloorID)
			if err != nil {
				return nil, err
			}
			tempResult = append(tempResult, floorData{
				labelData: entity.LabelData{
					Id:    v.FloorID,
					Label: floor.FloorName,
					Value: v.PowerConsumption,
				},
				index: floor.Index,
			})
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_floor_1m) {
			floor, err := getFloorInfo(v.FloorID)
			if err != nil {
				return nil, err
			}
			tempResult = append(tempResult, floorData{
				labelData: entity.LabelData{
					Id:    v.FloorID,
					Label: floor.FloorName,
					Value: v.PowerConsumption,
				},
				index: floor.Index,
			})
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_floor_1y) {
			floor, err := getFloorInfo(v.FloorID)
			if err != nil {
				return nil, err
			}
			tempResult = append(tempResult, floorData{
				labelData: entity.LabelData{
					Id:    v.FloorID,
					Label: floor.FloorName,
					Value: v.PowerConsumption,
				},
				index: floor.Index,
			})
		}
	}

	// Sort by floor index
	sort.Slice(tempResult, func(i, j int) bool {
		return tempResult[i].index < tempResult[j].index
	})

	// Convert back to []entity.LabelData
	result := make([]entity.LabelData, 0)
	for _, v := range tempResult {
		result = append(result, v.labelData)
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
			current[i].HB = float64(int((current[i].Value-hb[i].Value)/hb[i].Value*10000)) / 10000
		}

		// 计算同比
		if tb != nil && i < len(tb) && tb[i].Value != 0 {
			current[i].TB = float64(int((current[i].Value-tb[i].Value)/tb[i].Value*10000)) / 10000
		}
	}
}

func getBuildingInfo(buildingID string) (*model.Ecbuilding, error) {
	if building, ok := buildingCacheMap[buildingID]; ok {
		return building, nil
	}

	building, err := common.DbGetOne[model.Ecbuilding](
		context.Background(),
		common.GetDaprClient(),
		model.EcbuildingTableInfo.Name,
		fmt.Sprintf("id=%s", buildingID),
	)
	if err != nil {
		return nil, err
	}

	buildingCacheMap[buildingID] = building
	return building, nil
}

func getFloorInfo(floorID string) (*model.Ecfloor, error) {
	if floor, ok := floorCacheMap[floorID]; ok {
		return floor, nil
	}

	floor, err := common.DbGetOne[model.Ecfloor](
		context.Background(),
		common.GetDaprClient(),
		model.EcfloorTableInfo.Name,
		fmt.Sprintf("id=%s", floorID),
	)
	if err != nil {
		return nil, err
	}

	floorCacheMap[floorID] = floor
	return floor, nil
}
