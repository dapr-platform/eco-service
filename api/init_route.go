package api

import "github.com/go-chi/chi/v5"

func InitRoute(r chi.Router) {
	InitEcgatewayRoute(r)
	InitEcfloorRoute(r)
	InitEcbuildingRoute(r)
	InitEcparkRoute(r)
	InitManuCollectRoute(r)
	InitDashboardRoute(r)
}
