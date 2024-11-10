package api

import (
	"eco-service/service"
	"net/http"

	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
)

func InitManuCollectRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/manu_collect", ManuCollectHandler)
}

// @Summary Manually collect data
// @Description Manually collect data
// @Tags Manually collect data
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
