package entity

import (
	baseentity "github.com/curltech/go-colla-core/entity"
)

const (
	DataType_String    string = "String"
	DataType_Number    string = "Number"
	DataType_Integer   string = "Integer"
	DataType_Long      string = "Long"
	DataType_Float     string = "Float"
	DataType_Date      string = "Date"
	DataType_Timestamp string = "Timestamp"
	DataType_Uint64    string = "uint64"
	DataType_Int64     string = "int64"
	DataType_Uint      string = "uint"
	DataType_Int       string = "int"
	DataType_Bool      string = "bool"
)

type AttributeSpec struct {
	Specification `xorm:"extends"`
	DataType      string `xorm:"varchar(255)" json:"dataType,omitempty"`
	// 存储格式
	Pattern string `xorm:"varchar(255)" json:"pattern,omitempty"`
	/** 缺省值 */
	DefaultValue string `xorm:"varchar(255)" json:"defaultValue,omitempty"`
	/** 是否可选 */
	Required bool `json:"required,omitempty"`
	/** 被容许的值，是BaseCodeId */
	AllowedValue string `xorm:"varchar(255)" json:"allowedValue,omitempty"`
	Alias        string `xorm:"varchar(255)" json:"alias,omitempty"`
}

func (AttributeSpec) TableName() string {
	return "spec_attribute"
}

func (AttributeSpec) KeyName() string {
	return "SpecId"
}

func (AttributeSpec) IdName() string {
	return baseentity.FieldName_Id
}
