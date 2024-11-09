package api

import "github.com/go-chi/chi/v5"

func InitRoute(r chi.Router) {
	InitEco_building_info_1dRoute(r)
	InitEco_building_info_1mRoute(r)
	InitEco_building_info_1yRoute(r)
	InitEcgatewayRoute(r)
	InitManuCollectRoute(r)
}
