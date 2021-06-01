package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type ShareIndex struct {
	entity.BaseEntity `xorm:"extends"`
	/**
	 * 编号
	 */
	ShareId string `xorm:"varchar(255)" json:"shareId,omitempty"`
	/**
	 * 指标时间，按照季度2021-01
	 */
	IndexTime string `xorm:"varchar(255)" json:"indexTime,omitempty"`

	NetProfitGrowthRate    float64 `json:"kind,omitempty"`
	GrossRevenueGrowthRate float64 `json:"name,omitempty"`
	ReturnOnNetAsset       float64 `json:"name,omitempty"`
	GrossProfitMargin      float64 `json:"name,omitempty"`
	NetOperatingRate       float64 `json:"name,omitempty"`
	PriceEarningsRatio     float64 `json:"name,omitempty"`

	NetProfit             float64 `json:"name,omitempty"`
	TotalOperatingRevenue float64 `json:"name,omitempty"`
	BasicEarningsPerShare float64 `json:"name,omitempty"`

	GrowthRateOfTotalAsset      float64 `json:"name,omitempty"`
	GrowthRateOfOperatingProfit float64 `json:"name,omitempty"`
	GrowthRateOfNetAsset        float64 `json:"name,omitempty"`

	CapitalAccumulationFundPerShare            float64 `json:"name,omitempty"`
	UndistributedProfitPerShare                float64 `json:"name,omitempty"`
	NetAssetValuePerShare                      float64 `json:"name,omitempty"`
	OperatingCashFlowPerShare                  float64 `json:"name,omitempty"`
	GrowthOfNetCashFlowFromOperatingActivities float64 `json:"name,omitempty"`

	DeductionOfNonNetProfit      float64 `json:"name,omitempty"`
	AssetLiabilityRatio          float64 `json:"name,omitempty"`
	CurrentRatio                 float64 `json:"name,omitempty"`
	QuickRatio                   float64 `json:"name,omitempty"`
	InventoryTurnover            float64 `json:"name,omitempty"`
	TurnoverOfCurrentAsset       float64 `json:"name,omitempty"`
	TurnoverRateOfFixedAsset     float64 `json:"name,omitempty"`
	TurnoverOfTotalAsset         float64 `json:"name,omitempty"`
	GrowthRateOfCashFlowPerShare float64 `json:"name,omitempty"`

	CostProfitMargin   float64 `json:"name,omitempty"`
	ReturnOnTotalAsset float64 `json:"name,omitempty"`
}

func (ShareIndex) TableName() string {
	return "stk_index"
}

func (ShareIndex) KeyName() string {
	return entity.FieldName_Id
}

func (ShareIndex) IdName() string {
	return entity.FieldName_Id
}
