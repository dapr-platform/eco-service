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


Table: o_eco_floor
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] created_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] created_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 3] updated_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 4] updated_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 5] floor_name                                     VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 6] building_id                                    VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 7] park_id                                        VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 8] index                                          INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "KnSOxiEJClsAoisALeioQDfLM",    "created_by": "QoHTvPIWOeZQjCmThPHTtUHed",    "created_time": 81,    "updated_by": "oVKndODNeCPcLgnTysmuOxErs",    "updated_time": 7,    "floor_name": "lFWPVlAwxgBYrLpgHijJPSBJV",    "building_id": "GIlPYumJRQiOADRfIaInMNePI",    "park_id": "ykQRlwBVxyPoroFgyPMCvbrkv",    "index": 70}



*/

var (
	Ecfloor_FIELD_NAME_id = "id"

	Ecfloor_FIELD_NAME_created_by = "created_by"

	Ecfloor_FIELD_NAME_created_time = "created_time"

	Ecfloor_FIELD_NAME_updated_by = "updated_by"

	Ecfloor_FIELD_NAME_updated_time = "updated_time"

	Ecfloor_FIELD_NAME_floor_name = "floor_name"

	Ecfloor_FIELD_NAME_building_id = "building_id"

	Ecfloor_FIELD_NAME_park_id = "park_id"

	Ecfloor_FIELD_NAME_index = "index"
)

// Ecfloor struct is a row record of the o_eco_floor table in the  database
type Ecfloor struct {
	ID string `json:"id"` //主键ID

	CreatedBy string `json:"created_by"` //创建人

	CreatedTime common.LocalTime `json:"created_time"` //创建时间

	UpdatedBy string `json:"updated_by"` //更新人

	UpdatedTime common.LocalTime `json:"updated_time"` //更新时间

	FloorName string `json:"floor_name"` //楼层名称

	BuildingID string `json:"building_id"` //楼栋ID

	ParkID string `json:"park_id"` //园区ID

	Index int32 `json:"index"` //排序索引

}

var EcfloorTableInfo = &TableInfo{
	Name: "o_eco_floor",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            `主键ID`,
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
			Name:               "created_by",
			Comment:            `创建人`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "CreatedBy",
			GoFieldType:        "string",
			JSONFieldName:      "created_by",
			ProtobufFieldName:  "created_by",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "created_time",
			Comment:            `创建时间`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "CreatedTime",
			GoFieldType:        "common.LocalTime",
			JSONFieldName:      "created_time",
			ProtobufFieldName:  "created_time",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "updated_by",
			Comment:            `更新人`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "UpdatedBy",
			GoFieldType:        "string",
			JSONFieldName:      "updated_by",
			ProtobufFieldName:  "updated_by",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "updated_time",
			Comment:            `更新时间`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "UpdatedTime",
			GoFieldType:        "common.LocalTime",
			JSONFieldName:      "updated_time",
			ProtobufFieldName:  "updated_time",
			ProtobufType:       "uint64",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "floor_name",
			Comment:            `楼层名称`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "FloorName",
			GoFieldType:        "string",
			JSONFieldName:      "floor_name",
			ProtobufFieldName:  "floor_name",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
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
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
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
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "index",
			Comment:            `排序索引`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Index",
			GoFieldType:        "int32",
			JSONFieldName:      "index",
			ProtobufFieldName:  "index",
			ProtobufType:       "int32",
			ProtobufPos:        9,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Ecfloor) TableName() string {
	return "o_eco_floor"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Ecfloor) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Ecfloor) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Ecfloor) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Ecfloor) TableInfo() *TableInfo {
	return EcfloorTableInfo
}
