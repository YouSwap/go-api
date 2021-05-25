package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

type PairGraph struct {
	Id                     graphql.ID
	Token0                 TokenGraph
	Token1                 TokenGraph
	Reserve0               graphql.String
	Reserve1               graphql.String
	TotalSupply            graphql.String
	ReserveETH             graphql.String `graphql:"reserveETH"`
	ReserveUSD             graphql.String `graphql:"reserveUSD"`
	TrackedReserveETH      graphql.String `graphql:"trackedReserveETH"`
	Token0Price            graphql.String `graphql:"token0Price"`
	Token1Price            graphql.String `graphql:"token1Price"`
	VolumeToken0           graphql.String
	VolumeToken1           graphql.String
	VolumeUSD              graphql.String `graphql:"volumeUSD"`
	UntrackedVolumeUSD     graphql.String `graphql:"untrackedVolumeUSD"`
	TxCount                graphql.String
	CreatedAtTimestamp     graphql.String
	CreatedAtBlockNumber   graphql.String
	LiquidityProviderCount graphql.String
}

type Pair struct {
	Id                     string
	Token0                 Token
	Token1                 Token
	Reserve0               float64
	Reserve1               float64
	TotalSupply            float64
	ReserveETH             float64 `graphql:"reserveETH"`
	ReserveUSD             float64 `graphql:"reserveUSD"`
	TrackedReserveETH      float64 `graphql:"trackedReserveETH"`
	Token0Price            float64 `graphql:"token0Price"`
	Token1Price            float64 `graphql:"token1Price"`
	VolumeToken0           float64
	VolumeToken1           float64
	VolumeUSD              float64 `graphql:"volumeUSD"`
	UntrackedVolumeUSD     float64 `graphql:"untrackedVolumeUSD"`
	TxCount                int64
	CreatedAtTimestamp     int64
	CreatedAtBlockNumber   int64
	LiquidityProviderCount int64
	Chain                  string
}

var (
	PairList  map[string]map[string]*Pair
	pairCache map[string]*Pair
)

func init() {
	PairList = make(map[string]map[string]*Pair)
	pairCache = make(map[string]*Pair)
}

func GetAllPairs() map[string]map[string]*Pair {
	return PairList
}

func GetPair(chain, address string) (u *Pair, err error) {
	if p, ok := PairList[chain][address]; ok {
		return p, nil
	}
	return nil, errors.New("pair not exists")
}

func GetCachePair(address string) (*Pair, error) {
	if pair := pairCache[address]; pair != nil {
		return pair, nil
	}

	for k, _ := range PairList {
		pair := PairList[k][address]
		if pair != nil {
			pairCache[address] = pair
			return pair, nil
		}
	}
	return nil, errors.New("pair not exists")
}
