package stocks

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const (
	// Tencent Finance quote API (Tencent/gtimg)
	// Example: http://qt.gtimg.cn/q=usAAPL,hk00700,sh600519
	tencentFinanceQuoteApiUrl = "http://qt.gtimg.cn/q="
)

// TencentFinanceDataSource defines the structure of Tencent Finance data source
type TencentFinanceDataSource struct {
}

// BuildRequests builds the http requests
func (s *TencentFinanceDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	var tencentSymbols []string
	for _, symbol := range symbols {
		upperSymbol := strings.ToUpper(symbol)
		if strings.HasSuffix(upperSymbol, ".HK") {
			tencentSymbols = append(tencentSymbols, "hk"+strings.TrimSuffix(upperSymbol, ".HK"))
		} else if strings.HasSuffix(upperSymbol, ".SS") {
			tencentSymbols = append(tencentSymbols, "sh"+strings.TrimSuffix(upperSymbol, ".SS"))
		} else if strings.HasSuffix(upperSymbol, ".SZ") {
			tencentSymbols = append(tencentSymbols, "sz"+strings.TrimSuffix(upperSymbol, ".SZ"))
		} else if !strings.Contains(upperSymbol, ".") {
			// Assume US stock if no suffix
			tencentSymbols = append(tencentSymbols, "us"+upperSymbol)
		} else {
			// Fallback: try as is or prefix with us if not sh/sz/hk
			tencentSymbols = append(tencentSymbols, "us"+upperSymbol)
		}
	}

	u, err := url.Parse(tencentFinanceQuoteApiUrl + strings.Join(tencentSymbols, ","))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	return []*http.Request{req}, nil
}

// Parse parses the response content
func (s *TencentFinanceDataSource) Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error) {
	// Tencent returns data in GBK encoding, but we only care about numbers and symbols (ASCII)
	// Format: v_usAAPL="200~Apple~AAPL.OQ~269.48~...~USD~...";
	lines := strings.Split(string(content), ";")
	prices := make(models.LatestStockPriceSlice, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract content between quotes
		firstQuote := strings.Index(line, "\"")
		lastQuote := strings.LastIndex(line, "\"")
		if firstQuote == -1 || lastQuote == -1 || firstQuote >= lastQuote {
			continue
		}

		dataStr := line[firstQuote+1 : lastQuote]
		parts := strings.Split(dataStr, "~")
		if len(parts) < 4 {
			continue
		}

		// Field mapping for US stocks (Tencent format can vary slightly between markets but generally):
		// 1: Name
		// 2: Symbol
		// 3: Current Price
		// 34: Currency (for US stocks)

		rawSymbol := parts[2]
		priceStr := parts[3]

		// Clean symbol (e.g., AAPL.OQ -> AAPL)
		symbol := rawSymbol
		if dotIndex := strings.Index(symbol, "."); dotIndex != -1 {
			symbol = symbol[:dotIndex]
		}

		// If it was a HK stock, rawSymbol might be 00700
		if len(rawSymbol) == 5 && strings.HasPrefix(rawSymbol, "0") {
			symbol = rawSymbol + ".HK"
		}

		currency := "USD"
		if len(parts) > 34 && parts[34] != "" {
			currency = parts[34]
		} else if strings.Contains(line, "v_sh") || strings.Contains(line, "v_sz") {
			currency = "CNY"
		} else if strings.Contains(line, "v_hk") {
			currency = "HKD"
		}

		// Validate price
		if _, err := strconv.ParseFloat(priceStr, 64); err != nil {
			continue
		}

		prices = append(prices, &models.LatestStockPrice{
			Symbol:   strings.ToUpper(symbol),
			Price:    priceStr,
			Currency: strings.ToUpper(currency),
		})
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("no valid stock prices found in response")
	}

	return &models.LatestStockPriceResponse{
		DataSource:   "Tencent Finance",
		ReferenceUrl: "https://gu.qq.com/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: "USD",
		Prices:       prices,
	}, nil
}
