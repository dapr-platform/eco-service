package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEcgatewayRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/ecgateway/page", EcgatewayPageListHandler)
	r.Get(common.BASE_CONTEXT+"/ecgateway", EcgatewayListHandler)
	r.Post(common.BASE_CONTEXT+"/ecgateway", UpsertEcgatewayHandler)
	r.Delete(common.BASE_CONTEXT+"/ecgateway/{id}", DeleteEcgatewayHandler)
	r.Post(common.BASE_CONTEXT+"/ecgateway/batch-delete", batchDeleteEcgatewayHandler)
	r.Post(common.BASE_CONTEXT+"/ecgateway/batch-upsert", batchUpsertEcgatewayHandler)
	r.Get(common.BASE_CONTEXT+"/ecgateway/groupby", EcgatewayGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Ecgateway
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecgateway/groupby [get]
func EcgatewayGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_eco_gateway")
}

// @Summary batch update
// @Description batch update
// @Tags Ecgateway
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecgateway/batch-upsert [post]
func batchUpsertEcgatewayHandler(w http.ResponseWriter, r *http.Request) {

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

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.EcgatewayTableInfo.Name, model.Ecgateway_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Ecgateway
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
// @Param floor_no query string false "floor_no"
// @Param building_no query string false "building_no"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Ecgateway}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecgateway/page [get]
func EcgatewayPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Ecgateway](w, r, common.GetDaprClient(), "o_eco_gateway", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Ecgateway
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
// @Param floor_no query string false "floor_no"
// @Param building_no query string false "building_no"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Ecgateway} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecgateway [get]
func EcgatewayListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Ecgateway](w, r, common.GetDaprClient(), "o_eco_gateway", "id")
}

// @Summary save
// @Description save
// @Tags Ecgateway
// @Accept       json
// @Param item body model.Ecgateway true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecgateway} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecgateway [post]
func UpsertEcgatewayHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Ecgateway
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Ecgateway")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Ecgateway)
	}

	err = common.DbUpsert[model.Ecgateway](r.Context(), common.GetDaprClient(), val, model.EcgatewayTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Ecgateway
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecgateway} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecgateway/{id} [delete]
func DeleteEcgatewayHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Ecgateway")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_eco_gateway", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Ecgateway
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecgateway/batch-delete [post]
func batchDeleteEcgatewayHandler(w http.ResponseWriter, r *http.Request) {

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
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Ecgateway")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_eco_gateway", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
