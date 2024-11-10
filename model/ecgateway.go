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


Table: o_eco_gateway
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] created_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 2] created_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 3] updated_by                                     VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 4] updated_time                                   TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 5] model_name                                     VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 6] dev_name                                       VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[ 7] mac_addr                                       VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 8] cm_code                                        VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 9] project_code                                   VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[10] location                                       VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[11] floor_id                                       VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[12] building_id                                    VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[13] park_id                                        VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[14] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "wRHKrkPCLeOweiUkYXDZMhkrV",    "created_by": "xrmJGBoveyTgvnbyfTjrYYmGK",    "created_time": 21,    "updated_by": "ojctSCZQviQRSvlhvNehcXOAC",    "updated_time": 86,    "model_name": "OuKaMotRvmQIXGWaXtVoPXxBq",    "dev_name": "oOnpbAJdZOdAkcEZjDknTONvP",    "mac_addr": "WpyxVnNnuvsBMMWgJVhbQnlRw",    "cm_code": "wUpohZqCPrTUYBUdeIKSKqCjn",    "project_code": "wORVBRVFejZDDJioWAHcJHvQt",    "location": "IEbaMUJDokyuytaoYxwvtBSqi",    "floor_id": "jPKniXKcpsUbdiCGwSAfAhwRJ",    "building_id": "juBsfITcFZlwBdOcGTjVhBBAa",    "park_id": "SfZRiHfQixDNVJjvXNmoFrCvk",    "type": 51}



*/

var (
	Ecgateway_FIELD_NAME_id = "id"

	Ecgateway_FIELD_NAME_created_by = "created_by"

	Ecgateway_FIELD_NAME_created_time = "created_time"

	Ecgateway_FIELD_NAME_updated_by = "updated_by"

	Ecgateway_FIELD_NAME_updated_time = "updated_time"

	Ecgateway_FIELD_NAME_model_name = "model_name"

	Ecgateway_FIELD_NAME_dev_name = "dev_name"

	Ecgateway_FIELD_NAME_mac_addr = "mac_addr"

	Ecgateway_FIELD_NAME_cm_code = "cm_code"

	Ecgateway_FIELD_NAME_project_code = "project_code"

	Ecgateway_FIELD_NAME_location = "location"

	Ecgateway_FIELD_NAME_floor_id = "floor_id"

	Ecgateway_FIELD_NAME_building_id = "building_id"

	Ecgateway_FIELD_NAME_park_id = "park_id"

	Ecgateway_FIELD_NAME_type = "type"
)

// Ecgateway struct is a row record of the o_eco_gateway table in the  database
type Ecgateway struct {
	ID          string           `json:"id"`           //主键ID
	CreatedBy   string           `json:"created_by"`   //创建人
	CreatedTime common.LocalTime `json:"created_time"` //创建时间
	UpdatedBy   string           `json:"updated_by"`   //更新人
	UpdatedTime common.LocalTime `json:"updated_time"` //更新时间
	ModelName   string           `json:"model_name"`   //型号名称
	DevName     string           `json:"dev_name"`     //设备名称
	MacAddr     string           `json:"mac_addr"`     //MAC地址
	CmCode      string           `json:"cm_code"`      //通信码
	ProjectCode string           `json:"project_code"` //项目编码
	Location    string           `json:"location"`     //组织名称
	FloorID     string           `json:"floor_id"`     //楼层ID
	BuildingID  string           `json:"building_id"`  //楼栋ID
	ParkID      string           `json:"park_id"`      //园区ID
	Type        int32            `json:"type"`         //网关类型(1:AL,2:AP)

}

var EcgatewayTableInfo = &TableInfo{
	Name: "o_eco_gateway",
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
			Comment:            `设备名称`,
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
			Name:               "mac_addr",
			Comment:            `MAC地址`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       64,
			GoFieldName:        "MacAddr",
			GoFieldType:        "string",
			JSONFieldName:      "mac_addr",
			ProtobufFieldName:  "mac_addr",
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
			Name:               "project_code",
			Comment:            `项目编码`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       64,
			GoFieldName:        "ProjectCode",
			GoFieldType:        "string",
			JSONFieldName:      "project_code",
			ProtobufFieldName:  "project_code",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
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
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
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
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
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
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
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
			ProtobufPos:        14,
		},

		&ColumnInfo{
			Index:              14,
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
			ProtobufPos:        15,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Ecgateway) TableName() string {
	return "o_eco_gateway"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Ecgateway) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Ecgateway) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Ecgateway) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Ecgateway) TableInfo() *TableInfo {
	return EcgatewayTableInfo
}
