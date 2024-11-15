package api

import (
	"eco-service/service"
	"net/http"

	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
)

func InitManuCollectRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/manu_collect", ManuCollectHandler)
	r.Get(common.BASE_CONTEXT+"/check_collect_date", CheckCollectPowerHandler)
	r.Get(common.BASE_CONTEXT+"/manu_gen_demo_water_data", ManuGenDemoWaterDataHandler)
	r.Get(common.BASE_CONTEXT+"/debug_get_box_hour_stats", DebugGetBoxHourStatsHandler)
	r.Get(common.BASE_CONTEXT+"/manu_fill_gateway_hour_stats", ManuFillGatewayHourStatsHandler)
	r.Get(common.BASE_CONTEXT+"/manu_fill_park_water_hour_stats", ManuFillParkWaterHourStatsHandler)
	r.Get(common.BASE_CONTEXT+"/force_refresh_continuous_aggregate", ForceRefreshContinuousAggregateHandler)
}

// @Summary Manually refresh continuous aggregate
// @Description Manually refresh continuous aggregate
// @Tags Manually
// @Produce  json
// @Success 200 {object} common.Response "success"
// @Router /force_refresh_continuous_aggregate [get]
func ForceRefreshContinuousAggregateHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		err := service.ForceRefreshContinuousAggregate()
		if err != nil {
			common.Logger.Error("手动刷新连续聚合失败," + err.Error())
		}
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary Manually fill park water hour stats
// @Description Manually fill park water hour stats
// @Tags Manually
// @Produce  json
// @Param month query string true "month"
// @Param value query string true "value"
// @Success 200 {object} common.Response "success"
// @Router /manu_fill_park_water_hour_stats [get]
func ManuFillParkWaterHourStatsHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		err := service.ManuFillParkWaterHourStats(r.URL.Query().Get("month"), r.URL.Query().Get("value"))
		if err != nil {
			common.Logger.Error("手动收集数据失败," + err.Error())
			common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
			return
		}
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary Manually fill gateway hour stats
// @Description Manually fill gateway hour stats
// @Tags Manually
// @Produce  json
// @Param month query string true "month"
// @Param value query string true "value"
// @Success 200 {object} common.Response "success"
// @Router /manu_fill_gateway_hour_stats [get]
func ManuFillGatewayHourStatsHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		err := service.ManuFillGatewayHourStats(r.URL.Query().Get("month"), r.URL.Query().Get("value"))
		if err != nil {
			common.Logger.Error("手动收集数据失败," + err.Error())
			common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
			return
		}
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary Manually collect data
// @Description Manually collect data
// @Tags Manually
// @Produce  json
// @Param mac_addr query string true "mac_addr"
// @Param year query string true "year"
// @Param month query string true "month"
// @Param day query string true "day"
// @Success 200 {object} common.Response "success"
// @Router /debug_get_box_hour_stats [get]
func DebugGetBoxHourStatsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := service.DebugGetBoxHourStats(r.URL.Query().Get("mac_addr"), r.URL.Query().Get("year"), r.URL.Query().Get("month"), r.URL.Query().Get("day"))
	if err != nil {
		common.Logger.Error("手动收集数据失败," + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary Manually collect data
// @Description Manually collect data
// @Tags Manually
// @Produce  json
// @Param start query string false "Start time (2024-01-01)"
// @Param end query string false "End time (2024-01-01)"
// @Success 200 {object} common.Response "success"
// @Router /manu_collect [get]
func ManuCollectHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	go func() {
		err := service.ManuCollectGatewayHourlyStatsByDay(start, end)
		if err != nil {
			common.Logger.Error("手动收集数据失败," + err.Error())
		}
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary Manually generate demo water data
// @Description Manually generate demo water data
// @Tags Manually
// @Produce  json
// @Param start query string false "Start time (2024-01-01)"
// @Param end query string false "End time (2024-01-01)"
// @Success 200 {object} common.Response "success"
// @Router /manu_gen_demo_water_data [get]
func ManuGenDemoWaterDataHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		service.ManuGenDemoWaterData(r.URL.Query().Get("start"), r.URL.Query().Get("end"))
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary 查看采集到的数据时间分布
// @Description 查看采集到的数据时间分布
// @Tags Manually
// @Produce  json
// @Param start query string true "Start time (2024-01-01)"
// @Param end query string true "End time (2024-01-01)"
// @Param tablename query string true "tablename"
// @Success 200 {object} common.Response "success"
// @Router /check_collect_date [get]
func CheckCollectPowerHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	tablename := r.URL.Query().Get("tablename")
	data, err := service.CheckCollectData(start, end, tablename)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}
