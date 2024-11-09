package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEcfloorRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/ecfloor/page", EcfloorPageListHandler)
	r.Get(common.BASE_CONTEXT+"/ecfloor", EcfloorListHandler)
	r.Post(common.BASE_CONTEXT+"/ecfloor", UpsertEcfloorHandler)
	r.Delete(common.BASE_CONTEXT+"/ecfloor/{id}", DeleteEcfloorHandler)
	r.Post(common.BASE_CONTEXT+"/ecfloor/batch-delete", batchDeleteEcfloorHandler)
	r.Post(common.BASE_CONTEXT+"/ecfloor/batch-upsert", batchUpsertEcfloorHandler)
	r.Get(common.BASE_CONTEXT+"/ecfloor/groupby", EcfloorGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Ecfloor
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecfloor/groupby [get]
func EcfloorGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_eco_floor")
}

// @Summary batch update
// @Description batch update
// @Tags Ecfloor
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecfloor/batch-upsert [post]
func batchUpsertEcfloorHandler(w http.ResponseWriter, r *http.Request) {

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

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.EcfloorTableInfo.Name, model.Ecfloor_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Ecfloor
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param floor_name query string false "floor_name"
// @Param building_id query string false "building_id"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Ecfloor}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecfloor/page [get]
func EcfloorPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Ecfloor](w, r, common.GetDaprClient(), "o_eco_floor", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Ecfloor
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param floor_name query string false "floor_name"
// @Param building_id query string false "building_id"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Ecfloor} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecfloor [get]
func EcfloorListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Ecfloor](w, r, common.GetDaprClient(), "o_eco_floor", "id")
}

// @Summary save
// @Description save
// @Tags Ecfloor
// @Accept       json
// @Param item body model.Ecfloor true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecfloor} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecfloor [post]
func UpsertEcfloorHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Ecfloor
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Ecfloor")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Ecfloor)
	}

	err = common.DbUpsert[model.Ecfloor](r.Context(), common.GetDaprClient(), val, model.EcfloorTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Ecfloor
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecfloor} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecfloor/{id} [delete]
func DeleteEcfloorHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Ecfloor")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_eco_floor", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Ecfloor
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecfloor/batch-delete [post]
func batchDeleteEcfloorHandler(w http.ResponseWriter, r *http.Request) {

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
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Ecfloor")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_eco_floor", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
