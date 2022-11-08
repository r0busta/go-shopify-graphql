package graphqlclient

import (
	"fmt"
	"net/http"

	"github.com/r0busta/graphql"
)

const (
	shopifyBaseDomain        = "myshopify.com"
	shopifyAccessTokenHeader = "X-Shopify-Access-Token"

	defaultAPIProtocol = "https"
	defaultAPIEndpoint = "graphql.json"
	defaultAPIBasePath = "admin/api"
)

// Option is used to configure options.
type Option func(t *transport)

// WithVersion optionally sets the API version if the passed string is valid.
func WithVersion(apiVersion string) Option {
	return func(t *transport) {
		if apiVersion != "" {
			t.apiBasePath = fmt.Sprintf("%s/%s", defaultAPIBasePath, apiVersion)
		}
	}
}

// WithToken optionally sets access token.
func WithToken(token string) Option {
	return func(t *transport) {
		t.accessToken = token
	}
}

// WithPrivateAppAuth optionally sets private app credentials (API key and access token).
func WithPrivateAppAuth(apiKey string, accessToken string) Option {
	return func(t *transport) {
		t.apiKey = apiKey
		t.accessToken = accessToken
	}
}

type transport struct {
	accessToken string
	apiKey      string
	apiBasePath string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	isAccessTokenSet := t.accessToken != ""
	areBasicAuthCredentialsSet := t.apiKey != "" && isAccessTokenSet

	if areBasicAuthCredentialsSet {
		req.SetBasicAuth(t.apiKey, t.accessToken)
	} else if isAccessTokenSet {
		req.Header.Set(shopifyAccessTokenHeader, t.accessToken)
	}

	return http.DefaultTransport.RoundTrip(req)
}

// NewClient creates a new client (in fact, just a simple wrapper for a graphql.Client).
func NewClient(shopName string, opts ...Option) *graphql.Client {
	transport := &transport{
		apiBasePath: defaultAPIBasePath,
	}

	for _, opt := range opts {
		opt(transport)
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	url := buildAPIEndpoint(shopName, transport.apiBasePath)

	return graphql.NewClient(url, httpClient)
}

func buildAPIEndpoint(shopName string, apiPathPrefix string) string {
	return fmt.Sprintf("%s://%s.%s/%s/%s", defaultAPIProtocol, shopName, shopifyBaseDomain, apiPathPrefix, defaultAPIEndpoint)
}
