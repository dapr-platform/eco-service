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
[ 7] channel_no                                     VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 8] cm_code                                        VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 9] location                                       VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[10] floor_no                                       VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []
[11] building_no                                    VARCHAR(128)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 128     default: []


JSON Sample
-------------------------------------
{    "id": "KIUcvMRhULMrcNqMtwMMkSVyX",    "created_by": "EbmmCUkdYNNaNZRemoprvHYGQ",    "created_time": 52,    "updated_by": "YSVJQtcBYhutpjdjFqqQEvSop",    "updated_time": 83,    "model_name": "rcmZjphhaILkiQYNdlJFhMDPr",    "dev_name": "QMFujOGWECNhDFUiPVodtZWFn",    "channel_no": "owdgiPjLiOGyCfErphpSjenOO",    "cm_code": "hrIabecFmvgyenFvAyQTLDrnk",    "location": "hxoplLFKnZtfKAVEVfgdhKEsE",    "floor_no": "YqgoqumBCWorgPXYLaXCNRQJU",    "building_no": "CkgFxJvPcioAEaRxrcNwNXqyI"}



*/

var (
	Ecgateway_FIELD_NAME_id = "id"

	Ecgateway_FIELD_NAME_created_by = "created_by"

	Ecgateway_FIELD_NAME_created_time = "created_time"

	Ecgateway_FIELD_NAME_updated_by = "updated_by"

	Ecgateway_FIELD_NAME_updated_time = "updated_time"

	Ecgateway_FIELD_NAME_model_name = "model_name"

	Ecgateway_FIELD_NAME_dev_name = "dev_name"

	Ecgateway_FIELD_NAME_channel_no = "channel_no"

	Ecgateway_FIELD_NAME_cm_code = "cm_code"

	Ecgateway_FIELD_NAME_location = "location"

	Ecgateway_FIELD_NAME_floor_no = "floor_no"

	Ecgateway_FIELD_NAME_building_no = "building_no"
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
	ChannelNo   string           `json:"channel_no"`   //通道号
	CmCode      string           `json:"cm_code"`      //通信码
	Location    string           `json:"location"`     //组织名称
	FloorNo     string           `json:"floor_no"`     //楼层号
	BuildingNo  string           `json:"building_no"`  //楼栋号

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
			Name:               "floor_no",
			Comment:            `楼层号`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "FloorNo",
			GoFieldType:        "string",
			JSONFieldName:      "floor_no",
			ProtobufFieldName:  "floor_no",
			ProtobufType:       "string",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "building_no",
			Comment:            `楼栋号`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       128,
			GoFieldName:        "BuildingNo",
			GoFieldType:        "string",
			JSONFieldName:      "building_no",
			ProtobufFieldName:  "building_no",
			ProtobufType:       "string",
			ProtobufPos:        12,
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
