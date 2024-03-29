package downloader

import (
	"errors"
	"golang.org/x/net/http/httpproxy"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

type Session struct {
	client    *http.Client
	transport *http.Transport
}

// NewSession new a Session object, and set a default Client and Transport.
func NewSession(options ...*SessionOptions) *Session {
	var sessionOptions *SessionOptions
	if len(options) > 0 {
		sessionOptions = options[0]
	} else {
		sessionOptions = DefaultSessionOptions()
	}

	// set transport parameters.
	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   sessionOptions.DialTimeout,
			KeepAlive: sessionOptions.DialKeepAlive,
		}).DialContext,
		MaxIdleConns:          sessionOptions.MaxIdleConns,
		MaxIdleConnsPerHost:   sessionOptions.MaxIdleConnsPerHost,
		MaxConnsPerHost:       sessionOptions.MaxConnsPerHost,
		IdleConnTimeout:       sessionOptions.IdleConnTimeout,
		TLSHandshakeTimeout:   sessionOptions.TLSHandshakeTimeout,
		ExpectContinueTimeout: sessionOptions.ExpectContinueTimeout,
		Proxy:                 proxyFunc,
	}
	if sessionOptions.DisableDialKeepAlives {
		trans.DisableKeepAlives = true
	}

	client := &http.Client{
		Transport:     trans,
		CheckRedirect: redirectFunc,
	}

	// set CookieJar
	if sessionOptions.DisableCookieJar == false {
		cookieJarOptions := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}
		jar, err := cookiejar.New(&cookieJarOptions)
		if err != nil {
			return nil
		}
		client.Jar = jar
	}

	return &Session{
		client:    client,
		transport: trans,
	}
}

type SessionOptions struct {
	// DialTimeout is the maximum amount of time a dial will wait for
	// a connect to complete.
	//
	// When using TCP and dialing a host name with multiple IP
	// addresses, the timeout may be divided between them.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	DialTimeout time.Duration

	// DialKeepAlive specifies the interval between keep-alive
	// probes for an active network connection.
	//
	// Network protocols or operating systems that do
	// not support keep-alives ignore this field.
	// If negative, keep-alive probes are disabled.
	DialKeepAlive time.Duration

	// MaxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	//
	// Zero means no limit.
	MaxConnsPerHost int

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration

	// TLSHandshakeTimeout specifies the maximum amount of time waiting to
	// wait for a TLS handshake. Zero means no timeout.
	TLSHandshakeTimeout time.Duration

	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	ExpectContinueTimeout time.Duration

	// DisableCookieJar specifies whether disable session cookiejar.
	DisableCookieJar bool

	// DisableDialKeepAlives, if true, disables HTTP keep-alives and
	// will only use the connection to the server for a single
	// HTTP request.
	//
	// This is unrelated to the similarly named TCP keep-alives.
	DisableDialKeepAlives bool
}

// DefaultSessionOptions return a default SessionOptions object.
func DefaultSessionOptions() *SessionOptions {
	return &SessionOptions{
		DialTimeout:           30 * time.Second,
		DialKeepAlive:         30 * time.Second,
		MaxConnsPerHost:       0,					// 对每个host的最大连接数量，0表示不限制
		MaxIdleConns:          100,					// 所有host的连接池最大连接数量，默认100
		MaxIdleConnsPerHost:   2,					// 每个host的连接池最大空闲连接数,默认2
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCookieJar:      false,
		DisableDialKeepAlives: false,
	}
}

var (
	// proxyConfigOnce guards proxyConfig
	envProxyOnce      sync.Once
	envProxyFuncValue func(*url.URL) (*url.URL, error)
)

// proxyFunc get proxy from request context.
// If there is no proxy set, use default proxy from environment.
func proxyFunc(req *http.Request) (*url.URL, error) {
	httpURLStr := req.Context().Value("http")   // get http proxy url form context
	httpsURLStr := req.Context().Value("https") // get https proxy url form context

	// If there is no proxy set, use default proxy from environment.
	// This mitigates expensive lookups on some platforms (e.g. Windows).
	envProxyOnce.Do(func() {
		envProxyFuncValue = httpproxy.FromEnvironment().ProxyFunc()
	})

	if req.URL.Scheme == "http" { // set proxy for http site
		if httpURLStr != nil {
			httpURL, err := url.Parse(httpURLStr.(string))
			if err != nil {
				return nil, err
			}
			return httpURL, nil
		}
	} else if req.URL.Scheme == "https" { // set proxy for https site
		if httpsURLStr != nil {
			httpsURL, err := url.Parse(httpsURLStr.(string))
			if err != nil {
				return nil, err
			}
			return httpsURL, nil
		}
	}

	return envProxyFuncValue(req.URL)
}

// redirectFunc get redirectNum from request context and check redirect number.
func redirectFunc(req *http.Request, via []*http.Request) error {
	redirectNum := req.Context().Value("redirectNum").(int)
	if len(via) > redirectNum {
		return errors.New("RedirectError")
	}
	return nil
}
