package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type defaultTransport struct {
	defaultUserAgent string
	baseTransport    http.RoundTripper
}

func (t *defaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Values("User-Agent")) < 1 {
		req.Header.Set("User-Agent", t.defaultUserAgent)
	} else if req.Header.Get("User-Agent") == "" {
		req.Header.Del("User-Agent")
	}

	return t.baseTransport.RoundTrip(req)
}

// NewHttpClient creates and returns a new http client with specified settings
func NewHttpClient(requestTimeout uint32, proxy string, skipTLSVerify bool, defaultUserAgent string) *http.Client {
	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	SetProxyUrl(baseTransport, proxy)

	if skipTLSVerify {
		baseTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Transport: &defaultTransport{
			defaultUserAgent: defaultUserAgent,
			baseTransport:    baseTransport,
		},
		Timeout: time.Duration(requestTimeout) * time.Millisecond,
	}
}

// SetProxyUrl sets proxy url to http transport according to specified proxy setting
func SetProxyUrl(transport *http.Transport, proxy string) {
	if proxy == "none" {
		transport.Proxy = nil
	} else if proxy != "system" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			// If proxy URL is invalid, fall back to no proxy
			transport.Proxy = nil
		} else {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	} else {
		// Use system proxy, but check if environment variables are actually set
		// If no valid proxy environment variables are found, fall back to no proxy
		if hasValidProxyEnv() {
			transport.Proxy = http.ProxyFromEnvironment
		} else {
			transport.Proxy = nil
		}
	}
}

// hasValidProxyEnv checks if there are valid proxy environment variables set
func hasValidProxyEnv() bool {
	proxyVars := []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy"}
	for _, v := range proxyVars {
		if val := os.Getenv(v); val != "" && strings.TrimSpace(val) != "" {
			// Try to parse the proxy URL to validate it
			if _, err := url.Parse(val); err == nil {
				return true
			}
		}
	}
	return false
}
