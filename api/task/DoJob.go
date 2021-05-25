package task

import (
	"context"
	"github.com/beego/beego/v2/core/logs"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"sync"
	graphql "v1-go-api/graph"
	"v1-go-api/graph/models"
	m2 "v1-go-api/models"
)

var (
	QUERY_SWAP struct {
		Bundles         []*models.BundleGraph         `graphql:"bundles"`
		BuyBacks        []*models.BuyBackGraph        `graphql:"buyBacks"`
		Pairs           []*models.PairGraph           `graphql:"pairs(first: 1000)"`
		Tokens          []*models.TokenGraph          `graphql:"tokens(first: 1000)"`
		UniswapFactoris []*models.UniswapFactoryGraph `graphql:"uniswapFactories"`
		BuyBackList     []*models.SwapGraph           `graphql:"swaps(orderBy: timestamp orderDirection: desc first: 5 where:{to: $BLACKHOLD_ADDR})"`
	}
	QUERY_POOL struct {
		Pools []*models.PoolGraph `graphql:"pools"`
	}
	clients = make(map[string]*graphql.Client)
)

var lock sync.Mutex
var l sync.Mutex
var ctx = context.Background()

func Do() error {
	var err error = nil
	l.Lock()
	for k, v := range GraphDataList {
		//logger.Info("%s start", k)
		if err = doSwap(ctx, v); err != nil {
		}  else {
			err = doPool(ctx, v)
		}
		//logger.Info("%s end", k)
		_ = k
	}
	l.Unlock()
	return err
}
func doSwap(ctx context.Context, graphData *models.GraphData) error {
	lock.Lock()         // 加锁
	defer lock.Unlock() // 解锁
	YOU_ADDRESS := graphData.YouContract
	BLACKHOLD_ADDR := graphData.BlackholdContract

	var swapClients *graphql.Client
	if swapClients = clients[graphData.GraphSwap]; swapClients == nil {
		c := graphql.NewClient(graphData.GraphSwap, nil)
		clients[graphData.GraphSwap] = c
		swapClients = c
	}

	variables := map[string]interface{}{
		"BLACKHOLD_ADDR": BLACKHOLD_ADDR,
	}
	err := swapClients.Query(ctx, &QUERY_SWAP, variables)
	if err != nil {
		logs.Error("query object err %v", err)
		return err
	}

	for _, bundle := range QUERY_SWAP.Bundles {
		id := bundle.Id.(string)
		ethPrice, err := strconv.ParseFloat(string(bundle.EthPrice), 64)
		if err != nil {
			logs.Error("ethPrice err %v", err)
			continue
		}
		// 初始化
		if models.BundleList[graphData.Name] == nil {
			models.BundleList[graphData.Name] = make(map[string]*models.Bundle)
		}

		models.BundleList[graphData.Name][id] = &models.Bundle{
			Id:       id,
			EthPrice: ethPrice,
		}
	}

	ethPrice := models.BundleList[graphData.Name][graphData.BundleId].EthPrice
	for _, token := range QUERY_SWAP.Tokens {
		derivedEth := decimal.RequireFromString(string(token.DerivedETH))
		ethPrice := decimal.NewFromFloat(ethPrice)

		tokenPrice, exact := derivedEth.Mul(ethPrice).Float64()
		_ = exact

		id := token.Id.(string)
		symbol := string(token.Symbol)
		name := string(token.Name)
		decimals, err := strconv.ParseInt(string(token.Decimals), 10, 64)
		if err != nil {
			logs.Error("decimals err %v", err)
			continue
		}

		tradeVolumeStr := decimal.RequireFromString(string(token.TradeVolume))
		tradeVolume, err := strconv.ParseFloat(tradeVolumeStr.String(), 64)
		if err != nil {
			logs.Error("tradeVolume err %v", err)
			continue
		}

		totalSupply, err := strconv.ParseFloat(tradeVolumeStr.String(), 64)
		if err != nil {
			logs.Error("TradeVolume err %v", err)
			continue
		}

		tradeVolumeUSDStr := decimal.RequireFromString(string(token.TradeVolumeUSD))
		tradeVolumeUSD, err := strconv.ParseFloat(tradeVolumeUSDStr.String(), 64)
		if err != nil {
			logs.Error("TradeVolumeUSD err %v", err)
			continue
		}

		untrackedVolumeUSDStr := decimal.RequireFromString(string(token.UntrackedVolumeUSD))
		untrackedVolumeUSD, err := strconv.ParseFloat(untrackedVolumeUSDStr.String(), 64)
		if err != nil {
			logs.Error("UntrackedVolumeUSD err %v", err)
			continue
		}

		txCountStr := decimal.RequireFromString(string(token.TxCount))
		txCount, err := strconv.ParseInt(txCountStr.String(), 10, 64)

		totalLiquidityStr := decimal.RequireFromString(string(token.TotalLiquidity))
		totalLiquidity, err := strconv.ParseFloat(totalLiquidityStr.String(), 64)
		if err != nil {
			logs.Error("totalLiquidity err %v", err)
			continue
		}
		derivedETHStr := decimal.RequireFromString(string(token.DerivedETH))
		derivedETH, err := strconv.ParseFloat(derivedETHStr.String(), 64)
		if err != nil {
			logs.Error("totalLiquidity err %v", err)
			continue
		}
		if models.TokenList[graphData.Name] == nil {
			models.TokenList[graphData.Name] = make(map[string]*models.Token)
		}

		models.TokenList[graphData.Name][id] = &models.Token{
			Id:                 id,
			Symbol:             symbol,
			Name:               name,
			Decimals:           decimals,
			TotalSupply:        totalSupply,
			TradeVolume:        tradeVolume,
			TradeVolumeUSD:     tradeVolumeUSD,
			UntrackedVolumeUSD: untrackedVolumeUSD,
			TxCount:            txCount,
			TotalLiquidity:     totalLiquidity,
			DerivedETH:         derivedETH,
			Price:              tokenPrice,
		}
		if models.Cache.YOU[graphData.Name] == nil {
			models.Cache.YOU[graphData.Name] = make(map[string]interface{})
		}
		if id == YOU_ADDRESS {
			models.Cache.YOU[graphData.Name]["youPrice"] = tokenPrice
		}
	}
	for i, buy := range QUERY_SWAP.BuyBacks {
		id := buy.Id.(string)
		totalUSD, err := strconv.ParseFloat(string(buy.TotalUSD), 64)
		if err != nil {
			logs.Error("totalUSD err %v", err)
			continue
		}
		totalToken, err := strconv.ParseFloat(string(buy.TotalToken), 64)
		if err != nil {
			logs.Error("totalToken err %v", err)
			continue
		}
		if models.BuyBackList[graphData.Name] == nil {
			models.BuyBackList[graphData.Name] = make(map[string]*models.BuyBack)
		}
		models.BuyBackList[graphData.Name][id] = &models.BuyBack{
			Id:         id,
			TotalUSD:   totalUSD,
			TotalToken: totalToken,
		}
		if i == 0 {
			timeRate := totalToken / models.TokenList[graphData.Name][YOU_ADDRESS].TotalSupply
			models.Cache.YOU[graphData.Name]["buyBackFee"] = totalUSD
			models.Cache.YOU[graphData.Name]["buyBackNum"] = totalToken
			models.Cache.YOU[graphData.Name]["timeRate"] = timeRate
		}
	}

	for _, pair := range QUERY_SWAP.Pairs {
		id := pair.Id.(string)
		token0 := models.TokenList[graphData.Name][pair.Token0.Id.(string)]
		token1 := models.TokenList[graphData.Name][pair.Token1.Id.(string)]
		reserve0, err := strconv.ParseFloat(string(pair.Reserve0), 64)
		if err != nil {
			logs.Error("reserve0 err %v", err)
			continue
		}
		reserve1, err := strconv.ParseFloat(string(pair.Reserve1), 64)
		if err != nil {
			logs.Error("reserve1 err %v", err)
			continue
		}

		totalSupply, err := strconv.ParseFloat(string(pair.TotalSupply), 64)
		if err != nil {
			logs.Error("totalSupply err %v", err)
			continue
		}
		reserveETH, err := strconv.ParseFloat(string(pair.ReserveETH), 64)
		if err != nil {
			logs.Error("reserveETH err %v", err)
			continue
		}
		reserveUSD, err := strconv.ParseFloat(string(pair.ReserveUSD), 64)
		if err != nil {
			logs.Error("reserveUSD err %v", err)
			continue
		}
		trackedReserveETH, err := strconv.ParseFloat(string(pair.TrackedReserveETH), 64)
		if err != nil {
			logs.Error("trackedReserveETH err %v", err)
			continue
		}
		token0Price, err := strconv.ParseFloat(string(pair.Token0Price), 64)
		if err != nil {
			logs.Error("token0Price err %v", err)
			continue
		}
		token1Price, err := strconv.ParseFloat(string(pair.Token1Price), 64)
		if err != nil {
			logs.Error("token1Price err %v", err)
			continue
		}
		volumeToken0, err := strconv.ParseFloat(string(pair.VolumeToken0), 64)
		if err != nil {
			logs.Error("volumeToken0 err %v", err)
			continue
		}
		volumeToken1, err := strconv.ParseFloat(string(pair.VolumeToken1), 64)
		if err != nil {
			logs.Error("volumeToken1 err %v", err)
			continue
		}
		volumeUSD, err := strconv.ParseFloat(string(pair.VolumeUSD), 64)
		if err != nil {
			logs.Error("volumeUSD err %v", err)
			continue
		}
		untrackedVolumeUSD, err := strconv.ParseFloat(string(pair.UntrackedVolumeUSD), 64)
		if err != nil {
			logs.Error("untrackedVolumeUSD err %v", err)
			continue
		}
		txCount, err := strconv.ParseInt(string(pair.TxCount), 10, 64)
		if err != nil {
			logs.Error("txCount err %v", err)
			continue
		}
		createdAtTimestamp, err := strconv.ParseInt(string(pair.CreatedAtTimestamp), 10, 64)
		if err != nil {
			logs.Error("createdAtTimestamp err %v", err)
			continue
		}
		createdAtBlockNumber, err := strconv.ParseInt(string(pair.CreatedAtBlockNumber), 10, 64)
		if err != nil {
			logs.Error("createdAtBlockNumber err %v", err)
			continue
		}
		liquidityProviderCount, err := strconv.ParseInt(string(pair.LiquidityProviderCount), 10, 64)
		if err != nil {
			logs.Error("liquidityProviderCount err %v", err)
			continue
		}
		if models.PairList[graphData.Name] == nil {
			models.PairList[graphData.Name] = make(map[string]*models.Pair)
		}
		models.PairList[graphData.Name][id] = &models.Pair{
			Id:                     id,
			Token0:                 *token0,
			Token1:                 *token1,
			Reserve0:               reserve0,
			Reserve1:               reserve1,
			TotalSupply:            totalSupply,
			ReserveETH:             reserveETH,
			ReserveUSD:             reserveUSD,
			TrackedReserveETH:      trackedReserveETH,
			Token0Price:            token0Price,
			Token1Price:            token1Price,
			VolumeToken0:           volumeToken0,
			VolumeToken1:           volumeToken1,
			VolumeUSD:              volumeUSD,
			UntrackedVolumeUSD:     untrackedVolumeUSD,
			TxCount:                txCount,
			CreatedAtTimestamp:     createdAtTimestamp,
			CreatedAtBlockNumber:   createdAtBlockNumber,
			LiquidityProviderCount: liquidityProviderCount,
		}
	}

	for _, factory := range QUERY_SWAP.UniswapFactoris {
		id := factory.Id.(string)

		//pairCountStr := decimal.RequireFromString(string(factory.PairCount))
		//pairCount, err := strconv.ParseInt(pairCountStr.String(), 10, 64)
		//pairCount, err := strconv.ParseInt(string(factory.PairCount),10,64)

		pairCount := factory.PairCount
		if err != nil {
			logs.Error("pairCount err %v", err)
			continue
		}
		totalVolumeUSD, err := strconv.ParseFloat(string(factory.TotalVolumeUSD), 64)
		if err != nil {
			logs.Error("totalVolumeUSD err %v", err)
			continue
		}
		totalVolumeETH, err := strconv.ParseFloat(string(factory.TotalVolumeETH), 64)
		if err != nil {
			logs.Error("totalVolumeETH err %v", err)
			continue
		}
		untrackedVolumeUSD, err := strconv.ParseFloat(string(factory.UntrackedVolumeUSD), 64)
		if err != nil {
			logs.Error("untrackedVolumeUSD err %v", err)
			continue
		}
		totalLiquidityUSD, err := strconv.ParseFloat(string(factory.TotalLiquidityUSD), 64)
		if err != nil {
			logs.Error("totalLiquidityUSD err %v", err)
			continue
		}
		totalLiquidityETH, err := strconv.ParseFloat(string(factory.TotalLiquidityETH), 64)
		if err != nil {
			logs.Error("totalLiquidityETH err %v", err)
			continue
		}
		txCountStr := decimal.RequireFromString(string(factory.TxCount))
		txCount, err := strconv.ParseInt(txCountStr.String(), 10, 64)
		//txCount, err := strconv.ParseInt(string(factory.TxCount),10,64)
		if err != nil {
			logs.Error("txCount err %v", err)
			continue
		}

		models.UniswapFactoriesList[id+"_"+graphData.Name] = &models.UniswapFactory{
			Id:                 id,
			PairCount:          pairCount,
			TotalVolumeUSD:     totalVolumeUSD,
			TotalVolumeETH:     totalVolumeETH,
			UntrackedVolumeUSD: untrackedVolumeUSD,
			TotalLiquidityUSD:  totalLiquidityUSD,
			TotalLiquidityETH:  totalLiquidityETH,
			TxCount:            txCount,
		}
	}

	for i, swap := range QUERY_SWAP.BuyBackList {
		id := swap.Id.(string)
		sender := string(swap.Sender)
		from := string(swap.From)
		to := string(swap.To)

		logIndex, err := strconv.ParseInt(string(swap.LogIndex), 10, 64)
		if err != nil {
			logs.Error("logIndex err %v", err)
			continue
		}
		timestamp, err := strconv.ParseInt(string(swap.Timestamp), 10, 64)
		if err != nil {
			logs.Error("timestamp err %v", err)
			continue
		}

		amount0In, err := strconv.ParseFloat(string(swap.Amount0In), 64)
		if err != nil {
			logs.Error("amount0In err %v", err)
			continue
		}
		amount1In, err := strconv.ParseFloat(string(swap.Amount1In), 64)
		if err != nil {
			logs.Error("amount1In err %v", err)
			continue
		}
		amount0Out, err := strconv.ParseFloat(string(swap.Amount0Out), 64)
		if err != nil {
			logs.Error("amount0Out err %v", err)
			continue
		}
		amount1Out, err := strconv.ParseFloat(string(swap.Amount1Out), 64)
		if err != nil {
			logs.Error("amount1Out err %v", err)
			continue
		}
		amountUSD, err := strconv.ParseFloat(string(swap.AmountUSD), 64)
		if err != nil {
			logs.Error("amountUSD err %v", err)
			continue
		}

		tran := swap.Transaction
		id_t := tran.Id.(string)
		blockNumber_t, err := strconv.ParseInt(string(tran.BlockNumber), 10, 64)
		if err != nil {
			logs.Error("blockNumber_t err %v", err)
			continue
		}
		timestamp_t, err := strconv.ParseInt(string(tran.Timestamp), 10, 64)
		if err != nil {
			logs.Error("timestamp_t err %v", err)
			continue
		}
		//mints := string(tran.Mints)
		//burns := string(tran.Burns)
		//swaps := string(tran.Swaps)

		t := &models.Transaction{
			Id:          id_t,
			BlockNumber: blockNumber_t,
			Timestamp:   timestamp_t,
			//Mints:       mints,
			//Burns:       burns,
			//Swaps:       swaps,
		}

		pair := swap.Pair

		token0 := models.TokenList[graphData.Name][pair.Token0.Id.(string)]
		token1 := models.TokenList[graphData.Name][pair.Token1.Id.(string)]

		reserve0, err := strconv.ParseFloat(string(pair.Reserve0), 64)
		if err != nil {
			logs.Error("reserve0 err %v", err)
			continue
		}
		reserve1, err := strconv.ParseFloat(string(pair.Reserve1), 64)
		if err != nil {
			logs.Error("reserve1 err %v", err)
			continue
		}

		totalSupply, err := strconv.ParseFloat(string(pair.TotalSupply), 64)
		if err != nil {
			logs.Error("totalSupply err %v", err)
			continue
		}
		reserveETH, err := strconv.ParseFloat(string(pair.ReserveETH), 64)
		if err != nil {
			logs.Error("reserveETH err %v", err)
			continue
		}
		reserveUSD, err := strconv.ParseFloat(string(pair.ReserveUSD), 64)
		if err != nil {
			logs.Error("reserveUSD err %v", err)
			continue
		}
		trackedReserveETH, err := strconv.ParseFloat(string(pair.TrackedReserveETH), 64)
		if err != nil {
			logs.Error("trackedReserveETH err %v", err)
			continue
		}
		token0Price, err := strconv.ParseFloat(string(pair.Token0Price), 64)
		if err != nil {
			logs.Error("token0Price err %v", err)
			continue
		}
		token1Price, err := strconv.ParseFloat(string(pair.Token1Price), 64)
		if err != nil {
			logs.Error("token1Price err %v", err)
			continue
		}
		volumeToken0, err := strconv.ParseFloat(string(pair.VolumeToken0), 64)
		if err != nil {
			logs.Error("volumeToken0 err %v", err)
			continue
		}
		volumeToken1, err := strconv.ParseFloat(string(pair.VolumeToken1), 64)
		if err != nil {
			logs.Error("volumeToken1 err %v", err)
			continue
		}
		volumeUSD, err := strconv.ParseFloat(string(pair.VolumeUSD), 64)
		if err != nil {
			logs.Error("volumeUSD err %v", err)
			continue
		}
		untrackedVolumeUSD, err := strconv.ParseFloat(string(pair.UntrackedVolumeUSD), 64)
		if err != nil {
			logs.Error("untrackedVolumeUSD err %v", err)
			continue
		}
		txCount, err := strconv.ParseInt(string(pair.TxCount), 10, 64)
		if err != nil {
			logs.Error("txCount err %v", err)
			continue
		}
		createdAtTimestamp, err := strconv.ParseInt(string(pair.CreatedAtTimestamp), 10, 64)
		if err != nil {
			logs.Error("createdAtTimestamp err %v", err)
			continue
		}
		createdAtBlockNumber, err := strconv.ParseInt(string(pair.CreatedAtBlockNumber), 10, 64)
		if err != nil {
			logs.Error("createdAtBlockNumber err %v", err)
			continue
		}
		liquidityProviderCount, err := strconv.ParseInt(string(pair.LiquidityProviderCount), 10, 64)
		if err != nil {
			logs.Error("liquidityProviderCount err %v", err)
			continue
		}
		p := &models.Pair{
			Id:                     id,
			Token0:                 *token0,
			Token1:                 *token1,
			Reserve0:               reserve0,
			Reserve1:               reserve1,
			TotalSupply:            totalSupply,
			ReserveETH:             reserveETH,
			ReserveUSD:             reserveUSD,
			TrackedReserveETH:      trackedReserveETH,
			Token0Price:            token0Price,
			Token1Price:            token1Price,
			VolumeToken0:           volumeToken0,
			VolumeToken1:           volumeToken1,
			VolumeUSD:              volumeUSD,
			UntrackedVolumeUSD:     untrackedVolumeUSD,
			TxCount:                txCount,
			CreatedAtTimestamp:     createdAtTimestamp,
			CreatedAtBlockNumber:   createdAtBlockNumber,
			LiquidityProviderCount: liquidityProviderCount,
		}
		models.Cache.BuyBacks[i] = &models.Swap{
			Id:          id,
			Transaction: t,
			Timestamp:   timestamp,
			Pair:        p,
			Sender:      sender,
			From:        from,
			Amount0In:   amount0In,
			Amount1In:   amount1In,
			Amount0Out:  amount0Out,
			Amount1Out:  amount1Out,
			To:          to,
			LogIndex:    logIndex,
			AmountUSD:   amountUSD,
		}
	}
	return nil
}

func doPool(ctx context.Context, graphData *models.GraphData) error {
	lock.Lock()         // 加锁
	defer lock.Unlock() // 解锁

	variables := map[string]interface{}{}

	var poolClients *graphql.Client
	if poolClients = clients[graphData.GraphPool]; poolClients == nil {
		poolClients = graphql.NewClient(graphData.GraphPool, nil)
		clients[graphData.GraphPool] = poolClients
	}

	err := poolClients.Query(ctx, &QUERY_POOL, variables)
	if err != nil {
		logs.Error("query object err %v", err)
		return err
	}

	HecoTvlAmount := decimal.NewFromFloat(0)
	YouTvlAmount := decimal.NewFromFloat(0)
	for _, pool := range QUERY_POOL.Pools {
		id := pool.Id.(string)
		poolPool := string(pool.Pool)
		lpAddress := string(pool.Lpaddress)
		poolName := string(pool.Poolname)

		stakeTotalNow, err := strconv.ParseFloat(string(pool.Staketotaldnow), 64)
		if err != nil {
			logs.Error("stakeTotalNow err %v", err)
			continue
		}
		types, err := strconv.ParseFloat(string(pool.Type), 64)
		if err != nil {
			logs.Error("types err %v", err)
			continue
		}
		if types == 2 && stakeTotalNow > 0 {
			Decimals := models.TokenList[graphData.Name][lpAddress].Decimals
			price := models.TokenList[graphData.Name][lpAddress].Price
			tvlAmount := price * stakeTotalNow
			TokenTvlAmount := decimal.NewFromFloat(tvlAmount).Div(decimal.NewFromFloat(math.Pow10(int(Decimals))))
			YouTvlAmount = TokenTvlAmount.Add(YouTvlAmount)
			//logs.Info("YouTvlAmount==", YouTvlAmount)
		} else {
			if models.PairList[graphData.Name] == nil {
				models.PairList[graphData.Name] = make(map[string]*models.Pair)
			}
			if stakeTotalNow > 0 {
				if models.PairList[graphData.Name][lpAddress] == nil {
					continue
				}
				//reserveUSD / totalSupply * staketotaldnow
				reserveUSD := models.PairList[graphData.Name][lpAddress].ReserveUSD
				totalSupply := models.PairList[graphData.Name][lpAddress].TotalSupply
				tvlAmount := reserveUSD / totalSupply * stakeTotalNow
				HecoTvlAmount = HecoTvlAmount.Add(decimal.NewFromFloat(tvlAmount))
				//amount := decimal.NewFromFloat(tvlAmount).Div(decimal.NewFromFloat(math.Pow10(18)))
				//logs.Info(graphData.Name+"===lpAddress======", lpAddress+"=====tvlAmount======", amount)
			} else {
				//logs.Info("0000000000000000000")
			}
		}

		pair := models.PairList[graphData.Name][lpAddress]
		if pair == nil {
			continue
		}

		startBlockHeight, err := strconv.ParseInt(string(pool.Startblockheight), 10, 64)
		if err != nil {
			logs.Error("startBlockHeight err %v", err)
			continue
		}
		rewardTotal, err := strconv.ParseFloat(string(pool.Rewardtotal), 64)
		if err != nil {
			logs.Error("rewardTotal err %v", err)
			continue
		}
		rewardPerBlock, err := strconv.ParseInt(string(pool.Rewardperblock), 10, 64)
		if err != nil {
			logs.Error("rewardPerBlock err %v", err)
			continue
		}
		rewardMultiple, err := strconv.ParseFloat(string(pool.Rewardmultiple), 64)
		if err != nil {
			logs.Error("rewardMultiple err %v", err)
			continue
		}
		priority, err := strconv.ParseInt(string(pool.Priority), 10, 64)
		if err != nil {
			logs.Error("priority err %v", err)
			continue
		}

		rewardCanWithdrawNow, err := strconv.ParseFloat(string(pool.Rewardcanwithdrawnow), 64)
		if err != nil {
			logs.Error("rewardCanWithdrawNow err %v", err)
			continue
		}
		totalPower, err := strconv.ParseFloat(string(pool.Totalpower), 64)
		if err != nil {
			//logs.Error("totalPower err %v", err)
			//continue
			totalPower = 0
		}
		isFinshed := bool(pool.Isfinshed)

		models.PoolList[id+"_"+graphData.Name] = &models.Pool{
			Id:                   id,
			Pool:                 poolPool,
			Lpaddress:            lpAddress,
			Poolname:             poolName,
			Startblockheight:     startBlockHeight,
			Rewardtotal:          rewardTotal,
			Rewardperblock:       rewardPerBlock,
			Rewardmultiple:       rewardMultiple,
			Priority:             priority,
			Isfinshed:            isFinshed,
			Staketotaldnow:       stakeTotalNow,
			Rewardcanwithdrawnow: rewardCanWithdrawNow,
			Totalpower:           totalPower,
		}
	}
	tvlAmount := HecoTvlAmount.Div(decimal.NewFromFloat(math.Pow10(18))).Add(YouTvlAmount)

	_tvlAmount, err := strconv.ParseFloat(string(tvlAmount.String()), 64)
	m2.TvlList[graphData.Name] = &m2.Tvl{
		TotalAmount: _tvlAmount,
	}
	return nil
}
