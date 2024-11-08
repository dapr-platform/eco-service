package api

import (
	"eco-service/model"
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func InitEco_floor_1yRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/eco-floor-1y/page", Eco_floor_1yPageListHandler)
	r.Get(common.BASE_CONTEXT+"/eco-floor-1y", Eco_floor_1yListHandler)
	r.Post(common.BASE_CONTEXT+"/eco-floor-1y", UpsertEco_floor_1yHandler)
	r.Delete(common.BASE_CONTEXT+"/eco-floor-1y/{id}", DeleteEco_floor_1yHandler)
	r.Post(common.BASE_CONTEXT+"/eco-floor-1y/batch-delete", batchDeleteEco_floor_1yHandler)
	r.Post(common.BASE_CONTEXT+"/eco-floor-1y/batch-upsert", batchUpsertEco_floor_1yHandler)
	r.Get(common.BASE_CONTEXT+"/eco-floor-1y/groupby", Eco_floor_1yGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Eco_floor_1y
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y/groupby [get]
func Eco_floor_1yGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "f_eco_floor_1y")
}

// @Summary batch update
// @Description batch update
// @Tags Eco_floor_1y
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y/batch-upsert [post]
func batchUpsertEco_floor_1yHandler(w http.ResponseWriter, r *http.Request) {

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

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.Eco_floor_1yTableInfo.Name, model.Eco_floor_1y_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Eco_floor_1y
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param time query string false "time"
// @Param floor_id query string false "floor_id"
// @Param building_id query string false "building_id"
// @Param type query string false "type"
// @Param power_consumption query string false "power_consumption"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Eco_floor_1y}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y/page [get]
func Eco_floor_1yPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Eco_floor_1y](w, r, common.GetDaprClient(), "f_eco_floor_1y", "time")

}

// @Summary query objects
// @Description query objects
// @Tags Eco_floor_1y
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param time query string false "time"
// @Param floor_id query string false "floor_id"
// @Param building_id query string false "building_id"
// @Param type query string false "type"
// @Param power_consumption query string false "power_consumption"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Eco_floor_1y} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y [get]
func Eco_floor_1yListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Eco_floor_1y](w, r, common.GetDaprClient(), "f_eco_floor_1y", "time")
}

// @Summary save
// @Description save
// @Tags Eco_floor_1y
// @Accept       json
// @Param item body model.Eco_floor_1y true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Eco_floor_1y} "object"
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y [post]
func UpsertEco_floor_1yHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Eco_floor_1y
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Eco_floor_1y")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Eco_floor_1y)
	}

	err = common.DbUpsert[model.Eco_floor_1y](r.Context(), common.GetDaprClient(), val, model.Eco_floor_1yTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Eco_floor_1y
// @Param time  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Eco_floor_1y} "object"
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y/{id} [delete]
func DeleteEco_floor_1yHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Eco_floor_1y")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "f_eco_floor_1y", "time", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Eco_floor_1y
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /eco-floor-1y/batch-delete [post]
func batchDeleteEco_floor_1yHandler(w http.ResponseWriter, r *http.Request) {

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
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Eco_floor_1y")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "f_eco_floor_1y", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
