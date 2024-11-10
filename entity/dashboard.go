package entity

type LabelData struct {
	Id    string  `json:"id"`
	Label string  `json:"label"`
	Value float64 `json:"value"`
	TB    float64 `json:"tb"`
	HB    float64 `json:"hb"`
}
