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

	NetProfitGrowthRate    float64 `json:"NetProfitGrowthRate,omitempty"`
	GrossRevenueGrowthRate float64 `json:"GrossRevenueGrowthRate,omitempty"`
	ReturnOnNetAsset       float64 `json:"ReturnOnNetAsset,omitempty"`
	GrossProfitMargin      float64 `json:"GrossProfitMargin,omitempty"`
	NetOperatingRate       float64 `json:"NetOperatingRate,omitempty"`
	PriceEarningsRatio     float64 `json:"PriceEarningsRatio,omitempty"`

	NetProfit             float64 `json:"NetProfit,omitempty"`
	TotalOperatingRevenue float64 `json:"TotalOperatingRevenue,omitempty"`
	BasicEarningsPerShare float64 `json:"BasicEarningsPerShare,omitempty"`

	GrowthRateOfTotalAsset      float64 `json:"GrowthRateOfTotalAsset,omitempty"`
	GrowthRateOfOperatingProfit float64 `json:"GrowthRateOfOperatingProfit,omitempty"`
	GrowthRateOfNetAsset        float64 `json:"GrowthRateOfNetAsset,omitempty"`

	CapitalAccumulationFundPerShare            float64 `json:"CapitalAccumulationFundPerShare,omitempty"`
	UndistributedProfitPerShare                float64 `json:"UndistributedProfitPerShare,omitempty"`
	NetAssetValuePerShare                      float64 `json:"NetAssetValuePerShare,omitempty"`
	OperatingCashFlowPerShare                  float64 `json:"OperatingCashFlowPerShare,omitempty"`
	GrowthOfNetCashFlowFromOperatingActivities float64 `json:"GrowthOfNetCashFlowFromOperatingActivities,omitempty"`

	DeductionOfNonNetProfit      float64 `json:"DeductionOfNonNetProfit,omitempty"`
	AssetLiabilityRatio          float64 `json:"AssetLiabilityRatio,omitempty"`
	CurrentRatio                 float64 `json:"CurrentRatio,omitempty"`
	QuickRatio                   float64 `json:"QuickRatio,omitempty"`
	InventoryTurnover            float64 `json:"InventoryTurnover,omitempty"`
	TurnoverOfCurrentAsset       float64 `json:"TurnoverOfCurrentAsset,omitempty"`
	TurnoverRateOfFixedAsset     float64 `json:"TurnoverRateOfFixedAsset,omitempty"`
	TurnoverOfTotalAsset         float64 `json:"TurnoverOfTotalAsset,omitempty"`
	GrowthRateOfCashFlowPerShare float64 `json:"GrowthRateOfCashFlowPerShare,omitempty"`

	CostProfitMargin   float64 `json:"CostProfitMargin,omitempty"`
	ReturnOnTotalAsset float64 `json:"ReturnOnTotalAsset,omitempty"`
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
