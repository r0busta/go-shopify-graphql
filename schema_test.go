package shopify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/r0busta/graphql"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// Every query and mutation the client sends must be valid against the schema
// the models were generated from. This catches fields and arguments that
// Shopify removed or renamed in a new API version, which compiling against
// the models alone does not.
func TestQueriesValidateAgainstSchema(t *testing.T) {
	schema := loadModelSchema(t)

	queries := captureQueries(t)

	// Methods that send nothing when called with zero-value arguments. Their
	// queries are covered by the methods they delegate to.
	noQueries := map[string]bool{
		"Collection.CreateBulk": true, // loops over an empty slice
		"Metafield.DeleteBulk":  true, // returns early on an empty slice
	}

	names := make([]string, 0, len(queries))
	for n := range queries {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		if len(queries[name]) == 0 && !noQueries[name] {
			t.Errorf("%s sent no query; if that is expected, list it in noQueries", name)
		}
		for _, q := range queries[name] {
			if _, err := gqlparser.LoadQuery(schema, q); err != nil {
				t.Errorf("%s: %s\n%s", name, err.Error(), q)
			}
		}
	}
}

// loadModelSchema reads schema.graphql from the go-shopify-graphql-model
// module in use, so the test always checks against the version the models
// were generated from.
func loadModelSchema(t *testing.T) *ast.Schema {
	t.Helper()

	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/r0busta/go-shopify-graphql-model/v5").Output()
	if err != nil {
		t.Fatalf("locating go-shopify-graphql-model: %v", err)
	}
	path := filepath.Join(strings.TrimSpace(string(out)), "schema.graphql")

	sdl, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// The checked-in schema omits the schema block; Shopify's root types are
	// not the default names, so declare them.
	src := "schema { query: QueryRoot mutation: Mutation }\n" + string(sdl)
	schema, gerr := gqlparser.LoadSchema(&ast.Source{Name: path, Input: src})
	if gerr != nil {
		t.Fatal(gerr)
	}

	return schema
}

// captureQueries calls every service method with zero-value arguments
// against a transport that records each request and fails it, and returns
// the captured queries keyed by service and method name. Every method has an
// entry, even when it sent nothing.
func captureQueries(t *testing.T) map[string][]string {
	t.Helper()

	rec := &recordingTransport{queries: map[string][]string{}}
	client := NewClient(WithGraphQLClient(graphql.NewClient("http://localhost/graphql", &http.Client{Transport: rec})))

	cv := reflect.ValueOf(client).Elem()
	for i := 0; i < cv.NumField(); i++ {
		field := cv.Field(i)
		if !field.CanInterface() || field.Kind() != reflect.Interface || field.IsNil() {
			continue
		}
		svc := field.Elem()
		for j := 0; j < svc.NumMethod(); j++ {
			method := svc.Type().Method(j)
			name := cv.Type().Field(i).Name + "." + method.Name
			rec.begin(name)

			args := make([]reflect.Value, method.Type.NumIn()-1)
			for k := range args {
				argType := method.Type.In(k + 1)
				if argType == reflect.TypeOf((*context.Context)(nil)).Elem() {
					args[k] = reflect.ValueOf(context.Background())
				} else {
					args[k] = reflect.Zero(argType)
				}
			}

			done := make(chan struct{})
			go func() {
				defer close(done)
				defer func() { _ = recover() }()
				svc.Method(j).Call(args)
			}()
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				// A leaked goroutine would keep writing into rec, so stop here.
				t.Fatalf("%s did not return after the request failed", name)
			}
		}
	}

	return rec.queries
}

// recordingTransport records the GraphQL document of every request under the
// current method name and then fails the request, so no method blocks on a
// response. The one exception is the bulk operation poll: it is answered
// with a completed operation so that BulkQuery goes on to submit its query,
// which travels as the `query` variable of bulkOperationRunQuery and is
// recorded from there.
type recordingTransport struct {
	mu      sync.Mutex
	queries map[string][]string
	current string
}

func (r *recordingTransport) begin(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.current = name
	r.queries[name] = nil
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, _ := io.ReadAll(req.Body)
	var payload struct {
		Query     string                     `json:"query"`
		Variables map[string]json.RawMessage `json:"variables"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || payload.Query == "" {
		return nil, errors.New("recorded: request without a query")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.queries[r.current] = append(r.queries[r.current], payload.Query)

	if strings.Contains(payload.Query, "bulkOperationRunQuery") {
		var bulkQuery string
		if raw, ok := payload.Variables["query"]; ok {
			_ = json.Unmarshal(raw, &bulkQuery)
		}
		if bulkQuery != "" {
			r.queries[r.current] = append(r.queries[r.current], bulkQuery)
		}
	}

	if strings.Contains(payload.Query, "currentBulkOperation") {
		const completed = `{"data":{"currentBulkOperation":{"id":"gid://shopify/BulkOperation/1","status":"COMPLETED","objectCount":"0"}}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(completed)),
			Request:    req,
		}, nil
	}

	return nil, errors.New("recorded")
}
