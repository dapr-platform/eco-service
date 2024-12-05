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


Table: o_eco_building
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] created_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] created_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 3] updated_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 4] updated_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 5] building_name                                  VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 6] park_id                                        VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 7] index                                          INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "dfURNcSFsGjegUmpBNaidfOMT",    "created_by": "fLGBZkuTFJEkRfPaKLjxaGitq",    "created_time": 44,    "updated_by": "LrReVCjvnTGmljdfdwyvuOPgJ",    "updated_time": 25,    "building_name": "qhCcyCjOfrRnvqFgBmkMorNmG",    "park_id": "sVuxSFLGBtwJTEoLFrxigYkej",    "index": 8}



*/

var (
	Ecbuilding_FIELD_NAME_id = "id"

	Ecbuilding_FIELD_NAME_created_by = "created_by"

	Ecbuilding_FIELD_NAME_created_time = "created_time"

	Ecbuilding_FIELD_NAME_updated_by = "updated_by"

	Ecbuilding_FIELD_NAME_updated_time = "updated_time"

	Ecbuilding_FIELD_NAME_building_name = "building_name"

	Ecbuilding_FIELD_NAME_park_id = "park_id"

	Ecbuilding_FIELD_NAME_index = "index"
)

// Ecbuilding struct is a row record of the o_eco_building table in the  database
type Ecbuilding struct {
	ID string `json:"id"` //主键ID

	CreatedBy string `json:"created_by"` //创建人

	CreatedTime common.LocalTime `json:"created_time"` //创建时间

	UpdatedBy string `json:"updated_by"` //更新人

	UpdatedTime common.LocalTime `json:"updated_time"` //更新时间

	BuildingName string `json:"building_name"` //楼栋名称

	ParkID string `json:"park_id"` //园区ID

	Index int32 `json:"index"` //排序索引

}

var EcbuildingTableInfo = &TableInfo{
	Name: "o_eco_building",
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
			Name:               "building_name",
			Comment:            `楼栋名称`,
			Notes:              ``,
			Nullable:           false,
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
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
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
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
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
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Ecbuilding) TableName() string {
	return "o_eco_building"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Ecbuilding) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Ecbuilding) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Ecbuilding) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Ecbuilding) TableInfo() *TableInfo {
	return EcbuildingTableInfo
}
