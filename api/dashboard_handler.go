package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InitDashboardRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/dashboard/power-consumption", PowerConsumptionHandler)
}

// @Summary Get power consumption statistics
// @Description Get power consumption statistics by building, floor and gateway type with year-over-year and month-over-month comparisons
// @Tags Dashboard
// @Param period query string true "Period type (day/month/year)"
// @Produce json
// @Success 200 {object} common.Response{data=map[string]interface{}} "Power consumption statistics"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/power-consumption [get]
func PowerConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("period is required"))
		return
	}

	var tableName string
	switch period {
	case "day":
		tableName = "f_eco_building_1d"
	case "month":
		tableName = "f_eco_building_1m"
	case "year":
		tableName = "f_eco_building_1y"
	default:
		common.HttpResult(w, common.ErrParam.AppendMsg("invalid period"))
		return
	}

	// Query building level statistics with YoY and MoM comparisons
	common.CommonQuery[map[string]interface{}](w, r, common.GetDaprClient(), tableName, "building_id")
}
