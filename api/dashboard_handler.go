package api

import (
	"eco-service/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
)

//前端dashboard接口

// getQueryTime extracts and parses query time based on period
func getQueryTime(period string, queryTimeStr string) (time.Time, error) {
	var layout string
	switch period {
	case "day":
		layout = "2006-01-02"
	case "month":
		layout = "2006-01"
	case "year":
		layout = "2006"
	default:
		layout = "2006-01-02"
	}
	if queryTimeStr == "" {
		now := time.Now()
		queryTimeStr = now.Format(layout)
	}

	queryTime, err := time.Parse(layout, queryTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid query_time format for period %s", period)
	}

	return queryTime, nil
}

func InitDashboardRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/dashboard/building-power-consumption", BuildingPowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/building-type-power-consumption", BuildingTypePowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/building-floor-power-consumption", BuildingFloorPowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/park-carbon-emission", ParkCarbonEmissionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/park-standard-coal-emission", ParkStandardCoalEmissionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/park-power-consumption", ParkPowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/park-water-consumption", ParkWaterConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/park-water-consumption-range", ParkWaterConsumptionRangeHandler)
}

// @Summary 建筑用电量
// @Description 建筑用电量，根据粒度获取所有建筑的用电量
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/building-power-consumption [get]
func BuildingPowerConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.GetBuildingsPowerConsumption(period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 建筑细分用电量
// @Description 建筑细分用电量，根据粒度,和分类获取所有建筑的用电量。分类为1:照明，2:动力
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Param type query string true "type"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/building-type-power-consumption [get]
func BuildingTypePowerConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	typeStr := r.URL.Query().Get("type")
	if typeStr == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("type is required"))
		return
	}
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("type is invalid"))
		return
	}
	data, err := service.GetBuildingsPowerConsumption(period, queryTime, typeInt)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 建筑楼层用电量
// @Description 建筑楼层用电量，根据粒度获取建筑所有楼层的用电量
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param building_id query string true "building_id"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/building-floor-power-consumption [get]
func BuildingFloorPowerConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	buildingId := r.URL.Query().Get("building_id")
	if buildingId == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("building_id is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.GetBuildingFloorsPowerConsumption(buildingId, period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区碳排放
// @Description 园区碳排放，根据粒度获取园区历史的碳排放,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/park-carbon-emission [get]
func ParkCarbonEmissionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.GetParkCarbonEmissionSubRange(period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区标准煤排放
// @Description 园区标准煤排放，根据粒度获取园区历史的标准煤排放,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/park-standard-coal-emission [get]
func ParkStandardCoalEmissionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.GetParkStandardCoalSubRange(period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区用电量
// @Description 园区用电量，根据粒度获取园区历史的用电量,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param type query string false "type,1:照明,2:动力"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/park-power-consumption [get]
func ParkPowerConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	typeStr := r.URL.Query().Get("type")
	if typeStr == "" {
		typeStr = "0"
	}
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("type is invalid"))
		return
	}
	data, err := service.GetParkPowerConsumptionSubRange(period, queryTime, typeInt)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区用水量
// @Description 园区用水量，根据粒度获取园区当前的用水量,不同粒度返回不同的数据。
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/park-water-consumption [get]
func ParkWaterConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.GetParkWaterConsumption(period, queryTime)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区用水量范围
// @Description 园区用水量范围，根据粒度获取园区历史的用水量范围,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period, hour/day/month/year"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/park-water-consumption-range [get]
func ParkWaterConsumptionRangeHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTime, err := getQueryTime(period, r.URL.Query().Get("query_time"))
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.GetParkWaterConsumptionSubRange(period, queryTime)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}
