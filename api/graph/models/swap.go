package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

var (
	SwapList map[string]*Swap
)

func init() {
	SwapList = make(map[string]*Swap)
}

type SwapGraph struct {
	Id          graphql.ID
	Transaction TransactionGraph
	Timestamp   graphql.String
	Pair        PairGraph
	Sender      graphql.String
	From        graphql.String
	Amount0In   graphql.String `graphql:"amount0In"`
	Amount1In   graphql.String `graphql:"amount1In"`
	Amount0Out  graphql.String `graphql:"amount0Out"`
	Amount1Out  graphql.String `graphql:"amount1Out"`
	To          graphql.String
	LogIndex    graphql.String
	AmountUSD   graphql.String `graphql:"amountUSD"`
}
type Swap struct {
	Id          string
	Transaction *Transaction
	Timestamp   int64
	Pair        *Pair
	Sender      string
	From        string
	Amount0In   float64
	Amount1In   float64
	Amount0Out  float64
	Amount1Out  float64
	To          string
	LogIndex    int64
	AmountUSD   float64 `json:"amountUSD"`
}

func GetAllSwaps() map[string]*Swap {
	return SwapList
}

func GetSwap(id string) (u *Swap, err error) {
	if b, ok := SwapList[id]; ok {
		return b, nil
	}
	return nil, errors.New("swap not exists")
}
