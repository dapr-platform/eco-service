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
	PERIOD_HOUR  = "hour"
	PERIOD_DAY   = "day"
	PERIOD_MONTH = "month"
	PERIOD_YEAR  = "year"

	DATA_TYPE_CURRENT = "current"
	DATA_TYPE_HB      = "hb" // 环比
	DATA_TYPE_TB      = "tb" // 同比
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

func GetBuildingFloorsPowerConsumption(buildingID string, period string, queryTime time.Time) ([]entity.LabelData, error) {
	getters := map[string]BuildingDataGetter{
		PERIOD_DAY:   func(t time.Time) ([]entity.LabelData, error) { return getBuildingFloorDataDay(buildingID, t) },
		PERIOD_MONTH: func(t time.Time) ([]entity.LabelData, error) { return getBuildingFloorDataMonth(buildingID, t) },
		PERIOD_YEAR:  func(t time.Time) ([]entity.LabelData, error) { return getBuildingFloorDataYear(buildingID, t) },
	}

	if getter, ok := getters[period]; ok {
		return getter(queryTime)
	}

	return nil, fmt.Errorf("unsupported period: %s", period)
}

func GetBuildingsPowerConsumption(period string, queryTime time.Time) ([]entity.LabelData, error) {
	getters := map[string]BuildingDataGetter{
		PERIOD_DAY:   getBuildingDataDay,
		PERIOD_MONTH: getBuildingDataMonth,
		PERIOD_YEAR:  getBuildingDataYear,
	}

	if getter, ok := getters[period]; ok {
		return getter(queryTime)
	}

	return nil, fmt.Errorf("unsupported period: %s", period)
}

func getBuildingDataWithTimeOffset(period string, queryTime time.Time, years, months, days int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days)

	switch period {
	case PERIOD_DAY:
		tableName = model.Eco_building_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006-01-02")),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_building_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006-01")),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_building_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_building_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s", offsetTime.Format("2006")),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	result := make([]entity.LabelData, 0)

	switch period {
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_building_1d) {
			building, err := getBuildingInfo(v.BuildingID)
			if err != nil {
				return nil, err
			}
			result = append(result, entity.LabelData{
				Id:    v.BuildingID,
				Label: building.BuildingName,
				Value: v.PowerConsumption,
			})
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_building_1m) {
			building, err := getBuildingInfo(v.BuildingID)
			if err != nil {
				return nil, err
			}
			result = append(result, entity.LabelData{
				Id:    v.BuildingID,
				Label: building.BuildingName,
				Value: v.PowerConsumption,
			})
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_building_1y) {
			building, err := getBuildingInfo(v.BuildingID)
			if err != nil {
				return nil, err
			}
			result = append(result, entity.LabelData{
				Id:    v.BuildingID,
				Label: building.BuildingName,
				Value: v.PowerConsumption,
			})
		}
	}

	return result, nil
}

func getBuildingFloorDataWithTimeOffset(buildingID string, period string, queryTime time.Time, years, months, days int) ([]entity.LabelData, error) {
	var data interface{}
	var err error
	var tableName string

	offsetTime := queryTime.AddDate(years, months, days)

	switch period {
	case PERIOD_DAY:
		tableName = model.Eco_floor_1dTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1d](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s&building_id=%s", offsetTime.Format("2006-01-02"), buildingID),
		)
	case PERIOD_MONTH:
		tableName = model.Eco_floor_1mTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1m](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s&building_id=%s", offsetTime.Format("2006-01"), buildingID),
		)
	case PERIOD_YEAR:
		tableName = model.Eco_floor_1yTableInfo.Name
		data, err = common.DbQuery[model.Eco_floor_1y](
			context.Background(),
			common.GetDaprClient(),
			tableName,
			fmt.Sprintf("time=%s&building_id=%s", offsetTime.Format("2006"), buildingID),
		)
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	if err != nil {
		return nil, err
	}

	result := make([]entity.LabelData, 0)

	switch period {
	case PERIOD_DAY:
		for _, v := range data.([]model.Eco_floor_1d) {
			floor, err := getFloorInfo(v.FloorID)
			if err != nil {
				return nil, err
			}
			result = append(result, entity.LabelData{
				Id:    v.FloorID,
				Label: floor.FloorName,
				Value: v.PowerConsumption,
			})
		}
	case PERIOD_MONTH:
		for _, v := range data.([]model.Eco_floor_1m) {
			floor, err := getFloorInfo(v.FloorID)
			if err != nil {
				return nil, err
			}
			result = append(result, entity.LabelData{
				Id:    v.FloorID,
				Label: floor.FloorName,
				Value: v.PowerConsumption,
			})
		}
	case PERIOD_YEAR:
		for _, v := range data.([]model.Eco_floor_1y) {
			floor, err := getFloorInfo(v.FloorID)
			if err != nil {
				return nil, err
			}
			result = append(result, entity.LabelData{
				Id:    v.FloorID,
				Label: floor.FloorName,
				Value: v.PowerConsumption,
			})
		}
	}

	return result, nil
}

func getBuildingFloorDataDay(buildingID string, queryTime time.Time) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_DAY, queryTime, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一天)
	hb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_DAY, queryTime, 0, 0, -1)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(上月同天)
	tb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_DAY, queryTime, 0, -1, 0)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingFloorDataMonth(buildingID string, queryTime time.Time) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_MONTH, queryTime, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(上月)
	hb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_MONTH, queryTime, 0, -1, 0)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同月)
	tb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_MONTH, queryTime, -1, 0, 0)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingFloorDataYear(buildingID string, queryTime time.Time) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_YEAR, queryTime, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(去年)
	hb, err := getBuildingFloorDataWithTimeOffset(buildingID, PERIOD_YEAR, queryTime, -1, 0, 0)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, nil)
	return current, nil
}

func getBuildingDataDay(queryTime time.Time) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingDataWithTimeOffset(PERIOD_DAY, queryTime, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(前一天)
	hb, err := getBuildingDataWithTimeOffset(PERIOD_DAY, queryTime, 0, 0, -1)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(上月同天)
	tb, err := getBuildingDataWithTimeOffset(PERIOD_DAY, queryTime, 0, -1, 0)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingDataMonth(queryTime time.Time) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingDataWithTimeOffset(PERIOD_MONTH, queryTime, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(上月)
	hb, err := getBuildingDataWithTimeOffset(PERIOD_MONTH, queryTime, 0, -1, 0)
	if err != nil {
		return nil, err
	}

	// 获取同比数据(去年同月)
	tb, err := getBuildingDataWithTimeOffset(PERIOD_MONTH, queryTime, -1, 0, 0)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, tb)
	return current, nil
}

func getBuildingDataYear(queryTime time.Time) ([]entity.LabelData, error) {
	// 获取当前数据
	current, err := getBuildingDataWithTimeOffset(PERIOD_YEAR, queryTime, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取环比数据(去年)
	hb, err := getBuildingDataWithTimeOffset(PERIOD_YEAR, queryTime, -1, 0, 0)
	if err != nil {
		return nil, err
	}

	calculateRatios(current, hb, nil)
	return current, nil
}

func calculateRatios(current, hb, tb []entity.LabelData) {
	for i := range current {
		// 计算环比
		if i < len(hb) && hb[i].Value != 0 {
			current[i].HB = (current[i].Value - hb[i].Value) / hb[i].Value
		}

		// 计算同比
		if tb != nil && i < len(tb) && tb[i].Value != 0 {
			current[i].TB = (current[i].Value - tb[i].Value) / tb[i].Value
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
