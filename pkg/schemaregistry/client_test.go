package schemaregistry_test

import (
	. "event-schema-manager/pkg/schemaregistry"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	tests := []struct {
		input        string
		failExpected bool
	}{
		{"", true},
		{"lorem ipsum", true},
		{"http://localhost:8081", false},
	}

	for _, s := range tests {
		cli, err := NewClient(s.input)
		if s.failExpected {
			if err == nil {
				t.Fail()
			}
		} else {
			if err != nil || cli == nil {
				t.Fail()
			}
		}
	}
}
