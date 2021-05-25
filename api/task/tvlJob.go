package task

import (
	"v1-go-api/graph/models"
)

var (
	TVL_AMOUNT struct {
		Tvl []*models.TokenGraph `graphql:"tokens"`
	}
)

//func getTvlAmount(client *graphql.Client, ctx context.Context) error {
//	models.TvlList["1"] = &models.Tvl{
//		TotalAmount: 300,
//	}
//	return nil
//}

//func Tvls() error {
//	logs.Debug("get tvl %s", time.Now().String())
//	client := graphql.NewClient(SWAP_URL, nil)
//	ctx := context.Background()
//	return getTvlAmount(client, ctx)
//}
