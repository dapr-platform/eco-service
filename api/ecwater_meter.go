package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEcwater_meterRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/ecwater-meter/page", Ecwater_meterPageListHandler)
	r.Get(common.BASE_CONTEXT+"/ecwater-meter", Ecwater_meterListHandler)
	r.Post(common.BASE_CONTEXT+"/ecwater-meter", UpsertEcwater_meterHandler)
	r.Delete(common.BASE_CONTEXT+"/ecwater-meter/{id}", DeleteEcwater_meterHandler)
	r.Post(common.BASE_CONTEXT+"/ecwater-meter/batch-delete", batchDeleteEcwater_meterHandler)
	r.Post(common.BASE_CONTEXT+"/ecwater-meter/batch-upsert", batchUpsertEcwater_meterHandler)
	r.Get(common.BASE_CONTEXT+"/ecwater-meter/groupby", Ecwater_meterGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Ecwater_meter
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter/groupby [get]
func Ecwater_meterGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_eco_water_meter")
}

// @Summary batch update
// @Description batch update
// @Tags Ecwater_meter
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter/batch-upsert [post]
func batchUpsertEcwater_meterHandler(w http.ResponseWriter, r *http.Request) {

	var entities []map[string]any
	err := common.ReadRequestBody(r, &entities)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	if len(entities) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of entities is 0"))
		return
	}

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.Ecwater_meterTableInfo.Name, model.Ecwater_meter_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Ecwater_meter
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param model_name query string false "model_name"
// @Param dev_name query string false "dev_name"
// @Param channel_no query string false "channel_no"
// @Param cm_code query string false "cm_code"
// @Param location query string false "location"
// @Param building_id query string false "building_id"
// @Param park_id query string false "park_id"
// @Param type query string false "type"
// @Param total_value query string false "total_value"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Ecwater_meter}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter/page [get]
func Ecwater_meterPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Ecwater_meter](w, r, common.GetDaprClient(), "o_eco_water_meter", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Ecwater_meter
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param model_name query string false "model_name"
// @Param dev_name query string false "dev_name"
// @Param channel_no query string false "channel_no"
// @Param cm_code query string false "cm_code"
// @Param location query string false "location"
// @Param building_id query string false "building_id"
// @Param park_id query string false "park_id"
// @Param type query string false "type"
// @Param total_value query string false "total_value"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Ecwater_meter} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter [get]
func Ecwater_meterListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Ecwater_meter](w, r, common.GetDaprClient(), "o_eco_water_meter", "id")
}

// @Summary save
// @Description save
// @Tags Ecwater_meter
// @Accept       json
// @Param item body model.Ecwater_meter true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecwater_meter} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter [post]
func UpsertEcwater_meterHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Ecwater_meter
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Ecwater_meter")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Ecwater_meter)
	}

	err = common.DbUpsert[model.Ecwater_meter](r.Context(), common.GetDaprClient(), val, model.Ecwater_meterTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Ecwater_meter
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecwater_meter} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter/{id} [delete]
func DeleteEcwater_meterHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Ecwater_meter")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_eco_water_meter", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Ecwater_meter
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecwater-meter/batch-delete [post]
func batchDeleteEcwater_meterHandler(w http.ResponseWriter, r *http.Request) {

	var ids []string
	err := common.ReadRequestBody(r, &ids)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of ids is 0"))
		return
	}
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Ecwater_meter")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_eco_water_meter", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
