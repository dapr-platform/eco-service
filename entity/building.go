package entity

type BuildingPowerConsumptionInfoList []BuildingPowerConsumptionInfo

// 楼宇能耗信息
type BuildingPowerConsumptionInfo struct {
	ID         string  `json:"id"`
	BuildingID string  `json:"building_id"`
	BuildingName string  `json:"building_name"`
	Total      float64 `json:"total"`
	Types      []BuildingPowerConsumptionType `json:"types"`
	Floors     []FloorPowerConsumptionInfo    `json:"floors"`
}
type BuildingPowerConsumptionType struct {
	ID   int  `json:"id"`
	Total float64 `json:"total"`
}
type FloorPowerConsumptionInfo struct {
	ID         string  `json:"id"`
	FloorID    string  `json:"floor_id"`
	FloorName  string  `json:"floor_name"`
	Total      float64 `json:"total"`
	Types      []FloorPowerConsumptionType `json:"types"`
}
type FloorPowerConsumptionType struct {
	ID   int     `json:"id"`
	Total float64 `json:"total"`
}
