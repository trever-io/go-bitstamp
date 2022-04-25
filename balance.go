package bitstamp

import (
	"encoding/json"
	"net/http"
	"strings"
)

const BALANCE_ENDOINT = "v2/balance/"

type AssetBalance struct {
	Available string `json:"available"`
	Balance   string `json:"balance"`
	Reserved  string `json:"reserved"`
}

type Balance map[string]*AssetBalance

func (c *Client) Balance() (Balance, error) {
	b, err := c.privateRequest(BALANCE_ENDOINT, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}

	balance := Balance{}
	tmp := make(map[string]string)
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		return nil, err
	}

	for key, val := range tmp {
		parts := strings.Split(key, "_")
		if len(parts) <= 1 {
			continue
		}
		asset := parts[0]

		if strings.HasSuffix(key, "_balance") {
			if _, ok := balance[asset]; !ok {
				balance[asset] = &AssetBalance{}
			}

			balance[asset].Balance = val
		}

		if strings.HasSuffix(key, "_available") {
			if _, ok := balance[asset]; !ok {
				balance[asset] = &AssetBalance{}
			}

			balance[asset].Available = val
		}

		if strings.HasSuffix(key, "_reserved") {
			if _, ok := balance[asset]; !ok {
				balance[asset] = &AssetBalance{}
			}

			balance[asset].Reserved = val
		}
	}

	return balance, nil
}
