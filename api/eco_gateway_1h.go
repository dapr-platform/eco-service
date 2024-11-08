package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEco_gateway_1hRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/eco-gateway-1h/page", Eco_gateway_1hPageListHandler)
	r.Get(common.BASE_CONTEXT+"/eco-gateway-1h", Eco_gateway_1hListHandler)
	r.Post(common.BASE_CONTEXT+"/eco-gateway-1h", UpsertEco_gateway_1hHandler)
	r.Delete(common.BASE_CONTEXT+"/eco-gateway-1h/{id}", DeleteEco_gateway_1hHandler)
	r.Post(common.BASE_CONTEXT+"/eco-gateway-1h/batch-delete", batchDeleteEco_gateway_1hHandler)
	r.Post(common.BASE_CONTEXT+"/eco-gateway-1h/batch-upsert", batchUpsertEco_gateway_1hHandler)
	r.Get(common.BASE_CONTEXT+"/eco-gateway-1h/groupby", Eco_gateway_1hGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Eco_gateway_1h
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h/groupby [get]
func Eco_gateway_1hGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "f_eco_gateway_1h")
}

// @Summary batch update
// @Description batch update
// @Tags Eco_gateway_1h
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h/batch-upsert [post]
func batchUpsertEco_gateway_1hHandler(w http.ResponseWriter, r *http.Request) {

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

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.Eco_gateway_1hTableInfo.Name, model.Eco_gateway_1h_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Eco_gateway_1h
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param time query string false "time"
// @Param gateway_id query string false "gateway_id"
// @Param floor_id query string false "floor_id"
// @Param building_id query string false "building_id"
// @Param type query string false "type"
// @Param power_consumption query string false "power_consumption"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Eco_gateway_1h}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h/page [get]
func Eco_gateway_1hPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Eco_gateway_1h](w, r, common.GetDaprClient(), "f_eco_gateway_1h", "idtime")

}

// @Summary query objects
// @Description query objects
// @Tags Eco_gateway_1h
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param time query string false "time"
// @Param gateway_id query string false "gateway_id"
// @Param floor_id query string false "floor_id"
// @Param building_id query string false "building_id"
// @Param type query string false "type"
// @Param power_consumption query string false "power_consumption"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Eco_gateway_1h} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h [get]
func Eco_gateway_1hListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Eco_gateway_1h](w, r, common.GetDaprClient(), "f_eco_gateway_1h", "idtime")
}

// @Summary save
// @Description save
// @Tags Eco_gateway_1h
// @Accept       json
// @Param item body model.Eco_gateway_1h true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Eco_gateway_1h} "object"
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h [post]
func UpsertEco_gateway_1hHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Eco_gateway_1h
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Eco_gateway_1h")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Eco_gateway_1h)
	}

	err = common.DbUpsert[model.Eco_gateway_1h](r.Context(), common.GetDaprClient(), val, model.Eco_gateway_1hTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Eco_gateway_1h
// @Param idtime  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Eco_gateway_1h} "object"
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h/{id} [delete]
func DeleteEco_gateway_1hHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Eco_gateway_1h")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "f_eco_gateway_1h", "idtime", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Eco_gateway_1h
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /eco-gateway-1h/batch-delete [post]
func batchDeleteEco_gateway_1hHandler(w http.ResponseWriter, r *http.Request) {

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
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Eco_gateway_1h")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "f_eco_gateway_1h", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
