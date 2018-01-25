package kraken

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Beldur/kraken-go-api-client"
	"github.com/lightyeario/kelp/support/exchange"
	"github.com/lightyeario/kelp/support/exchange/assets"
	"github.com/stretchr/testify/assert"
)

var testKrakenExchange exchange.Exchange = krakenExchange{
	assetConverter: assets.KrakenAssetConverter,
	api:            krakenapi.New("", ""),
	delimiter:      "",
	isSimulated:    true,
}

func TestGetTickerPrice(t *testing.T) {
	pair := assets.TradingPair{Base: assets.XLM, Quote: assets.BTC}
	pairs := []assets.TradingPair{pair}

	m, e := testKrakenExchange.GetTickerPrice(pairs)
	if !assert.NoError(t, e) {
		return
	}
	assert.Equal(t, 1, len(m))

	ticker := m[pair]
	assert.True(t, ticker.AskPrice.AsFloat() < 1, ticker.AskPrice.AsString())
}

func TestGetAccountBalances(t *testing.T) {
	assetList := []assets.Asset{
		assets.USD,
		assets.XLM,
		assets.BTC,
		assets.LTC,
		assets.ETH,
		assets.REP,
	}
	m, e := testKrakenExchange.GetAccountBalances(assetList)
	if !assert.NoError(t, e) {
		return
	}
	assert.Equal(t, 6, len(m))

	// print balances here for convenience
	for assetKey, balanceValue := range m {
		fmt.Printf("Balance %s = %.8f\n", assetKey, balanceValue.AsFloat())
	}

	for _, a := range assetList {
		bal := m[a]
		assert.True(t, bal.AsFloat() > 0, bal.AsString())
	}
}

func TestGetOrderBook(t *testing.T) {
	pair := assets.TradingPair{Base: assets.XLM, Quote: assets.BTC}
	ob, e := testKrakenExchange.GetOrderBook(&pair, 10)
	if !assert.NoError(t, e) {
		return
	}

	assert.True(t, len(ob.Asks()) > 0, len(ob.Asks()))
	assert.True(t, len(ob.Bids()) > 0, len(ob.Bids()))
	assert.True(t, ob.Asks()[0].OrderAction.IsSell())
	assert.True(t, ob.Asks()[0].OrderType.IsLimit())
	assert.True(t, ob.Bids()[0].OrderAction.IsBuy())
	assert.True(t, ob.Bids()[0].OrderType.IsLimit())
}

func TestGetTrades(t *testing.T) {
	pair := assets.TradingPair{Base: assets.XLM, Quote: assets.BTC}
	trades, e := testKrakenExchange.GetTrades(&pair, nil)
	if !assert.NoError(t, e) {
		return
	}

	cursor := trades.Cursor.(int64)
	assert.True(t, cursor > 0, strconv.FormatInt(cursor, 10))
	assert.True(t, len(trades.Trades) > 0)
}

func TestGetTradeHistory(t *testing.T) {
	tradeHistoryResult, e := testKrakenExchange.GetTradeHistory(nil, nil)
	if !assert.NoError(t, e) {
		return
	}
	assert.True(t, len(tradeHistoryResult.Trades) > 0)
}

func TestGetOpenOrders(t *testing.T) {
	m, e := testKrakenExchange.GetOpenOrders()
	if !assert.NoError(t, e) {
		return
	}

	// print open orders here for convenience
	for pair, openOrders := range m {
		fmt.Printf("Open Orders for pair: %s\n", pair.String())
		for _, o := range openOrders {
			fmt.Printf("    %s\n", o.String())
		}
	}

	assert.True(t, len(m) > 0, "there were no open orders")
}