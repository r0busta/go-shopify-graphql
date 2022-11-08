package graphqlclient

import (
	"testing"
)

func TestAPIURLWithVersion(t *testing.T) {
	transport := &transport{}
	WithVersion("2019-10")(transport)

	expected := "admin/api/2019-10"
	actual := transport.apiBasePath
	if actual != expected {
		t.Errorf("WithVersion apiBasePath = %s, expected %s", actual, expected)
	}
}

func TestAPIURLWithEmptyVersion(t *testing.T) {
	transport := &transport{}
	WithVersion("")(transport)

	expected := ""
	actual := transport.apiBasePath
	if actual != expected {
		t.Errorf("WithVersion apiBasePath = %s, expected %s", actual, expected)
	}
}

func TestBuildAPIEndpoint(t *testing.T) {
	expected := "https://store.myshopify.com/admin/api/graphql.json"
	actual := buildAPIEndpoint("store", "admin/api")
	if actual != expected {
		t.Errorf("buildAPIEndpoint = %s, expected %s", actual, expected)
	}
}
