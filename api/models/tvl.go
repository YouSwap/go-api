package models

import (
	"errors"
)

var (
	TvlList map[string]*Tvl
)

func init() {
	TvlList = make(map[string]*Tvl)
}

type Tvl struct {
	TotalAmount float64
}

func GetAllTvls() map[string]*Tvl {
	return TvlList
}

func GetTvl(id string) (u *Tvl, err error) {
	if t, ok := TvlList[id]; ok {
		return t, nil
	}
	return nil, errors.New("tvl not exists")
}

// 矿池总抵押计算：reserveUSD / totalSupply * staketotaldnow
// 单币种抵押计算：youPrice * staketotaldnow
