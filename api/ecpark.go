package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEcparkRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/ecpark/page", EcparkPageListHandler)
	r.Get(common.BASE_CONTEXT+"/ecpark", EcparkListHandler)
	r.Post(common.BASE_CONTEXT+"/ecpark", UpsertEcparkHandler)
	r.Delete(common.BASE_CONTEXT+"/ecpark/{id}", DeleteEcparkHandler)
	r.Post(common.BASE_CONTEXT+"/ecpark/batch-delete", batchDeleteEcparkHandler)
	r.Post(common.BASE_CONTEXT+"/ecpark/batch-upsert", batchUpsertEcparkHandler)
	r.Get(common.BASE_CONTEXT+"/ecpark/groupby", EcparkGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Ecpark
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecpark/groupby [get]
func EcparkGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_eco_park")
}

// @Summary batch update
// @Description batch update
// @Tags Ecpark
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecpark/batch-upsert [post]
func batchUpsertEcparkHandler(w http.ResponseWriter, r *http.Request) {

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

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.EcparkTableInfo.Name, model.Ecpark_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Ecpark
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param park_name query string false "park_name"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Ecpark}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecpark/page [get]
func EcparkPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Ecpark](w, r, common.GetDaprClient(), "o_eco_park", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Ecpark
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param park_name query string false "park_name"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Ecpark} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecpark [get]
func EcparkListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Ecpark](w, r, common.GetDaprClient(), "o_eco_park", "id")
}

// @Summary save
// @Description save
// @Tags Ecpark
// @Accept       json
// @Param item body model.Ecpark true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecpark} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecpark [post]
func UpsertEcparkHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Ecpark
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Ecpark")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Ecpark)
	}

	err = common.DbUpsert[model.Ecpark](r.Context(), common.GetDaprClient(), val, model.EcparkTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Ecpark
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecpark} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecpark/{id} [delete]
func DeleteEcparkHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Ecpark")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_eco_park", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Ecpark
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecpark/batch-delete [post]
func batchDeleteEcparkHandler(w http.ResponseWriter, r *http.Request) {

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
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Ecpark")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_eco_park", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
