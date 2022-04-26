package bitstamp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const WITHDRAW_ENDPOINT = "v2/%v_withdrawal/"
const WITHDRAWAL_REQUEST_ENDPOINT = "v2/withdrawal-requests/"

var optionalMapping map[string]string = map[string]string{
	"xlm":  "memo_id",
	"xrp":  "destination_tag",
	"hbar": "memo_id",
}

type WithdrawResponse struct {
	Id int64 `json:"id"`
}

type WithdrawalRequest struct {
	Id            int64  `json:"id"`
	Status        int    `json:"status"`
	DateTime      string `json:"datetime"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	Address       string `json:"address"`
	Type          int    `json:"type"`
	TransactionId string `json:"transaction_id"`
}

type WithdrawalRequestsReponse []*WithdrawalRequest

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

func (c *Client) WithdrawalRequests(timedelta int64) (WithdrawalRequestsReponse, error) {
	var buf io.Reader
	data := url.Values{}

	if timedelta != 0 {
		data.Set("timedelta", strconv.FormatInt(timedelta, 10))
	}

	if len(data) > 0 {
		buf = strings.NewReader(data.Encode())
	}

	b, err := c.privateRequest(WITHDRAWAL_REQUEST_ENDPOINT, http.MethodPost, buf)
	if err != nil {
		return nil, err
	}

	resp := WithdrawalRequestsReponse{}
	err = json.Unmarshal(b, &resp)
	return resp, err
}
