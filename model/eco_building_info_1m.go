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


Table: v_eco_building_info_1m
[ 0] time                                           TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 1] building_id                                    VARCHAR(32)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] building_name                                  VARCHAR(128)         null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 3] id                                             TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] total                                          NUMERIC              null: true   primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []
[ 5] types                                          JSON                 null: true   primary: false  isArray: false  auto: false  col: JSON            len: -1      default: []
[ 6] floors                                         JSON                 null: true   primary: false  isArray: false  auto: false  col: JSON            len: -1      default: []


JSON Sample
-------------------------------------
{    "time": 97,    "building_id": "wUZSMwSHvFVuiTtDIgLAMMSWf",    "building_name": "emrVYfMbXtnVXmkaXXfXqJMOs",    "id": "tDstXhkrovMyxIcAZZVLWxiVw",    "total": 0.11281511326088584,    "types": 73,    "floors": 75}


Comments
-------------------------------------
[ 0] Warning table: v_eco_building_info_1m does not have a primary key defined, setting col position 1 time as primary key
Warning table: v_eco_building_info_1m primary key column time is nullable column, setting it as NOT NULL




*/

var (
	Eco_building_info_1m_FIELD_NAME_time = "time"

	Eco_building_info_1m_FIELD_NAME_building_id = "building_id"

	Eco_building_info_1m_FIELD_NAME_building_name = "building_name"

	Eco_building_info_1m_FIELD_NAME_id = "id"

	Eco_building_info_1m_FIELD_NAME_total = "total"

	Eco_building_info_1m_FIELD_NAME_types = "types"

	Eco_building_info_1m_FIELD_NAME_floors = "floors"
)

// Eco_building_info_1m struct is a row record of the v_eco_building_info_1m table in the  database
type Eco_building_info_1m struct {
	Time         common.LocalTime `json:"time"`          //time
	BuildingID   string           `json:"building_id"`   //building_id
	BuildingName string           `json:"building_name"` //building_name
	ID           string           `json:"id"`            //id
	Total        float64          `json:"total"`         //total
	Types        any              `json:"types"`         //types
	Floors       any              `json:"floors"`        //floors

}

var Eco_building_info_1mTableInfo = &TableInfo{
	Name: "v_eco_building_info_1m",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:   0,
			Name:    "time",
			Comment: `time`,
			Notes: `Warning table: v_eco_building_info_1m does not have a primary key defined, setting col position 1 time as primary key
Warning table: v_eco_building_info_1m primary key column time is nullable column, setting it as NOT NULL
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
			Name:               "building_name",
			Comment:            `building_name`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "BuildingName",
			GoFieldType:        "string",
			JSONFieldName:      "building_name",
			ProtobufFieldName:  "building_name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "id",
			Comment:            `id`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "total",
			Comment:            `total`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "NUMERIC",
			DatabaseTypePretty: "NUMERIC",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "NUMERIC",
			ColumnLength:       -1,
			GoFieldName:        "Total",
			GoFieldType:        "float64",
			JSONFieldName:      "total",
			ProtobufFieldName:  "total",
			ProtobufType:       "float",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "types",
			Comment:            `types`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "JSON",
			DatabaseTypePretty: "JSON",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "JSON",
			ColumnLength:       -1,
			GoFieldName:        "Types",
			GoFieldType:        "any",
			JSONFieldName:      "types",
			ProtobufFieldName:  "types",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "floors",
			Comment:            `floors`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "JSON",
			DatabaseTypePretty: "JSON",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "JSON",
			ColumnLength:       -1,
			GoFieldName:        "Floors",
			GoFieldType:        "any",
			JSONFieldName:      "floors",
			ProtobufFieldName:  "floors",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Eco_building_info_1m) TableName() string {
	return "v_eco_building_info_1m"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Eco_building_info_1m) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Eco_building_info_1m) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Eco_building_info_1m) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Eco_building_info_1m) TableInfo() *TableInfo {
	return Eco_building_info_1mTableInfo
}
