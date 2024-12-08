package api

import (
	"eco-service/client"
	"eco-service/service"
	"net/http"

	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cast"
)

func InitManuCollectRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/manu_collect", ManuCollectHandler)
	r.Get(common.BASE_CONTEXT+"/check_collect_date", CheckCollectPowerHandler)
	r.Get(common.BASE_CONTEXT+"/manu_collect_water_data", ManuCollectWaterDataHandler)
	r.Get(common.BASE_CONTEXT+"/debug_get_box_hour_stats", DebugGetBoxHourStatsHandler)
	r.Get(common.BASE_CONTEXT+"/debug_get_month_stats", DebugGetMonthStatsHandler)
	r.Post(common.BASE_CONTEXT+"/debug_method_invoke", DebugMethodInvokeHandler)
	r.Get(common.BASE_CONTEXT+"/manu_fill_gateway_hour_stats", ManuFillGatewayHourStatsHandler)
	r.Get(common.BASE_CONTEXT+"/manu_fill_park_water_hour_stats", ManuFillParkWaterHourStatsHandler)
	r.Get(common.BASE_CONTEXT+"/manu_fill_power_collect_iot_data", ManuFillPowerCollectIotDataHandler)
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
// @Param cm_code query string true "cm_code"
// @Param start query string true "start"
// @Param end query string true "end"
// @Param value query string true "value"
// @Success 200 {object} common.Response "success"
// @Router /manu_fill_park_water_hour_stats [get]
func ManuFillParkWaterHourStatsHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		err := service.ManuFillParkWaterHourStats(r.URL.Query().Get("cm_code"), r.URL.Query().Get("start"), r.URL.Query().Get("end"), r.URL.Query().Get("value"))
		if err != nil {
			common.Logger.Error("手动收集数据失败," + err.Error())
			common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
			return
		}
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}


// @Summary Manually fill power collect iot data
// @Description Manually fill power collect iot data
// @Tags Manually
// @Produce  json
// @Param cm_code query string true "cm_code"
// @Param start query string true "start"
// @Param end query string true "end"
// @Param value query string true "value"
// @Success 200 {object} common.Response "success"
// @Router /manu_fill_power_collect_iot_data [get]
func ManuFillPowerCollectIotDataHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		err := service.ManuFillPowerCollectIotData(r.URL.Query().Get("cm_code"), r.URL.Query().Get("start"), r.URL.Query().Get("end"), r.URL.Query().Get("value"))
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

// @Summary Manually invoke method
// @Description Manually invoke method
// @Tags Manually
// @Produce  json
// @Param method query string true "method"
// @Param params body string true "params"
// @Success 200 {object} common.Response "success"
// @Router /debug_method_invoke [post]
func DebugMethodInvokeHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	err := common.ReadRequestBody(r, &body)
	if err != nil {
		common.Logger.Error("解析参数失败," + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	data, err := client.GetFunc(r.URL.Query().Get("method"), body)
	if err != nil {
		common.Logger.Error("手动收集数据失败," + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}

// @Summary Manually collect data
// @Description Manually collect data for month
// @Tags Manually
// @Produce  json
// @Param mac_addr query string true "mac_addr"
// @Param year query string true "year"
// @Param month query string true "month"
// @Success 200 {object} common.Response "success"
// @Router /debug_get_month_stats [get]
func DebugGetMonthStatsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := service.DebugGetBoxMonthStats(r.URL.Query().Get("mac_addr"), r.URL.Query().Get("year"), r.URL.Query().Get("month"))
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
// @Param mac_addr query string false "mac_addr"
// @Success 200 {object} common.Response "success"
// @Router /manu_collect [get]
func ManuCollectHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	macAddr := r.URL.Query().Get("mac_addr")
	go func() {
		err := service.ManuCollectGatewayHourlyStatsByDay(start, end, macAddr)
		if err != nil {
			common.Logger.Error("手动收集数据失败," + err.Error())
		}
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary Manually collect water data
// @Description Manually collect water data
// @Tags Manually
// @Produce  json
// @Success 200 {object} common.Response "success"
// @Router /manu_collect_water_data [get]
func ManuCollectWaterDataHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		service.CollectWaterMeterRealData()
	}()
	common.HttpResult(w, common.OK.WithData("后台运行，请查看日志"))
}

// @Summary 查看采集到的数据时间分布
// @Description 查看采集到的数据时间分布
// @Tags Manually
// @Produce  json
// @Param start query string true "Start time (2024-01-01)"
// @Param end query string true "End time (2024-01-01)"
// @Param collect_type query string true "collect_type,0:gateway,1:park_water"
// @Success 200 {object} common.Response "success"
// @Router /check_collect_date [get]
func CheckCollectPowerHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	collectType := r.URL.Query().Get("collect_type")
	data, err := service.CheckCollectData(start, end, cast.ToInt(collectType))
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))
}
