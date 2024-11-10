package api

import (
	"eco-service/entity"
	"eco-service/service"
	"net/http"
	"strconv"
	"time"

	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/rand"
)

//前端dashboard接口

func InitDashboardRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/building-power-consumption", BuildingPowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/building-type-power-consumption", BuildingTypePowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/building-floor-power-consumption", BuildingFloorPowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/park-carbon-emission", ParkCarbonEmissionHandler)
	r.Get(common.BASE_CONTEXT+"/park-standard-coal-emission", ParkStandardCoalEmissionHandler)
	r.Get(common.BASE_CONTEXT+"/park-power-consumption", ParkPowerConsumptionHandler)
	r.Get(common.BASE_CONTEXT+"/park-water-consumption", ParkWaterConsumptionHandler)
}

// @Summary 建筑用电量
// @Description 建筑用电量，根据粒度获取所有建筑的用电量
// @Tags Dashboard
// @Param period query string true "period"
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
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)

	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
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
// @Param period query string true "period"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Param type query string true "type"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/building-power-consumption [get]
func BuildingTypePowerConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)

	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
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
// @Param period query string true "period"
// @Param building_id query string true "building_id"
// @Param query_time query string false "query_time,格式2024-01-01,不传则默认当天"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.LabelData} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/building-power-consumption [get]
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
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
		return
	}
	data, err := service.GetBuildingFloorsPowerConsumption(period, buildingId, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区碳排放
// @Description 园区碳排放，根据粒度获取园区历史的碳排放,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period"
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
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
		return
	}
	data, err := service.GetParkCarbonEmissionRange(period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区标准煤排放
// @Description 园区标准煤排放，根据粒度获取园区历史的标准煤排放,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period"
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
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
		return
	}
	data, err := service.GetParkStandardCoalRange(period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区用电量
// @Description 园区用电量，根据粒度获取园区历史的用电量,不同粒度返回不同数量的数据。日粒度，返回24小时数据。月粒度，返回31天数据。年粒度，返回12个月数据。包括同比环比
// @Tags Dashboard
// @Param period query string true "period"
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
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
		return
	}
	data, err := service.GetParkPowerConsumptionRange(period, queryTime, 0)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary 园区用水量
// @Description 园区用水量，根据粒度获取园区当前的用水量,不同粒度返回不同的数据。
// @Tags Dashboard
// @Param period query string true "period"
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
	queryTimeStr := r.URL.Query().Get("query_time")
	if queryTimeStr == "" {
		queryTimeStr = time.Now().Format("2006-01-02")
	}
	queryTime, err := time.Parse("2006-01-02", queryTimeStr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("query_time is invalid"))
		return
	}
	_ = queryTime
	label := ""
	switch period {
	case service.PERIOD_DAY:
		label = queryTime.Format("2006-01-02")
	case service.PERIOD_MONTH:
		label = queryTime.Format("2006-01")
	case service.PERIOD_YEAR:
		label = queryTime.Format("2006")
	}
	demoData := []entity.LabelData{
		{Label: label, Value: rand.Intn(1000)},
	}
	common.HttpResult(w, common.OK.WithData(demoData))
}
