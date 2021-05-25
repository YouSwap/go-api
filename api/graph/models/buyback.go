package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

type BuyBackGraph struct {
	Id         graphql.ID
	TotalUSD   graphql.String `graphql:"totalUSD" json:"totalUSD"`
	TotalToken graphql.String
}

type BuyBack struct {
	Id         string
	TotalUSD   float64 `graphql:"totalUSD" json:"totalUSD"`
	TotalToken float64
	Chain      string
}

var (
	BuyBackList  map[string]map[string]*BuyBack
	buyBackCache map[string]*BuyBack
)

func init() {
	BuyBackList = make(map[string]map[string]*BuyBack)
	buyBackCache = make(map[string]*BuyBack)
}

func GetAllBuyBacks() map[string]map[string]*BuyBack {
	return BuyBackList
}

func GetBuyBack(id, chain string) (u *BuyBack, err error) {
	if b, ok := BuyBackList[chain][id]; ok {
		return b, nil
	}
	return nil, errors.New("buyback not exists")
}

func GetCacheBuyBack(id string) (b *BuyBack, err error) {
	if buyBack := buyBackCache[id]; buyBack != nil {
		return buyBack, nil
	}
	for k, _ := range BuyBackList {
		if buyback := BuyBackList[k][id]; buyback != nil {
			buyBackCache[id] = buyback
			return buyback, nil
		}
	}
	return nil, errors.New("buyback not exists")
}
