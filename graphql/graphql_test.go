package graphqlclient

import (
	"fmt"
	"testing"
)

func TestWithVersion(t *testing.T) {
	_ = NewClient("myshop", WithVersion("2019-10"))
	expected := fmt.Sprintf("admin/api/%s", "2019-10")
	if defaultAPIPathPrefix != expected {
		t.Errorf("WithVersion apiPathPrefix = %s, expected %s", defaultAPIPathPrefix, expected)
	}
}

func TestWithVersionEmptyVersion(t *testing.T) {
	_ = NewClient("myshop", WithVersion(""))
	expected := "admin/api"
	if defaultAPIPathPrefix != expected {
		t.Errorf("WithVersion apiPathPrefix = %s, expected %s", defaultAPIPathPrefix, expected)
	}
}

func TestWithoutVersion(t *testing.T) {
	_ = NewClient("myshop")
	expected := "admin/api"
	if defaultAPIPathPrefix != expected {
		t.Errorf("WithVersion apiPathPrefix = %s, expected %s", defaultAPIPathPrefix, expected)
	}
}
