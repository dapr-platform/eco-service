package api

import (
	"eco-service/service"
	"net/http"

	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
)

func InitTestRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/test", TestHandler)
}

// @Summary Test
// @Description Test
// @Tags Test
// @Produce  json
// @Success 200 {object} common.Response "success"
// @Router /test [get]
func TestHandler(w http.ResponseWriter, r *http.Request) {
	err := service.CollectGatewayHourlyStats()
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
	} else {
		common.HttpResult(w, common.OK)
	}
}
