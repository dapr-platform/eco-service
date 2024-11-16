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


Table: f_eco_building_1m
[ 0] time                                           TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 1] building_id                                    VARCHAR(32)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] park_id                                        VARCHAR(32)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 3] type                                           INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 4] power_consumption                              NUMERIC              null: true   primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []


JSON Sample
-------------------------------------
{    "time": 45,    "building_id": "QTtUDduxkDTuJAxkmwCPrhCHD",    "park_id": "XkuiJUvrfsWdSbEyIXJdELWfd",    "type": 82,    "power_consumption": 0.7088988431618581}


Comments
-------------------------------------
[ 0] Warning table: f_eco_building_1m does not have a primary key defined, setting col position 1 time as primary key
Warning table: f_eco_building_1m primary key column time is nullable column, setting it as NOT NULL




*/

var (
	Eco_building_1m_FIELD_NAME_time = "time"

	Eco_building_1m_FIELD_NAME_building_id = "building_id"

	Eco_building_1m_FIELD_NAME_park_id = "park_id"

	Eco_building_1m_FIELD_NAME_type = "type"

	Eco_building_1m_FIELD_NAME_power_consumption = "power_consumption"
)

// Eco_building_1m struct is a row record of the f_eco_building_1m table in the  database
type Eco_building_1m struct {
	Time             common.LocalTime `json:"time"`              //time
	BuildingID       string           `json:"building_id"`       //building_id
	ParkID           string           `json:"park_id"`           //park_id
	Type             int32            `json:"type"`              //type
	PowerConsumption float64          `json:"power_consumption"` //power_consumption

}

var Eco_building_1mTableInfo = &TableInfo{
	Name: "f_eco_building_1m",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:   0,
			Name:    "time",
			Comment: `time`,
			Notes: `Warning table: f_eco_building_1m does not have a primary key defined, setting col position 1 time as primary key
Warning table: f_eco_building_1m primary key column time is nullable column, setting it as NOT NULL
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
			Name:               "building_id",
			Comment:            `building_id`,
			Notes:              ``,
			Nullable:           true,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "power_consumption",
			Comment:            `power_consumption`,
			Notes:              ``,
			Nullable:           true,
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
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Eco_building_1m) TableName() string {
	return "f_eco_building_1m"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Eco_building_1m) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Eco_building_1m) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Eco_building_1m) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Eco_building_1m) TableInfo() *TableInfo {
	return Eco_building_1mTableInfo
}
