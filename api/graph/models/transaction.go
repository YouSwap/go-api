package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

var (
	TransactionList map[string]*Transaction
)

func init() {
	TransactionList = make(map[string]*Transaction)
}

type TransactionGraph struct {
	Id          graphql.ID
	BlockNumber graphql.String
	Timestamp   graphql.String
	//Mints 		graphql.String
	//Burns 		graphql.String
	//Swaps 		graphql.String
}
type Transaction struct {
	Id          string
	BlockNumber int64
	Timestamp   int64
	//Mints 		string
	//Burns 		string
	//Swaps 		string
}

func GetAllTransactions() map[string]*Transaction {
	return TransactionList
}

func GetTransaction(id string) (u *Transaction, err error) {
	if b, ok := TransactionList[id]; ok {
		return b, nil
	}
	return nil, errors.New("transaction not exists")
}
