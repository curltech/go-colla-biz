package entity

import (
	"errors"
	"github.com/curltech/go-colla-biz/spec/entity"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/util/convert"
)

const (
	AttributeType_SpecId   string = "SpecId"
	AttributeType_Kind     string = "Kind"
	AttributeType_Path     string = "Path"
	AttributeType_DataType string = "DataType"
	AttributeType_Pattern  string = "Pattern"
	AttributeType_Value    string = "Value"
)

const AttributeSpecNumber = 15

type Property struct {
	entity.InternalFixedActual `xorm:"extends"`
	CurrentIndex               int `json:"CurrentIndex,omitempty"`
	SerialId                   int `json:"SerialId,omitempty"`

	SpecId0 uint64 `json:"specId0,omitempty"`
	Value0  string `xorm:"varchar(255)" json:"value0,omitempty"`

	SpecId1 uint64 `json:"specId1,omitempty"`
	Value1  string `xorm:"varchar(255)" json:"value1,omitempty"`

	SpecId2 uint64 `json:"specId2,omitempty"`
	Value2  string `xorm:"varchar(255)" json:"value2,omitempty"`

	SpecId3 uint64 `json:"specId3,omitempty"`
	Value3  string `xorm:"varchar(255)" json:"value3,omitempty"`

	SpecId4 uint64 `json:"specId4,omitempty"`
	Value4  string `xorm:"varchar(255)" json:"value4,omitempty"`

	SpecId5 uint64 `json:"specId5,omitempty"`
	Value5  string `xorm:"varchar(255)" json:"value5,omitempty"`

	SpecId6 uint64 `json:"specId6,omitempty"`
	Value6  string `xorm:"varchar(255)" json:"value6,omitempty"`

	SpecId7 uint64 `json:"specId7,omitempty"`
	Value7  string `xorm:"varchar(255)" json:"value7,omitempty"`

	SpecId8 uint64 `json:"specId8,omitempty"`
	Value8  string `xorm:"varchar(255)" json:"value8,omitempty"`

	SpecId9 uint64 `json:"specId9,omitempty"`
	Value9  string `xorm:"varchar(255)" json:"value9,omitempty"`

	SpecId10 uint64 `json:"specId10,omitempty"`
	Value10  string `xorm:"varchar(255)" json:"value10,omitempty"`

	SpecId11 uint64 `json:"specId11,omitempty"`
	Value11  string `xorm:"varchar(255)" json:"value11,omitempty"`

	SpecId12 uint64 `json:"specId12,omitempty"`
	Value12  string `xorm:"varchar(255)" json:"value12,omitempty"`

	SpecId13 uint64 `json:"specId13,omitempty"`
	Value13  string `xorm:"varchar(255)" json:"value13,omitempty"`

	SpecId14 uint64 `json:"specId14,omitempty"`
	Value14  string `xorm:"varchar(255)" json:"value14,omitempty"`

	computed map[string][]string `xorm:"-" json:"-"`
	position map[string]int      `xorm:"-" json:"-"`
}

func (Property) TableName() string {
	return "atl_property"
}

func (Property) KeyName() string {
	return "ActualId"
}

func (Property) IdName() string {
	return baseentity.FieldName_Id
}

func NewProperty() *Property {
	var property = Property{}

	return &property
}

func (this *Property) Get(attributeType string, index int) (interface{}, error) {
	if index >= AttributeSpecNumber {
		return nil, errors.New("")
	}
	if attributeType == AttributeType_SpecId {
		switch index {
		case 0:
			return this.SpecId0, nil
		case 1:
			return this.SpecId1, nil
		case 2:
			return this.SpecId2, nil
		case 3:
			return this.SpecId3, nil
		case 4:
			return this.SpecId4, nil
		case 5:
			return this.SpecId5, nil
		case 6:
			return this.SpecId6, nil
		case 7:
			return this.SpecId7, nil
		case 8:
			return this.SpecId8, nil
		case 9:
			return this.SpecId9, nil
		case 10:
			return this.SpecId10, nil
		case 11:
			return this.SpecId11, nil
		case 12:
			return this.SpecId12, nil
		case 13:
			return this.SpecId13, nil
		case 14:
			return this.SpecId14, nil
		default:
			return nil, errors.New("")
		}
	} else if attributeType == AttributeType_Value {
		switch index {
		case 0:
			return this.Value0, nil
		case 1:
			return this.Value1, nil
		case 2:
			return this.Value2, nil
		case 3:
			return this.Value3, nil
		case 4:
			return this.Value4, nil
		case 5:
			return this.Value5, nil
		case 6:
			return this.Value6, nil
		case 7:
			return this.Value7, nil
		case 8:
			return this.Value8, nil
		case 9:
			return this.Value9, nil
		case 10:
			return this.Value10, nil
		case 11:
			return this.Value11, nil
		case 12:
			return this.Value12, nil
		case 13:
			return this.Value13, nil
		case 14:
			return this.Value14, nil
		default:
			return nil, errors.New("")
		}
	} else {
		this.initComputed()
		c, ok := this.computed[attributeType]
		if ok {
			v := c[index]
			return v, nil
		}
	}

	return nil, errors.New("")
}

func (this *Property) getValue(index int) (interface{}, error) {
	_, err := this.Get(AttributeType_SpecId, index)
	if err != nil {
		return nil, errors.New("")
	}
	dataType, _ := this.Get(AttributeType_DataType, index)
	v, _ := this.Get(AttributeType_Value, index)

	return convert.ToObject(v.(string), dataType.(string))
}

func (this *Property) GetValue(kind string) (interface{}, error) {
	serialId, ok := this.position[kind]
	if ok {
		return this.getValue(serialId)
	}

	return nil, errors.New("")
}

func (this *Property) Set(attributeType string, index int, value interface{}) error {
	if index >= AttributeSpecNumber {
		return errors.New("")
	}
	if attributeType == AttributeType_SpecId {
		v := value.(uint64)
		switch index {
		case 0:
			this.SpecId0 = v
		case 1:
			this.SpecId1 = v
		case 2:
			this.SpecId2 = v
		case 3:
			this.SpecId3 = v
		case 4:
			this.SpecId4 = v
		case 5:
			this.SpecId5 = v
		case 6:
			this.SpecId6 = v
		case 7:
			this.SpecId7 = v
		case 8:
			this.SpecId8 = v
		case 9:
			this.SpecId9 = v
		case 10:
			this.SpecId10 = v
		case 11:
			this.SpecId11 = v
		case 12:
			this.SpecId12 = v
		case 13:
			this.SpecId13 = v
		case 14:
			this.SpecId14 = v
		default:
			return errors.New("")
		}
	} else if attributeType == AttributeType_Value {
		v := value.(string)
		switch index {
		case 0:
			this.Value0 = v
		case 1:
			this.Value1 = v
		case 2:
			this.Value2 = v
		case 3:
			this.Value3 = v
		case 4:
			this.Value4 = v
		case 5:
			this.Value5 = v
		case 6:
			this.Value6 = v
		case 7:
			this.Value7 = v
		case 8:
			this.Value8 = v
		case 9:
			this.Value9 = v
		case 10:
			this.Value10 = v
		case 11:
			this.Value11 = v
		case 12:
			this.Value12 = v
		case 13:
			this.Value13 = v
		case 14:
			this.Value14 = v
		default:
			return errors.New("")
		}
	} else {
		this.initComputed()
		v := value.(string)
		c, ok := this.computed[attributeType]
		if ok {
			c[index] = v
		} else {
			return errors.New("ErrorComputedValue")
		}
		if attributeType == AttributeType_Kind {
			this.position[v] = index
		}
	}

	return nil
}

func (this *Property) initComputed() {
	if this.computed == nil {
		this.computed = make(map[string][]string, 4)
		this.computed[AttributeType_Kind] = make([]string, AttributeSpecNumber)
		this.computed[AttributeType_Path] = make([]string, AttributeSpecNumber)
		this.computed[AttributeType_DataType] = make([]string, AttributeSpecNumber)
		this.computed[AttributeType_Pattern] = make([]string, AttributeSpecNumber)
		this.position = make(map[string]int, AttributeSpecNumber)
	}
}

func (this *Property) PutValue(serialId int, value interface{}) (bool, error) {
	dataType, _ := this.Get(AttributeType_DataType, serialId)
	pattern, _ := this.Get(AttributeType_Pattern, serialId)
	v, err := convert.ToString(value, dataType.(string), pattern.(string))
	if err != nil {
		return false, err
	}
	old, err := this.Get(AttributeType_Value, serialId)
	if err != nil {
		return false, err
	}
	if old == v {
		return false, nil
	} else {
		err = this.Set(AttributeType_Value, serialId, v)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

func (this *Property) SetValue(kind string, value interface{}) (bool, error) {
	serialId, ok := this.position[kind]
	if ok {
		return this.PutValue(serialId, value)
	}

	return false, errors.New("NotExist")
}

func (this *Property) Contain(kind string) bool {
	_, ok := this.position[kind]
	if ok {
		return true
	}

	return false
}
