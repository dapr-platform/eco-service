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


Table: f_eco_gateway_1h
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] time                                           TIMESTAMP            null: false  primary: true   isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 2] gateway_id                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 3] floor_id                                       VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 4] building_id                                    VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 5] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 6] power_consumption                              NUMERIC              null: false  primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "tGFRRitXTCAyBpMwcyxOBktpW",    "time": 58,    "gateway_id": "lppcpqXjsGeQaHXBVCfASPXDO",    "floor_id": "vfmNTFlTfJMHwUKvTrsuDlZre",    "building_id": "alwABjiYbdroJuMDvQZCkqHhm",    "type": 21,    "power_consumption": 0.26447821174157077}



*/

var (
	Eco_gateway_1h_FIELD_NAME_id = "id"

	Eco_gateway_1h_FIELD_NAME_time = "time"

	Eco_gateway_1h_FIELD_NAME_gateway_id = "gateway_id"

	Eco_gateway_1h_FIELD_NAME_floor_id = "floor_id"

	Eco_gateway_1h_FIELD_NAME_building_id = "building_id"

	Eco_gateway_1h_FIELD_NAME_type = "type"

	Eco_gateway_1h_FIELD_NAME_power_consumption = "power_consumption"
)

// Eco_gateway_1h struct is a row record of the f_eco_gateway_1h table in the  database
type Eco_gateway_1h struct {
	ID               string           `json:"id"`                //id
	Time             common.LocalTime `json:"time"`              //时间
	GatewayID        string           `json:"gateway_id"`        //网关ID
	FloorID          string           `json:"floor_id"`          //楼层ID
	BuildingID       string           `json:"building_id"`       //楼栋ID
	Type             int32            `json:"type"`              //网关类型(1:AL,2:AP)
	PowerConsumption float64          `json:"power_consumption"` //用电量(kWh)

}

var Eco_gateway_1hTableInfo = &TableInfo{
	Name: "f_eco_gateway_1h",
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
			Name:               "gateway_id",
			Comment:            `网关ID`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "GatewayID",
			GoFieldType:        "string",
			JSONFieldName:      "gateway_id",
			ProtobufFieldName:  "gateway_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "floor_id",
			Comment:            `楼层ID`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "FloorID",
			GoFieldType:        "string",
			JSONFieldName:      "floor_id",
			ProtobufFieldName:  "floor_id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "building_id",
			Comment:            `楼栋ID`,
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
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "type",
			Comment:            `网关类型(1:AL,2:AP)`,
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
			Name:               "power_consumption",
			Comment:            `用电量(kWh)`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "NUMERIC",
			DatabaseTypePretty: "NUMERIC",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "NUMERIC",
			ColumnLength:       -1,
			GoFieldName:        "PowerConsumption",
			GoFieldType:        "float64",
			JSONFieldName:      "power_consumption",
			ProtobufFieldName:  "power_consumption",
			ProtobufType:       "float",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Eco_gateway_1h) TableName() string {
	return "f_eco_gateway_1h"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Eco_gateway_1h) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Eco_gateway_1h) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Eco_gateway_1h) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Eco_gateway_1h) TableInfo() *TableInfo {
	return Eco_gateway_1hTableInfo
}
