package stocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const (
	// Yahoo Finance public quote API
	yahooFinanceQuoteApiUrl = "https://query1.finance.yahoo.com/v7/finance/quote"
	yahooFinanceCookieUrl   = "https://fc.yahoo.com"
	yahooFinanceCrumbUrl    = "https://query1.finance.yahoo.com/v1/test/getcrumb"
)

// YahooFinanceDataSource defines the structure of Yahoo Finance data source
type YahooFinanceDataSource struct {
	cookie    string
	crumb     string
	expiresAt time.Time
	mu        sync.RWMutex
}

// YahooFinanceQuoteResponse represents the response from Yahoo Finance API
type YahooFinanceQuoteResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol             string  `json:"symbol"`
			RegularMarketPrice float64 `json:"regularMarketPrice"`
			Currency           string  `json:"currency"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteResponse"`
}

// BuildRequests builds the http requests
func (s *YahooFinanceDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	if err := s.ensureCrumb(); err != nil {
		// Log error but try anyway? No, it will fail 401.
		// However, we can return the error.
		return nil, fmt.Errorf("failed to get yahoo crumb: %v", err)
	}

	s.mu.RLock()
	cookie := s.cookie
	crumb := s.crumb
	s.mu.RUnlock()

	u, err := url.Parse(yahooFinanceQuoteApiUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("symbols", strings.Join(symbols, ","))
	q.Set("fields", "symbol,regularMarketPrice,currency")
	q.Set("crumb", crumb)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)

	return []*http.Request{req}, nil
}

func (s *YahooFinanceDataSource) ensureCrumb() error {
	s.mu.RLock()
	if s.crumb != "" && time.Now().Before(s.expiresAt) {
		s.mu.RUnlock()
		return nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check
	if s.crumb != "" && time.Now().Before(s.expiresAt) {
		return nil
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	// 1. Get Cookie
	req1, err := http.NewRequest("GET", yahooFinanceCookieUrl, nil)
	if err != nil {
		return err
	}
	req1.Header.Set("User-Agent", userAgent)
	resp1, err := client.Do(req1)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	// We expect 404 or 200, but we just want the cookies
	// Read Set-Cookie header
	cookies := resp1.Cookies()
	var bCookie string
	for _, c := range cookies {
		if c.Name == "B" {
			bCookie = c.Value
			break
		}
	}

	// Sometimes fc.yahoo.com redirects or behaves oddly.
	// If standard cookies are not found in jar (not using jar here), we parse header manually?
	// resp1.Cookies() parses Set-Cookie headers.

	if bCookie == "" {
		// Fallback: Check if response has any cookies or if we can proceed without it (unlikely)
		// Some sources say https://finance.yahoo.com is better for cookie?
		// But let's try to proceed.
	}

	// Construct cookie string for next request
	cookieStr := ""
	for _, c := range cookies {
		cookieStr += c.Name + "=" + c.Value + ";"
	}

	// 2. Get Crumb
	req2, err := http.NewRequest("GET", yahooFinanceCrumbUrl, nil)
	if err != nil {
		return err
	}
	req2.Header.Set("User-Agent", userAgent)
	req2.Header.Set("Cookie", cookieStr)

	resp2, err := client.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		return fmt.Errorf("get crumb failed with status: %d", resp2.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp2.Body)
	if err != nil {
		return err
	}

	crumb := strings.TrimSpace(string(bodyBytes))
	if crumb == "" {
		return fmt.Errorf("empty crumb received")
	}

	s.cookie = cookieStr
	s.crumb = crumb
	s.expiresAt = time.Now().Add(24 * time.Hour) // Cache for 24h?

	return nil
}

// Parse parses the response content
func (s *YahooFinanceDataSource) Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error) {
	var response YahooFinanceQuoteResponse
	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, err
	}

	if response.QuoteResponse.Error != nil {
		return nil, fmt.Errorf("yahoo finance api error: %v", response.QuoteResponse.Error)
	}

	prices := make(models.LatestStockPriceSlice, 0, len(response.QuoteResponse.Result))

	for _, result := range response.QuoteResponse.Result {
		prices = append(prices, &models.LatestStockPrice{
			Symbol:   strings.ToUpper(result.Symbol),
			Price:    strconv.FormatFloat(result.RegularMarketPrice, 'f', -1, 64),
			Currency: strings.ToUpper(result.Currency),
		})
	}

	return &models.LatestStockPriceResponse{
		DataSource:   "Yahoo Finance",
		ReferenceUrl: "https://finance.yahoo.com/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: "USD", // This is a bit complex as Yahoo returns prices in various currencies based on exchange, but for now we standardise on response field if needed. The proposal says "Total Account Value = Held Quantity × Real-time Market Price", and "converts the valuation to user's Default Currency".
		Prices:       prices,
	}, nil
}
