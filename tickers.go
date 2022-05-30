package bitstamp

import (
	"encoding/json"
	"net/http"
)

const TICKER_ENDPOINT = "v2/ticker/"

type TickerResponse struct {
	Last      string `json:"last"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Timestamp string `json:"timestamp"`
	Open      string `json:"open"`
}

func (c *Client) Ticker(currencyPair string) (*TickerResponse, error) {
	b, err := c.publicRequest(TICKER_ENDPOINT+currencyPair+"/", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp := new(TickerResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
