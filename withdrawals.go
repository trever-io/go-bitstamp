package bitstamp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const WITHDRAW_ENDPOINT = "v2/%v_withdrawal/"

var optionalMapping map[string]string = map[string]string{
	"xlm":  "memo_id",
	"xrp":  "destination_tag",
	"hbar": "memo_id",
}

type WithdrawResponse struct {
	Id int64 `json:"id"`
}

func (c *Client) Withdraw(asset string, address string, amount string, optionalParam string) (*WithdrawResponse, error) {
	data := url.Values{}
	data.Set("address", address)
	data.Set("amount", amount)

	if optionalParam != "" {
		optionalTag, ok := optionalMapping[asset]
		if !ok {
			return nil, fmt.Errorf("no optional tag mapping for %v", asset)
		}

		data.Set(optionalTag, optionalParam)
	}

	b, err := c.privateRequest(fmt.Sprintf(WITHDRAW_ENDPOINT, asset), http.MethodPost, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	resp := new(WithdrawResponse)
	err = json.Unmarshal(b, resp)

	return resp, err
}
