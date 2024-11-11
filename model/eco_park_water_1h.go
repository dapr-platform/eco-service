package model

import (
	"database/sql"
	"github.com/dapr-platform/common"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = common.LocalTime{}
)

/*
DB Table Details
-------------------------------------


Table: f_eco_park_water_1h
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] time                                           TIMESTAMP            null: false  primary: true   isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 2] park_id                                        VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 3] water_consumption                              NUMERIC              null: false  primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "ttSOyNQGHdZDcdaNFOmIdcGJG",    "time": 53,    "park_id": "ZdmLyKVDsVVmrGqihsYGlPmYp",    "water_consumption": 0.8183260153519903}



*/

var (
	Eco_park_water_1h_FIELD_NAME_id = "id"

	Eco_park_water_1h_FIELD_NAME_time = "time"

	Eco_park_water_1h_FIELD_NAME_park_id = "park_id"

	Eco_park_water_1h_FIELD_NAME_water_consumption = "water_consumption"
)

// Eco_park_water_1h struct is a row record of the f_eco_park_water_1h table in the  database
type Eco_park_water_1h struct {
	ID               string           `json:"id"`                //id
	Time             common.LocalTime `json:"time"`              //时间
	ParkID           string           `json:"park_id"`           //园区ID
	WaterConsumption float64          `json:"water_consumption"` //用水量(m³)

}

var Eco_park_water_1hTableInfo = &TableInfo{
	Name: "f_eco_park_water_1h",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            `id`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "time",
			Comment:            `时间`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "Time",
			GoFieldType:        "common.LocalTime",
			JSONFieldName:      "time",
			ProtobufFieldName:  "time",
			ProtobufType:       "uint64",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "park_id",
			Comment:            `园区ID`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "ParkID",
			GoFieldType:        "string",
			JSONFieldName:      "park_id",
			ProtobufFieldName:  "park_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "water_consumption",
			Comment:            `用水量(m³)`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "NUMERIC",
			DatabaseTypePretty: "NUMERIC",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "NUMERIC",
			ColumnLength:       -1,
			GoFieldName:        "WaterConsumption",
			GoFieldType:        "float64",
			JSONFieldName:      "water_consumption",
			ProtobufFieldName:  "water_consumption",
			ProtobufType:       "float",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Eco_park_water_1h) TableName() string {
	return "f_eco_park_water_1h"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Eco_park_water_1h) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Eco_park_water_1h) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Eco_park_water_1h) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Eco_park_water_1h) TableInfo() *TableInfo {
	return Eco_park_water_1hTableInfo
}
