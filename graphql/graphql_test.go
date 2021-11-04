package graphql

import (
	"context"
	"encoding/json"
	_errors "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

// func TestDo(t *testing.T) {
// }

type Response struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func MakeHTTPCall(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := &Response{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, err
	}
	return r, nil
}

func TestDo(t *testing.T) {
	testTable := []struct {
		name             string
		server           *httptest.Server
		expectedResponse *Response
		expectedErr      error
	}{
		{
			name: "happy-server-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 1, "name": "kyle", "description": "novice gopher"}`))
			})),
			expectedResponse: &Response{
				ID:          1,
				Name:        "kyle",
				Description: "novice gopher",
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.server.Close()
			resp, err := MakeHTTPCall(tc.server.URL)
			if !reflect.DeepEqual(resp, tc.expectedResponse) {
				t.Errorf("expected (%v), got (%v)", tc.expectedResponse, resp)
			}
			if !_errors.Is(err, tc.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
			}
		})
	}
}

func TestQuery(t *testing.T) {
	timeNow := time.Now()
	testTable := []struct {
		name             string
		server           *httptest.Server
		expectedResponse *Response
		expectedErr      error
	}{
		{
			name: "throtled_query",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				os.Create("/home/dragonborn/work/ES-HL/pull/go-shopify-graphql-modified/go-shopify-graphql/graphql/testttttttttttttttttt.go")
				fmt.Println("in mock serverrrrrrrrrrrrrrrrrrrrrr")
				var in struct {
					Query     string                 `json:"query"`
					Variables map[string]interface{} `json:"variables,omitempty"`
				}
				_ = json.NewDecoder(r.Body).Decode(&in)
				if timeNow.Sub(time.Now()) < 2*time.Second {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"id": 1, "name": "kyle", "description": "novice gopher"}`))
				}
				// if in.Query == "throtled_query" {
				// 	w.WriteHeader(http.StatusOK)
				// 	w.Write([]byte(`{"id": 1, "name": "kyle", "description": "novice gopher"}`))
				// }
				// if in.Query == "nice_query" {
				// 	w.WriteHeader(http.StatusOK)
				// 	w.Write([]byte(`{"id": 1, "name": "kyle", "description": "novice gopher"}`))
				// }
			})),
			expectedResponse: &Response{
				ID:          1,
				Name:        "kyle",
				Description: "novice gopher",
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.server.Close()
			c := NewClient(tc.server.URL, tc.server.Client())
			var m map[string]interface{}
			var v interface{}
			t1 := time.Now()
			fmt.Println("hehehe")
			_ = c.do(context.Background(), tc.name, m, v)
			t2 := time.Now()
			if t1.Sub(t2) < 2*time.Second {
				t.Error("too much time")
			}
			// resp, err := MakeHTTPCall(tc.server.URL)
			// if !reflect.DeepEqual(resp, tc.expectedResponse) {
			// 	t.Errorf("expected (%v), got (%v)", tc.expectedResponse, resp)
			// }
			// if !_errors.Is(err, tc.expectedErr) {
			// 	t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
			// }
		})
	}
}

func TestSleepWhenRateLimitQuery(t *testing.T) {
	timeNow := time.Now()
	testTable := []struct {
		name             string
		server           *httptest.Server
		expectedResponse *Response
		expectedErr      error
		Duration         time.Duration
	}{
		{
			name: "throtled_query",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if timeNow.Sub(time.Now())-time.Second < 5 {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
					   "errors":[
						  {
							 "message":"Throttled"
						  }
					   ],
					   "extensions":{
						  "cost":{
							 "requestedQueryCost":452,
							 "actualQueryCost":null,
							 "throttleStatus":{
								"maximumAvailable":1000.0,
								"currentlyAvailable":202,
								"restoreRate":50.0
							 }
						  }
					   }
					}`))
				} else {
					w.WriteHeader(http.StatusOK)
				}
			})),
			expectedResponse: &Response{
				ID:          1,
				Name:        "kyle",
				Description: "novice gopher",
			},
			Duration:    5 * time.Second,
			expectedErr: nil,
		},
	}
	tc := testTable[0]
	t.Run(tc.name, func(t *testing.T) {
		defer tc.server.Close()
		c := NewClient(tc.server.URL, tc.server.Client())
		var m map[string]interface{}
		var v interface{}
		t1 := time.Now()
		_ = c.do(context.Background(), tc.name, m, v)
		t2 := time.Now()
		if t1.Sub(t2)-tc.Duration < 1 {
			t.Log("too much time")
		}
	})

}

// type API struct {
// 	Client  *http.Client
// 	baseURL string
// }

// func (api *API) DoStuff() ([]byte, error) {
// 	resp, err := api.Client.Get(api.baseURL + "/some/path")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	// handling error and doing stuff with body that needs to be unit tested
// 	return body, err
// }

// func TestDoStuffWithTestServer(t *testing.T) {
// 	// Start a local HTTP server
// 	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 		// Test request parameters
// 		os.Create("/home/dragonborn/work/ES-HL/pull/go-shopify-graphql-modified/go-shopify-graphql/graphql/testttttttttttttttttt.go")
// 		// equals(t, req.URL.String(), "/some/path")
// 		// Send response to be tested
// 		rw.Write([]byte(`OK`))
// 	}))
// 	// Close the server when test finishes
// 	defer server.Close()

// 	// Use Client & URL from our local test server
// 	api := API{server.Client(), server.URL}
// 	api.DoStuff()

// 	// ok(t, err)
// 	// equals(t, []byte("OK"), body)

// }
