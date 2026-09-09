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
	if len(queries) == 0 {
		t.Fatal("no queries captured")
	}

	names := make([]string, 0, len(queries))
	for n := range queries {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
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
// against a transport that records the request body and fails the request,
// and returns the captured queries keyed by service and method name.
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
			rec.current = name

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
				t.Errorf("%s did not return after the request failed", name)
			}
		}
	}

	return rec.queries
}

type recordingTransport struct {
	queries map[string][]string
	current string
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, _ := io.ReadAll(req.Body)
	var payload struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Query != "" {
		r.queries[r.current] = append(r.queries[r.current], payload.Query)
	}
	return nil, errors.New("recorded")
}
