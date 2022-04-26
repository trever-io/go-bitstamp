package bitstamp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const INTERNAL_API = "https://www.bitstamp.net/api-internal/market/popular-assets/"

type mostPopular struct {
	Data struct {
		LastListed []string `json:"lastListed"`
		MostViewed []struct {
			Asset string `json:"asset"`
		} `json:"mostViewed"`
	}
}

func (c *Client) Assets() ([]string, error) {
	resp, err := http.Get(INTERNAL_API)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := new(mostPopular)
	err = json.Unmarshal(b, tmp)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(tmp.Data.MostViewed))
	for i, v := range tmp.Data.MostViewed {
		result[i] = v.Asset
	}

	return result, nil
}
