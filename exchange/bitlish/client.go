package bitlish

import (
	"arbitrage/exchange"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BitlishAPI struct {
	domain string

	httpClient *http.Client
}

func MakeBitlishAPI() *BitlishAPI {
	return &BitlishAPI{
		domain: "https://bitlish.com/api",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (b *BitlishAPI) GetPrices() map[string]exchange.Price {
	result := make(map[string]exchange.Price, 0)

	response := make(map[string]priceItem, 0)
	err := b.makeGetRequest("/v1/tickers/", &response)
	if err != nil {
		log.Printf("[WARN] bitlish get price request has failed: %s", err.Error())
		return result
	}

	for currPair, price := range response {
		sellCost, parseErr := strconv.ParseFloat(price.Ask, 64)
		buyCost, parseErr := strconv.ParseFloat(price.Bid, 64)
		runeCurrPair := []rune(currPair)
		parsedCurrPair := strings.ToUpper(string(runeCurrPair[:3]) + "_" + string(runeCurrPair[3:]))
		if parseErr != nil {
			log.Printf("[WARN] bitlish price parse has failed: %s, %s", parsedCurrPair, err.Error())
			continue
		}

		result[parsedCurrPair] = exchange.Price{
			Source: "BITLISH",
			Sell:   sellCost,
			Buy:    buyCost,
		}
	}

	return result
}

func (b *BitlishAPI) makeGetRequest(url string, result interface{}) error {
	r, err := b.httpClient.Get(b.domain + url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(result)
}

type priceItem struct {
	Ask     string `json:"ask"`
	Bid     string `json:"bid"`
	First   string `json:"first"`
	Last    string `json:"last"`
	Max     string `json:"max"`
	Min     string `json:"min"`
	Prc     string `json:"prc"`
	Sum     string `json:"sum"`
	Updated string `json:"updated"`
}
