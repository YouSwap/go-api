package models

type GraphData struct {
	Igrone            bool   `json:"igrone"`
	GraphSwap         string `json:"graph_swap"`
	GraphPool         string `json:"graph_pool"`
	GraphBlock        string `json:"graph_block"`
	GraphRelation     string `json:"graph_relation"`
	YouContract       string `json:"you_contract"`
	FactoryContract   string `json:"factory_contract"`
	BlackholdContract string `json:"blackhold_contract"`
	YearBlockNum      int64  `json:"year_block_num"`
	BundleId          string `json:"bundle_id"`
	Name              string `json:"name"`
}
