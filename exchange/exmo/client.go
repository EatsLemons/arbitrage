package exmo

import (
	"arbitrage/exchange"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ExmoAPI struct {
	domain string

	httpClient *http.Client
}

func MakeExmoAPI() *ExmoAPI {
	return &ExmoAPI{
		domain: "https://api.exmo.com",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (e *ExmoAPI) GetPrices() map[string]exchange.Price {
	result := make(map[string]exchange.Price, 0)

	response := make(map[string]priceItem, 0)
	err := e.makeGetRequest("/v1/ticker/", &response)
	if err != nil {
		log.Printf("[WARN] exmo get price request has failed: %s", err.Error())
		return result
	}

	for currPair, price := range response {
		sellCost, parseErr := strconv.ParseFloat(price.SellPrice, 64)
		buyCost, parseErr := strconv.ParseFloat(price.BuyPrice, 64)
		if parseErr != nil {
			log.Printf("[WARN] exmo price parse has failed: %s, %s", currPair, err.Error())
			continue
		}

		result[currPair] = exchange.Price{
			Source: "EXMO",
			Sell:   sellCost,
			Buy:    buyCost,
		}
	}

	return result
}

func (e *ExmoAPI) makeGetRequest(url string, result interface{}) error {
	r, err := e.httpClient.Get(e.domain + url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(result)
}

type priceItem struct {
	BuyPrice  string `json:"buy_price"`
	SellPrice string `json:"sell_price"`
	LastTrade string `json:"last_trade"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Avg       string `json:"avg"`
	Vol       string `json:"vol"`
	VolCurr   string `json:"vol_curr"`
	Updated   int    `json:"updated"`
}
