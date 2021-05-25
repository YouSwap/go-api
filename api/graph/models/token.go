package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

type TokenGraph struct {
	Id                 graphql.ID
	Symbol             graphql.String
	Name               graphql.String
	Decimals           graphql.String
	TotalSupply        graphql.String
	TradeVolume        graphql.String
	TradeVolumeUSD     graphql.String `graphql:"tradeVolumeUSD"`
	UntrackedVolumeUSD graphql.String `graphql:"untrackedVolumeUSD"`
	TxCount            graphql.String
	TotalLiquidity     graphql.String
	DerivedETH         graphql.String `graphql:"derivedETH"`
}

type Token struct {
	Id                 string  `json:"id"`
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name"`
	Decimals           int64   `json:"decimals"`
	TotalSupply        float64 `json:"total_supply"`
	TradeVolume        float64 `json:"trade_volume"`
	TradeVolumeUSD     float64 `json:"tradeVolumeUSD"`
	UntrackedVolumeUSD float64 `json:"untrackedVolumeUSD"`
	TxCount            int64   `json:"tx_count"`
	TotalLiquidity     float64 `json:"total_liquidity"`
	DerivedETH         float64 `json:"derivedETH"`
	Price              float64 `json:"price"`
	Chain              string  `json:"chain"`
}

var (
	TokenList  map[string]map[string]*Token
	tokenCache map[string]*Token
)

func init() {
	TokenList = make(map[string]map[string]*Token)
	tokenCache = make(map[string]*Token)
}

func GetAllTokens() map[string]map[string]*Token {
	return TokenList
}

func GetToken(chain, address string) (u *Token, err error) {
	if t, ok := TokenList[chain][address]; ok {
		return t, nil
	}
	return nil, errors.New("token not exists")
}

func GetCacheToken(address string) (t *Token, err error) {
	if token := tokenCache[address]; token != nil {
		return token, nil
	}
	for k, _ := range TokenList {
		token := TokenList[k][address]
		if token != nil {
			tokenCache[address] = token
			return token, nil
		}
	}
	return nil, errors.New("token not exists")
}
