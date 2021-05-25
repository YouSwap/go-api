package models

var (
	Cache = &struct {
		YOU      map[string]map[string]interface{}
		BuyBacks [5]*Swap
		Pools    [5]*Pool
	}{
		YOU:      make(map[string]map[string]interface{}),
		BuyBacks: [5]*Swap{},
		Pools:    [5]*Pool{},
	}
)
