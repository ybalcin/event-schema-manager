package schemaregistry

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type (
	doFn func(req *http.Request) (*http.Response, error)
)

func (d doFn) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

func mustEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, but got %#v", expected, actual)
	}
}

func mockHTTPHandler(status int, reqBody, respBody interface{}) doFn {
	d := doFn(func(req *http.Request) (*http.Response, error) {
		if reqBody != nil {
			_, err := json.Marshal(reqBody)
			if err != nil {
				return nil, err
			}
		}
		var resp http.Response
		resp.Header = http.Header{"Content-Type": []string{"application/json"}}
		resp.StatusCode = status
		if respBody != nil {
			bs, err := json.Marshal(respBody)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		}
		return &resp, nil
	})
	return d
}

func httpSuccess(reqBody, respBody interface{}) doFn {
	return mockHTTPHandler(200, reqBody, respBody)
}

func httpError(status, errCode int, errMessage string) doFn {
	return mockHTTPHandler(status, nil, ResourceError{ErrorCode: errCode, Message: errMessage})
}

func TestNewClient(t *testing.T) {
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

func TestIsSubjectNotFound(t *testing.T) {
	errSubNotFound := ResourceError{
		ErrorCode: subjectNotFoundCode,
	}
	errOtherCode := ResourceError{
		ErrorCode: 123,
	}
	notResourceErr := errors.New("")

	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errSubNotFound, true},
		{errOtherCode, false},
		{notResourceErr, false},
	}

	for _, c := range tests {
		if c.expected != IsSubjectNotFound(c.err) {
			t.Fail()
		}
	}
}

func TestIsSchemaNotFound(t *testing.T) {
	errSchemaNotFound := ResourceError{
		ErrorCode: schemaNotFoundCode,
	}
	errOtherCode := ResourceError{
		ErrorCode: 123,
	}
	notResourceErr := errors.New("")

	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errSchemaNotFound, true},
		{errOtherCode, false},
		{notResourceErr, false},
	}

	for _, c := range tests {
		if c.expected != IsSchemaNotFound(c.err) {
			t.Fail()
		}
	}
}

func TestIsVersionNotFound(t *testing.T) {
	errVersionNotFound := ResourceError{
		ErrorCode: versionNotFound,
	}
	errOtherCode := ResourceError{
		ErrorCode: 123,
	}
	notResourceErr := errors.New("")

	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errVersionNotFound, true},
		{errOtherCode, false},
		{notResourceErr, false},
	}

	for _, c := range tests {
		if c.expected != IsVersionNotFound(c.err) {
			t.Fail()
		}
	}
}

func TestResourceError_Error(t *testing.T) {
	resErr := ResourceError{}

	tests := []ResourceError{
		resErr,
	}

	for _, r := range tests {
		if r.Error() == "" {
			t.Fail()
		}
	}
}

func TestClient_Subjects(t *testing.T) {
	expected := []string{"sub1", "sub2"}
	cli := Client{
		baseUrl: "",
		client:  httpSuccess(nil, expected),
	}
	subs, err := cli.Subjects()
	if err != nil {
		t.Error()
	}
	mustEqual(t, subs, expected)
}
