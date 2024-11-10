package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEcbuildingRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/ecbuilding/page", EcbuildingPageListHandler)
	r.Get(common.BASE_CONTEXT+"/ecbuilding", EcbuildingListHandler)
	r.Post(common.BASE_CONTEXT+"/ecbuilding", UpsertEcbuildingHandler)
	r.Delete(common.BASE_CONTEXT+"/ecbuilding/{id}", DeleteEcbuildingHandler)
	r.Post(common.BASE_CONTEXT+"/ecbuilding/batch-delete", batchDeleteEcbuildingHandler)
	r.Post(common.BASE_CONTEXT+"/ecbuilding/batch-upsert", batchUpsertEcbuildingHandler)
	r.Get(common.BASE_CONTEXT+"/ecbuilding/groupby", EcbuildingGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Ecbuilding
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding/groupby [get]
func EcbuildingGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_eco_building")
}

// @Summary batch update
// @Description batch update
// @Tags Ecbuilding
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding/batch-upsert [post]
func batchUpsertEcbuildingHandler(w http.ResponseWriter, r *http.Request) {

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

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.EcbuildingTableInfo.Name, model.Ecbuilding_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Ecbuilding
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param building_name query string false "building_name"
// @Param park_id query string false "park_id"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Ecbuilding}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding/page [get]
func EcbuildingPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Ecbuilding](w, r, common.GetDaprClient(), "o_eco_building", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Ecbuilding
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param building_name query string false "building_name"
// @Param park_id query string false "park_id"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Ecbuilding} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding [get]
func EcbuildingListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Ecbuilding](w, r, common.GetDaprClient(), "o_eco_building", "id")
}

// @Summary save
// @Description save
// @Tags Ecbuilding
// @Accept       json
// @Param item body model.Ecbuilding true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecbuilding} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding [post]
func UpsertEcbuildingHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Ecbuilding
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Ecbuilding")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Ecbuilding)
	}

	err = common.DbUpsert[model.Ecbuilding](r.Context(), common.GetDaprClient(), val, model.EcbuildingTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Ecbuilding
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Ecbuilding} "object"
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding/{id} [delete]
func DeleteEcbuildingHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Ecbuilding")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_eco_building", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Ecbuilding
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /ecbuilding/batch-delete [post]
func batchDeleteEcbuildingHandler(w http.ResponseWriter, r *http.Request) {

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
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Ecbuilding")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_eco_building", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
