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


Table: o_eco_water_meter
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] created_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] created_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 3] updated_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 4] updated_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 5] model_name                                     VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 6] dev_name                                       VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 7] channel_no                                     VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 8] cm_code                                        VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 9] location                                       VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[10] building_id                                    VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[11] park_id                                        VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[12] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[13] total_value                                    NUMERIC              null: false  primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "jTSbrwlJhFlaTopogWlnVoDVw",    "created_by": "yRrBPEWHPomLNHmZVHkMQhiLW",    "created_time": 70,    "updated_by": "wLwjcOdbeuebwexrTBGlVmqbq",    "updated_time": 55,    "model_name": "qKYODMIkAwxVobuVBfkyJeGHC",    "dev_name": "wWhLleduZgMwAUbBcRhLDdrNV",    "channel_no": "NQbNDoxobbLUOqNyDFBfWyjuq",    "cm_code": "MKNSJtBCXtAPOuIBOFYebhuAw",    "location": "TJpitwXKRgPdcxRCXdJxMNYyM",    "building_id": "NhqAkQKJfyfPifxoKqrIsXAUs",    "park_id": "RYfKTbPJHsIEBSXMamBJqpGrT",    "type": 81,    "total_value": 0.5168010617691216}



*/

var (
	Ecwater_meter_FIELD_NAME_id = "id"

	Ecwater_meter_FIELD_NAME_created_by = "created_by"

	Ecwater_meter_FIELD_NAME_created_time = "created_time"

	Ecwater_meter_FIELD_NAME_updated_by = "updated_by"

	Ecwater_meter_FIELD_NAME_updated_time = "updated_time"

	Ecwater_meter_FIELD_NAME_model_name = "model_name"

	Ecwater_meter_FIELD_NAME_dev_name = "dev_name"

	Ecwater_meter_FIELD_NAME_channel_no = "channel_no"

	Ecwater_meter_FIELD_NAME_cm_code = "cm_code"

	Ecwater_meter_FIELD_NAME_location = "location"

	Ecwater_meter_FIELD_NAME_building_id = "building_id"

	Ecwater_meter_FIELD_NAME_park_id = "park_id"

	Ecwater_meter_FIELD_NAME_type = "type"

	Ecwater_meter_FIELD_NAME_total_value = "total_value"
)

// Ecwater_meter struct is a row record of the o_eco_water_meter table in the  database
type Ecwater_meter struct {
	ID string `json:"id"` //主键ID

	CreatedBy string `json:"created_by"` //创建人

	CreatedTime common.LocalTime `json:"created_time"` //创建时间

	UpdatedBy string `json:"updated_by"` //更新人

	UpdatedTime common.LocalTime `json:"updated_time"` //更新时间

	ModelName string `json:"model_name"` //型号名称

	DevName string `json:"dev_name"` //dev_name

	ChannelNo string `json:"channel_no"` //通道号

	CmCode string `json:"cm_code"` //通信码

	Location string `json:"location"` //组织名称

	BuildingID string `json:"building_id"` //楼栋ID

	ParkID string `json:"park_id"` //园区ID

	Type int32 `json:"type"` //水表类型(1:低区,2:高区)

	TotalValue float64 `json:"total_value"` //总用水量

}

var Ecwater_meterTableInfo = &TableInfo{
	Name: "o_eco_water_meter",
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
			Name:               "model_name",
			Comment:            `型号名称`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "ModelName",
			GoFieldType:        "string",
			JSONFieldName:      "model_name",
			ProtobufFieldName:  "model_name",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "dev_name",
			Comment:            `dev_name`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "DevName",
			GoFieldType:        "string",
			JSONFieldName:      "dev_name",
			ProtobufFieldName:  "dev_name",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "channel_no",
			Comment:            `通道号`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       64,
			GoFieldName:        "ChannelNo",
			GoFieldType:        "string",
			JSONFieldName:      "channel_no",
			ProtobufFieldName:  "channel_no",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "cm_code",
			Comment:            `通信码`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       64,
			GoFieldName:        "CmCode",
			GoFieldType:        "string",
			JSONFieldName:      "cm_code",
			ProtobufFieldName:  "cm_code",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "location",
			Comment:            `组织名称`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "Location",
			GoFieldType:        "string",
			JSONFieldName:      "location",
			ProtobufFieldName:  "location",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
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
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
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
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "type",
			Comment:            `水表类型(1:低区,2:高区)`,
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
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "total_value",
			Comment:            `总用水量`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "NUMERIC",
			DatabaseTypePretty: "NUMERIC",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "NUMERIC",
			ColumnLength:       -1,
			GoFieldName:        "TotalValue",
			GoFieldType:        "float64",
			JSONFieldName:      "total_value",
			ProtobufFieldName:  "total_value",
			ProtobufType:       "float",
			ProtobufPos:        14,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Ecwater_meter) TableName() string {
	return "o_eco_water_meter"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Ecwater_meter) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Ecwater_meter) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Ecwater_meter) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Ecwater_meter) TableInfo() *TableInfo {
	return Ecwater_meterTableInfo
}
