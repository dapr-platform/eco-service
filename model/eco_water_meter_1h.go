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


Table: f_eco_water_meter_1h
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] time                                           TIMESTAMP            null: false  primary: true   isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 2] water_meter_id                                 VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 3] building_id                                    VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 4] park_id                                        VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 5] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 6] water_consumption                              NUMERIC              null: false  primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "GFPHkaseqHkveTlhxJgDSNjwE",    "time": 55,    "water_meter_id": "NQYoFtkpEgJPYbTlvTlwUpHuW",    "building_id": "rryUjWgOniZyBqqKbGcPynehr",    "park_id": "mXKurxyhgEuKDwZbNNugtuqmP",    "type": 0,    "water_consumption": 0.21503290081053947}



*/

var (
	Eco_water_meter_1h_FIELD_NAME_id = "id"

	Eco_water_meter_1h_FIELD_NAME_time = "time"

	Eco_water_meter_1h_FIELD_NAME_water_meter_id = "water_meter_id"

	Eco_water_meter_1h_FIELD_NAME_building_id = "building_id"

	Eco_water_meter_1h_FIELD_NAME_park_id = "park_id"

	Eco_water_meter_1h_FIELD_NAME_type = "type"

	Eco_water_meter_1h_FIELD_NAME_water_consumption = "water_consumption"
)

// Eco_water_meter_1h struct is a row record of the f_eco_water_meter_1h table in the  database
type Eco_water_meter_1h struct {
	ID               string           `json:"id"`                //id
	Time             common.LocalTime `json:"time"`              //time
	WaterMeterID     string           `json:"water_meter_id"`    //water_meter_id
	BuildingID       string           `json:"building_id"`       //building_id
	ParkID           string           `json:"park_id"`           //park_id
	Type             int32            `json:"type"`              //type
	WaterConsumption float64          `json:"water_consumption"` //water_consumption

}

var Eco_water_meter_1hTableInfo = &TableInfo{
	Name: "f_eco_water_meter_1h",
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
			Comment:            `time`,
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
			Name:               "water_meter_id",
			Comment:            `water_meter_id`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "WaterMeterID",
			GoFieldType:        "string",
			JSONFieldName:      "water_meter_id",
			ProtobufFieldName:  "water_meter_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "building_id",
			Comment:            `building_id`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "BuildingID",
			GoFieldType:        "string",
			JSONFieldName:      "building_id",
			ProtobufFieldName:  "building_id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "park_id",
			Comment:            `park_id`,
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
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "type",
			Comment:            `type`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Type",
			GoFieldType:        "int32",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "water_consumption",
			Comment:            `water_consumption`,
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
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Eco_water_meter_1h) TableName() string {
	return "f_eco_water_meter_1h"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Eco_water_meter_1h) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Eco_water_meter_1h) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Eco_water_meter_1h) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Eco_water_meter_1h) TableInfo() *TableInfo {
	return Eco_water_meter_1hTableInfo
}
