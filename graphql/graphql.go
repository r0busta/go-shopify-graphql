package graphqlclient

import (
	"fmt"
	"net/http"

	"github.com/r0busta/graphql"
)

const (
	shopifyBaseDomain        = "myshopify.com"
	shopifyAccessTokenHeader = "X-Shopify-Access-Token"
)

var (
	apiProtocol   = "https"
	apiPathPrefix = "admin/api"
	apiEndpoint   = "graphql.json"
)

// Option is used to configure options.
type Option func(t *transport)

// WithVersion optionally sets the API version if the passed string is valid.
func WithVersion(apiVersion string) Option {
	return func(t *transport) {
		if apiVersion != "" {
			apiPathPrefix = fmt.Sprintf("admin/api/%s", apiVersion)
		} else {
			apiPathPrefix = "admin/api"
		}
	}
}

// WithToken optionally sets oauth token.
func WithToken(token string) Option {
	return func(t *transport) {
		t.accessToken = token
	}
}

// WithPrivateAppAuth optionally sets private app credentials.
func WithPrivateAppAuth(apiKey string, password string) Option {
	return func(t *transport) {
		t.apiKey = apiKey
		t.password = password
	}
}

type transport struct {
	accessToken string
	apiKey      string
	password    string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.accessToken != "" {
		req.Header.Set(shopifyAccessTokenHeader, t.accessToken)
	} else if t.apiKey != "" && t.password != "" {
		req.SetBasicAuth(t.apiKey, t.password)
	}

	return http.DefaultTransport.RoundTrip(req)
}

// NewClient creates a new client (in fact, just a simple wrapper for a graphql.Client).
func NewClient(shopName string, opts ...Option) *graphql.Client {
	transport := &transport{}

	for _, opt := range opts {
		opt(transport)
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	url := buildAPIEndpoint(shopName)

	return graphql.NewClient(url, httpClient)
}

func buildAPIEndpoint(shopName string) string {
	return fmt.Sprintf("%s://%s.%s/%s/%s", apiProtocol, shopName, shopifyBaseDomain, apiPathPrefix, apiEndpoint)
}
