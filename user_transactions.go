package bitstamp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const USER_TRANSACTIONS_ENDPOINT = "v2/user_transactions/"

type UserTransaction struct {
	Id       int64     `json:"id"`
	DateTime time.Time `json:"datetime"`
	Fee      string    `json:"fee"`
	Type     string    `json:"type"`
	Quantity string    `json:"quantity"`
	Asset    string    `json:"asset"`
}

func (u *UserTransaction) UnmarshalJSON(data []byte) error {
	tmp := make(map[string]interface{})
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	for key, val := range tmp {
		if key == "datetime" {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot unmarshal datetime")
			}

			u.DateTime, err = time.Parse("2006-01-02 15:04:05.000000", str)
			if err != nil {
				return err
			}
			continue
		}

		if key == "fee" {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot unmarshal fee")
			}

			u.Fee = str
			continue
		}

		if key == "id" {
			f, ok := val.(float64)
			if !ok {
				return fmt.Errorf("cannot unmarshal id")
			}

			u.Id = int64(f)
			continue
		}

		if key == "type" {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot unmarshal type")
			}

			u.Type = str
			continue
		}

		str, ok := val.(string)
		if !ok {
			continue
		}

		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}

		if f == 0.0 {
			continue
		}

		u.Asset = key
		u.Quantity = str
	}

	if u.Asset == "" || u.Quantity == "" {
		return fmt.Errorf("invalid transaction body")
	}

	return nil
}

type UserTransactions []*UserTransaction

func (c *Client) UserTransactions(since int64) (UserTransactions, error) {
	var buf *strings.Reader
	data := url.Values{}

	if since != 0 {
		data.Set("since_timestamp", strconv.FormatInt(since, 10))
	}

	if len(data) > 0 {
		buf = strings.NewReader(data.Encode())
	}

	b, err := c.privateRequest(USER_TRANSACTIONS_ENDPOINT, http.MethodPost, buf)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	tx := UserTransactions{}
	err = json.Unmarshal(b, &tx)
	return tx, err
}
