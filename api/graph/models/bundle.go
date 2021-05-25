package models

import (
	"errors"
	graphql "v1-go-api/graph"
)

type BundleGraph struct {
	Id       graphql.ID
	EthPrice graphql.String
}
type Bundle struct {
	Id       string
	EthPrice float64
	Chain    string
}

var (
	BundleList  map[string]map[string]*Bundle
	bundleCache map[string]*Bundle
)

func init() {
	BundleList = make(map[string]map[string]*Bundle)
	bundleCache = make(map[string]*Bundle)
}

func GetAllBundles() map[string]map[string]*Bundle {
	return BundleList
}

func GetBundle(chain, id string) (u *Bundle, err error) {
	if b, ok := BundleList[chain][id]; ok {
		return b, nil
	}
	return nil, errors.New("bundle not exists")
}

//func GetCacheBundle()  {
//
//}
