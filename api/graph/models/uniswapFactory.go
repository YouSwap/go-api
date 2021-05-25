package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

var (
	UniswapFactoriesList map[string]*UniswapFactory
)

func init() {
	UniswapFactoriesList = make(map[string]*UniswapFactory)
}

type UniswapFactoryGraph struct {
	Id                 graphql.ID
	PairCount          int64          `graphql:"pairCount"`
	TotalVolumeUSD     graphql.String `graphql:"totalVolumeUSD"`
	TotalVolumeETH     graphql.String `graphql:"totalVolumeETH"`
	UntrackedVolumeUSD graphql.String `graphql:"untrackedVolumeUSD"`
	TotalLiquidityUSD  graphql.String `graphql:"totalLiquidityUSD"`
	TotalLiquidityETH  graphql.String `graphql:"totalLiquidityETH"`
	TxCount            graphql.String `graphql:"txCount"`
}
type UniswapFactory struct {
	Id                 string
	PairCount          int64
	TotalVolumeUSD     float64 `json:"totalVolumeUSD"`
	TotalVolumeETH     float64 `json:"totalVolumeETH"`
	UntrackedVolumeUSD float64 `json:"untrackedVolumeUSD"`
	TotalLiquidityUSD  float64 `json:"totalLiquidityUSD"`
	TotalLiquidityETH  float64 `json:"totalLiquidityETH"`
	TxCount            int64   `json:"txCount"`
	Chain              string  `json:"chain"`
}

func GetAllUniswapFactories() map[string]*UniswapFactory {
	return UniswapFactoriesList
}

func GetUniswapFactory(id string) (u *UniswapFactory, err error) {
	if b, ok := UniswapFactoriesList[id]; ok {
		return b, nil
	}
	return nil, errors.New("bundle not exists")
}
