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
const FIAT_WITHDRAW_ENDPOINT = "v2/withdrawal/open/"
const WITHDRAWAL_REQUEST_ENDPOINT = "v2/withdrawal-requests/"

var optionalMapping map[string]string = map[string]string{
	"xlm":  "memo_id",
	"xrp":  "destination_tag",
	"hbar": "memo_id",
}

type WithdrawResponse struct {
	Id     int64                  `json:"id"`
	Status string                 `json:"status"`
	Reason map[string]interface{} `json:"reason"`
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

func (c *Client) FiatWithdraw(asset string, amount string, name string, iban string, bic string, address string, postalCode string, city string, country string, fiatType string) (*WithdrawResponse, error) {

	data := url.Values{}
	data.Set("amount", amount)
	data.Set("account_currency", asset)
	data.Set("name", name)
	data.Set("iban", iban)
	data.Set("bic", bic)
	data.Set("address", address)
	data.Set("postal_code", postalCode)
	data.Set("city", city)
	data.Set("country", country)
	data.Set("type", fiatType)

	b, err := c.privateRequest(FIAT_WITHDRAW_ENDPOINT, http.MethodPost, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	resp := new(WithdrawResponse)
	err = json.Unmarshal(b, resp)

	if resp.Status == "error" {
		for k, v := range resp.Reason {
			err = fmt.Errorf("%s: %s", k, v)
		}
		return nil, err
	}

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
