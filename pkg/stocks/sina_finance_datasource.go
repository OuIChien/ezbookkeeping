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
	// Sina Finance quote API
	sinaFinanceQuoteApiUrl = "http://hq.sinajs.cn/list="
)

// SinaFinanceDataSource defines the structure of Sina Finance data source
type SinaFinanceDataSource struct {
}

// BuildRequests builds the http requests
func (s *SinaFinanceDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	var sinaSymbols []string
	for _, symbol := range symbols {
		upperSymbol := strings.ToUpper(symbol)
		if strings.HasSuffix(upperSymbol, ".HK") {
			sinaSymbols = append(sinaSymbols, "hk"+strings.TrimSuffix(upperSymbol, ".HK"))
		} else if strings.HasSuffix(upperSymbol, ".SS") {
			sinaSymbols = append(sinaSymbols, "sh"+strings.TrimSuffix(upperSymbol, ".SS"))
		} else if strings.HasSuffix(upperSymbol, ".SZ") {
			sinaSymbols = append(sinaSymbols, "sz"+strings.TrimSuffix(upperSymbol, ".SZ"))
		} else if !strings.Contains(upperSymbol, ".") {
			sinaSymbols = append(sinaSymbols, "gb_"+strings.ToLower(upperSymbol))
		} else {
			sinaSymbols = append(sinaSymbols, "gb_"+strings.ToLower(upperSymbol))
		}
	}

	u, err := url.Parse(sinaFinanceQuoteApiUrl + strings.Join(sinaSymbols, ","))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://finance.sina.com.cn/")

	return []*http.Request{req}, nil
}

// Parse parses the response content
func (s *SinaFinanceDataSource) Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error) {
	// Sina returns data in GB18030 encoding
	// Format: var hq_str_gb_aapl="Name,Price,Change,Time,...";
	lines := strings.Split(string(content), ";")
	prices := make(models.LatestStockPriceSlice, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "=") {
			continue
		}

		// Extract symbol from variable name
		varPart := strings.Split(line, "=")[0]
		sinaSymbol := ""
		if parts := strings.Split(varPart, "hq_str_"); len(parts) > 1 {
			sinaSymbol = parts[1]
		}

		// Extract content between quotes
		firstQuote := strings.Index(line, "\"")
		lastQuote := strings.LastIndex(line, "\"")
		if firstQuote == -1 || lastQuote == -1 || firstQuote >= lastQuote {
			continue
		}

		dataStr := line[firstQuote+1 : lastQuote]
		parts := strings.Split(dataStr, ",")
		if len(parts) < 2 {
			continue
		}

		// 0: Name
		// 1: Price
		priceStr := parts[1]

		symbol := sinaSymbol
		currency := "USD"

		if strings.HasPrefix(sinaSymbol, "gb_") {
			symbol = strings.ToUpper(strings.TrimPrefix(sinaSymbol, "gb_"))
			currency = "USD"
		} else if strings.HasPrefix(sinaSymbol, "hk") {
			symbol = strings.ToUpper(strings.TrimPrefix(sinaSymbol, "hk")) + ".HK"
			currency = "HKD"
		} else if strings.HasPrefix(sinaSymbol, "sh") {
			symbol = strings.ToUpper(strings.TrimPrefix(sinaSymbol, "sh")) + ".SS"
			currency = "CNY"
		} else if strings.HasPrefix(sinaSymbol, "sz") {
			symbol = strings.ToUpper(strings.TrimPrefix(sinaSymbol, "sz")) + ".SZ"
			currency = "CNY"
		}

		// Validate price
		if _, err := strconv.ParseFloat(priceStr, 64); err != nil {
			continue
		}

		prices = append(prices, &models.LatestStockPrice{
			Symbol:   symbol,
			Price:    priceStr,
			Currency: currency,
		})
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("no valid stock prices found in response")
	}

	return &models.LatestStockPriceResponse{
		DataSource:   "Sina Finance",
		ReferenceUrl: "https://finance.sina.com.cn/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: "USD",
		Prices:       prices,
	}, nil
}
