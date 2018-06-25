package exchange

import (
	"math"
)

type Exchange struct {
	Prices map[string][]Price

	markets []Market
}

func MakeExchange() *Exchange {
	return &Exchange{
		markets: make([]Market, 0),
	}
}

func (e *Exchange) AddMarket(m Market) {
	e.markets = append(e.markets, m)
}

func (e *Exchange) UpdatePrices() {
	e.Prices = make(map[string][]Price)

	for _, market := range e.markets {
		marketPrices := market.GetPrices()
		for currenciesKey, price := range marketPrices {
			if _, exists := e.Prices[currenciesKey]; !exists {
				e.Prices[currenciesKey] = make([]Price, 0)
			}

			e.Prices[currenciesKey] = append(e.Prices[currenciesKey], price)
		}
	}
}

func (e *Exchange) FindProfitCurrPairs() []string {
	hashResult := make(map[string]string, 0)

	for currPair, prices := range e.Prices {
		if len(prices) > 1 {
			for currIndex, currPrice := range prices {
				for comparsionIndex, comparsionPrice := range prices {
					if comparsionIndex == currIndex {
						continue
					}

					if currPrice.Buy < comparsionPrice.Sell {
						percentDiff := math.Abs(currPrice.Buy-comparsionPrice.Sell) / ((currPrice.Buy + comparsionPrice.Sell) * 2) * 100
						if percentDiff > 1 {
							hashResult[currPair] = currPair
						}
					}

					if currPrice.Sell > comparsionPrice.Buy {
						percentDiff := math.Abs(currPrice.Sell-comparsionPrice.Buy) / ((currPrice.Sell + comparsionPrice.Buy) * 2) * 100
						if percentDiff > 1 {
							hashResult[currPair] = currPair
						}
					}
				}
			}
		}
	}

	result := make([]string, 0, len(hashResult))
	for k := range hashResult {
		result = append(result, k)
	}

	return result
}

type Market interface {
	GetPrices() map[string]Price
}

type Price struct {
	Buy    float64
	Sell   float64
	Source string
}
