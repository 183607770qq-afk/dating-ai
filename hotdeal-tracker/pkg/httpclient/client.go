package httpclient

import (
	"bytes"
	"fmt"
	"hotdeal-tracker/internal/config"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/proxy"
	"compress/gzip"
	"compress/flate"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
}

type HTTPClient struct {
	client     *http.Client
	userAgent  string
	timeout    time.Duration
	proxyPool  []string
	useProxy   bool
	proxyIndex int
	random     *rand.Rand
}

func NewHTTPClient(cfg *config.CrawlerConfig) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
				ForceAttemptHTTP2:   true,
			},
		},
		userAgent: cfg.UserAgent,
		timeout:   time.Duration(cfg.Timeout) * time.Second,
		proxyPool: cfg.ProxyPool,
		useProxy:  len(cfg.ProxyPool) > 0,
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *HTTPClient) randomDelay() {
	delay := time.Duration(c.random.Intn(3000)+2000) * time.Millisecond
	time.Sleep(delay)
}

func (c *HTTPClient) Get(url string, headers map[string]string) ([]byte, error) {
	c.randomDelay()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	ua := c.userAgent
	if ua == "" {
		ua = userAgents[c.random.Intn(len(userAgents))]
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if c.useProxy {
		c.setupProxy()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := c.decodeResponseBody(resp)
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "..."
		}
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, bodyStr)
	}

	return c.decodeResponseBody(resp)
}

func (c *HTTPClient) decodeResponseBody(resp *http.Response) ([]byte, error) {
	contentEncoding := resp.Header.Get("Content-Encoding")
	
	var reader io.Reader = resp.Body

	switch strings.ToLower(contentEncoding) {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		flReader := flate.NewReader(resp.Body)
		defer flReader.Close()
		reader = flReader
	case "br":
		reader = brotli.NewReader(resp.Body)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	return c.fixCharset(body, contentType)
}

func (c *HTTPClient) fixCharset(body []byte, contentType string) ([]byte, error) {
	if strings.Contains(contentType, "charset=") {
		return body, nil
	}

	reader, err := charset.NewReaderLabel("utf-8", bytes.NewReader(body))
	if err != nil {
		return body, nil
	}

	return io.ReadAll(reader)
}

func (c *HTTPClient) setupProxy() {
	if len(c.proxyPool) == 0 {
		return
	}

	proxyStr := c.proxyPool[c.proxyIndex]
	c.proxyIndex = (c.proxyIndex + 1) % len(c.proxyPool)

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return
	}

	transport, ok := c.client.Transport.(*http.Transport)
	if !ok {
		return
	}

	if strings.HasPrefix(strings.ToLower(proxyStr), "socks5://") {
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
		if err == nil {
			transport.DialContext = dialer.(proxy.ContextDialer).DialContext
			transport.Proxy = nil
		}
	} else {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
}

func (c *HTTPClient) Post(uri string, data map[string]string, headers map[string]string) ([]byte, error) {
	formData := make(url.Values)
	for key, value := range data {
		formData.Set(key, value)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json, text/html, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if c.useProxy {
		c.setupProxy()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return c.decodeResponseBody(resp)
}

func (c *HTTPClient) getNextProxy() (*url.URL, error) {
	if len(c.proxyPool) == 0 {
		return nil, fmt.Errorf("no proxies available")
	}

	proxyStr := c.proxyPool[c.proxyIndex]
	c.proxyIndex = (c.proxyIndex + 1) % len(c.proxyPool)

	return url.Parse(proxyStr)
}

func (c *HTTPClient) SetProxy(proxyURL string) error {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	transport, ok := c.client.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("failed to get transport")
	}

	transport.Proxy = http.ProxyURL(proxy)
	return nil
}

func (c *HTTPClient) ClearProxy() {
	transport, ok := c.client.Transport.(*http.Transport)
	if ok {
		transport.Proxy = nil
	}
}

func (c *HTTPClient) GetWithChromeDP(url string) ([]byte, error) {
	return nil, fmt.Errorf("ChromeDP is not available in this build")
}