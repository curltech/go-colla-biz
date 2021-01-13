package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

const (
	UserStatus_Enabled   string = "Enabled"
	UserStatus_Expired   string = "Expired"
	UserStatus_Locked    string = "Locked"
	UserStatus_Disabled  string = "Disabled"
	UserStatus_Discarded string = "Discarded"
)

/**
 * The enum Temp password authenticate status.
 */
const (
	AuthenticateStatus_Success    string = "Success"
	AuthenticateStatus_Expired    string = "Expired"
	AuthenticateStatus_NotMatch   string = "NotMatch"
	AuthenticateStatus_NotExist   string = "NotExist"
	AuthenticateStatus_Locked     string = "Locked"
	AuthenticateStatus_NoPassword string = "NoPassword"
)

type User struct {
	entity.StatusEntity `xorm:"extends"`
	UserId              string     `xorm:"varchar(255)" json:"userId,omitempty"`
	NickName            string     `xorm:"varchar(255)" json:"nickName,omitempty"`
	UserName            string     `xorm:"varchar(255)" json:"userName,omitempty"`
	Name                string     `xorm:"varchar(255)" json:"name,omitempty"`
	Password            string     `xorm:"varchar(255)" json:"-,omitempty"`
	PlainPassword       string     `xorm:"-" json:"plainPassword,omitempty"`
	ConfirmPassword     string     `xorm:"-" json:"confirmPassword,omitempty"`
	PublicKey           string     `xorm:"varchar(2048)" json:"publicKey,omitempty"`
	SecurityContext     string     `xorm:"varchar(255)" json:"securityContext,omitempty"`
	Description         string     `xorm:"varchar(255)" json:"description,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	Email               string     `xorm:"varchar(255)" json:"email,omitempty"`
	IdentifyType        string     `xorm:"varchar(255)" json:"identifyType,omitempty"`
	IdentifyNumber      string     `xorm:"varchar(255)" json:"identifyNumber,omitempty"`
	Mobile              string     `xorm:"varchar(255)" json:"mobile,omitempty"`
	OwnedStructureId    string     `xorm:"varchar(255)" json:"ownedStructureIdownedStructureId,omitempty"`
	OwnedStructureName  string     `xorm:"varchar(255)" json:"ownedStructureName,omitempty"`
	OwnedStructurePath  string     `xorm:"varchar(255)" json:"ownedStructurePath,omitempty"`
	EmployeeId          string     `xorm:"varchar(255)" json:"employeeId,omitempty"`
	IpAddress           string     `xorm:"varchar(255)" json:"ipAddress,omitempty"`
	WechatId            string     `xorm:"varchar(255)" json:"wechatId,omitempty"`
	Avatar              string     `xorm:"varchar(4096)" json:"avatar,omitempty"`
}

func (User) TableName() string {
	return "rbac_user"
}

func (User) KeyName() string {
	return "UserId"
}

func (User) IdName() string {
	return entity.FieldName_Id
}
