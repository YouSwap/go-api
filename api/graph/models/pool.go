package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

type PoolGraph struct {
	Id                   graphql.ID
	Pool                 graphql.String
	Lpaddress            graphql.String
	Poolname             graphql.String
	Startblockheight     graphql.String
	Rewardtotal          graphql.String
	Rewardperblock       graphql.String
	Rewardmultiple       graphql.String
	Priority             graphql.String
	Isfinshed            graphql.Boolean
	Staketotaldnow       graphql.String
	Rewardcanwithdrawnow graphql.String
	Totalpower           graphql.String
	Type                 graphql.String
}

type Pool struct {
	Id                   string
	Pool                 string
	Lpaddress            string
	Poolname             string
	Startblockheight     int64
	Rewardtotal          float64
	Rewardperblock       int64
	Rewardmultiple       float64
	Priority             int64
	Isfinshed            bool
	Staketotaldnow       float64
	Rewardcanwithdrawnow float64
	Totalpower           float64
	Apy                  float64
	Type                 int
}

var (
	PoolList map[string]*Pool
)

func init() {
	PoolList = make(map[string]*Pool)
}

func GetAllPools() map[string]*Pool {
	return PoolList
}

func GetPool(id string) (u *Pool, err error) {
	if p, ok := PoolList[id]; ok {
		return p, nil
	}
	return nil, errors.New("pair not exists")
}
