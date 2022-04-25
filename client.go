package bitstamp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const API_URL = "https://www.bitstamp.net/api/"
const X_AUTH = "X-Auth"
const X_AUTH_VERSION = "X-Auth-Version"
const X_AUTH_NONCE = "X-Auth-Nonce"
const X_AUTH_TIMESTAMP = "X-Auth-Timestamp"
const X_AUTH_SIGNATURE = "X-Auth-Signature"
const CONTENT_TYPE = "Content-Type"

type Client struct {
	apiKey    string
	apiSecret string
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

type ApiError struct {
	Code int
	Body string
}

func newApiError(code int, message []byte) *ApiError {
	return &ApiError{
		Code: code,
		Body: string(message),
	}
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API Error: Code(%d) %v", e.Code, e.Body)
}

func (c *Client) addAuthorization(req *http.Request) error {
	nonce := uuid.New()
	timestamp := time.Now().UnixMilli()
	var buf *bytes.Buffer

	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body.Close()

		buf = bytes.NewBuffer(b)
		req.Body = io.NopCloser(buf)
	}

	req.Header.Add(X_AUTH, fmt.Sprintf("BITSTAMP %v", c.apiKey))
	req.Header.Add(X_AUTH_VERSION, "v2")
	req.Header.Add(X_AUTH_NONCE, nonce.String())
	req.Header.Add(X_AUTH_TIMESTAMP, strconv.FormatInt(timestamp, 10))
	if buf != nil {
		req.Header.Add(CONTENT_TYPE, "application/x-www-form-urlencoded")
	}

	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte("BITSTAMP "))
	h.Write([]byte(c.apiKey))
	h.Write([]byte(req.Method))
	h.Write([]byte(req.URL.Host))
	h.Write([]byte(req.URL.Path))
	h.Write([]byte(req.URL.RawQuery))

	if buf != nil {
		h.Write([]byte("application/x-www-form-urlencoded"))
	}

	h.Write([]byte(nonce.String()))
	h.Write([]byte(strconv.FormatInt(timestamp, 10)))
	h.Write([]byte("v2"))

	if buf != nil {
		h.Write(buf.Bytes())
	}

	sig := hex.EncodeToString(h.Sum(nil))
	req.Header.Add(X_AUTH_SIGNATURE, sig)
	return nil
}

func (c *Client) privateRequest(endpoint string, method string, body io.Reader) ([]byte, error) {
	if c.apiKey == "" || c.apiSecret == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	url := fmt.Sprintf("%v%v", API_URL, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	err = c.addAuthorization(req)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		apiErr := newApiError(resp.StatusCode, b)
		return nil, apiErr
	}

	return b, nil
}
