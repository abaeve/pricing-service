package handler

import (
	"context"
	"github.com/abaeve/pricing-service/model"
	"github.com/abaeve/pricing-service/proto"
)

type StatisticsProvider interface {
	Stats(regionId, typeId int32) *model.TypeStat
}

type pricesHandler struct {
	sp StatisticsProvider
}

func (ph *pricesHandler) GetItemPrice(ctx context.Context, req *pricing.ItemPriceRequest, rsp *pricing.ItemPriceResponse) error {
	stats := ph.sp.Stats(req.RegionId, req.ItemId)

	rsp.Item = &pricing.Item{
		ItemId: req.ItemId,
		Buy: &pricing.ItemPrice{
			Min: stats.Buy.Min,
			Avg: stats.Buy.Avg,
			Max: stats.Buy.Max,
			Vol: stats.Buy.Vol,
			Ord: stats.Buy.Ord,
		},
		Sell: &pricing.ItemPrice{
			Min: stats.Sell.Min,
			Avg: stats.Sell.Avg,
			Max: stats.Sell.Max,
			Vol: stats.Sell.Vol,
			Ord: stats.Sell.Ord,
		},
	}

	return nil
}

func NewPricesHandler(sp StatisticsProvider) pricing.PricesHandler {
	return &pricesHandler{
		sp: sp,
	}
}
