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


Table: f_eco_park_water_1y
[ 0] time                                           TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 1] park_id                                        VARCHAR(32)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] type                                           INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 3] water_consumption                              NUMERIC              null: true   primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []


JSON Sample
-------------------------------------
{    "time": 46,    "park_id": "DKRhvHHcCKXtjiORAUMlrpqeh",    "type": 69,    "water_consumption": 0.7636568433237264}


Comments
-------------------------------------
[ 0] Warning table: f_eco_park_water_1y does not have a primary key defined, setting col position 1 time as primary key
Warning table: f_eco_park_water_1y primary key column time is nullable column, setting it as NOT NULL




*/

var (
	Eco_park_water_1y_FIELD_NAME_time = "time"

	Eco_park_water_1y_FIELD_NAME_park_id = "park_id"

	Eco_park_water_1y_FIELD_NAME_type = "type"

	Eco_park_water_1y_FIELD_NAME_water_consumption = "water_consumption"
)

// Eco_park_water_1y struct is a row record of the f_eco_park_water_1y table in the  database
type Eco_park_water_1y struct {
	Time             common.LocalTime `json:"time"`              //time
	ParkID           string           `json:"park_id"`           //park_id
	Type             int32            `json:"type"`              //type
	WaterConsumption float64          `json:"water_consumption"` //water_consumption

}

var Eco_park_water_1yTableInfo = &TableInfo{
	Name: "f_eco_park_water_1y",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:   0,
			Name:    "time",
			Comment: `time`,
			Notes: `Warning table: f_eco_park_water_1y does not have a primary key defined, setting col position 1 time as primary key
Warning table: f_eco_park_water_1y primary key column time is nullable column, setting it as NOT NULL
`,
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
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "park_id",
			Comment:            `park_id`,
			Notes:              ``,
			Nullable:           true,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "type",
			Comment:            `type`,
			Notes:              ``,
			Nullable:           true,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "water_consumption",
			Comment:            `water_consumption`,
			Notes:              ``,
			Nullable:           true,
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
func (e *Eco_park_water_1y) TableName() string {
	return "f_eco_park_water_1y"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Eco_park_water_1y) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Eco_park_water_1y) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Eco_park_water_1y) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Eco_park_water_1y) TableInfo() *TableInfo {
	return Eco_park_water_1yTableInfo
}
