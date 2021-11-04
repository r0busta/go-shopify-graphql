package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// Client is a GraphQL client.
type Client struct {
	url        string // GraphQL server URL.
	httpClient *http.Client
}

// NewClient creates a GraphQL client targeting the specified GraphQL server URL.
// If httpClient is nil, then http.DefaultClient is used.
func NewClient(url string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		url:        url,
		httpClient: httpClient,
	}
}

// QueryString executes a single GraphQL query request,
// using the given raw query `q` and populating the response into the `v`.
// `q` should be a correct GraphQL request string that corresponds to the GraphQL schema.
func (c *Client) QueryString(ctx context.Context, q string, variables map[string]interface{}, v interface{}) error {
	return c.do(ctx, q, variables, v)
}

// Query executes a single GraphQL query request,
// with a query derived from q, populating the response into it.
// q should be a pointer to struct that corresponds to the GraphQL schema.
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	query := constructQuery(q, variables)
	return c.do(ctx, query, variables, q)
}

// Mutate executes a single GraphQL mutation request,
// with a mutation derived from m, populating the response into it.
// m should be a pointer to struct that corresponds to the GraphQL schema.
func (c *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	query := constructMutation(m, variables)
	fmt.Println(query)
	// return nil
	return c.do(ctx, query, variables, m)
}

// do executes a single GraphQL operation.
func (c *Client) do(ctx context.Context, query string, variables map[string]interface{}, v interface{}) error {
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}
	resp, err := ctxhttp.Post(ctx, c.httpClient, c.url, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}
	var out struct {
		Data       *json.RawMessage
		Errors     errors
		Extensions interface{} // Unused.
	}
	err = json.NewDecoder(resp.Body).Decode(&out)

	if len(out.Errors) > 0 {
		if out.Errors[0].Message == "Throttled" {
			b, err := json.Marshal(out.Extensions)
			if err != nil {
				return err
			}
			var extensions dict
			err = json.Unmarshal(b, &extensions)
			if err != nil {
				return err
			}
			requestedQueryCost := extensions.d("cost").s("requestedQueryCost")
			throttleStatus := extensions.d("cost").d("throttleStatus")
			currentlyAvailable := throttleStatus.s("currentlyAvailable")
			restoreRate := throttleStatus.s("restoreRate")
			if currentlyAvailable < requestedQueryCost {
				timeSleep := int((requestedQueryCost - currentlyAvailable) / restoreRate)
				time.Sleep(time.Duration(timeSleep) * time.Second)
			}
		}
	}

	if err != nil {
		// TODO: Consider including response body in returned error, if deemed helpful.
		return err
	}
	// xx := make(map[string]interface{})
	if out.Data != nil {
		err := json.Unmarshal(*out.Data, v)
		if err != nil {
			// TODO: Consider including response body in returned error, if deemed helpful.
			return err
		}
	}
	if len(out.Errors) > 0 {
		return out.Errors
	}
	return nil
}

// Accessing Nested Map of Type map[string]interface{} in Golang
// struct is the best option, but if you insist,
// you can add a type declaration for a map, then you can add methods to help with the type assertions:
type dict map[string]interface{}

// convert first index to map that has interface value.
func (d dict) d(k string) dict {
	if d[k] == nil {
		return nil
	}
	return d[k].(map[string]interface{})
}

// convert the item value to string
func (d dict) s(k string) float64 {
	if d[k] == nil {
		return 0
	}
	return d[k].(float64)
}

// errors represents the "errors" array in a response from a GraphQL server.
// If returned via error interface, the slice is expected to contain at least 1 element.
//
// Specification: https://facebook.github.io/graphql/#sec-Errors.
type errors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

// Error implements error interface.
func (e errors) Error() string {
	return e[0].Message
}

type operationType uint8

const (
	queryOperation operationType = iota
	mutationOperation
	//subscriptionOperation // Unused.
)
