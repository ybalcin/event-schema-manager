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
	mustErrorCodeEqual(t, actual, expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, but got %#v", expected, actual)
	}
}

func mustNotNil(t *testing.T, actual interface{}) {
	if actual == nil {
		t.Errorf("%#v must not nil", actual)
	}
}

func mustErrorCodeEqual(t *testing.T, actual, expected interface{}) {
	if actErr, ok := actual.(ResourceError); ok {
		if expErr, ok := expected.(ResourceError); ok {
			if actErr.ErrorCode != expErr.ErrorCode {
				t.Errorf("actual error code: %#v is not equal to expected: %#v", actual, expected)
			}
		}
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
		if !isOk(&resp) {
			return &resp, respBody.(ResourceError)
		}
		return &resp, nil
	})
	return d
}

func mockHttpSuccess(reqBody, respBody interface{}) doFn {
	return mockHTTPHandler(200, reqBody, respBody)
}

func mockHttpError(status, errCode int, reqBody interface{}, errMessage string) doFn {
	return mockHTTPHandler(status, reqBody, ResourceError{ErrorCode: errCode, Message: errMessage})
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
	cli := Client{client: mockHttpSuccess(nil, expected)}
	subs, err := cli.Subjects()
	if err != nil {
		t.Error(err)
	}
	mustEqual(t, subs, expected)
}

func TestClient_Versions(t *testing.T) {
	expected := []int{1, 2, 3}

	type testItem struct {
		subject     string
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{"", errRequired("subject"), nil}, // if subject is empty should return err
		{testSubject, ResourceError{ErrorCode: subjectNotFoundCode}, mockHttpError(http.StatusNotFound, subjectNotFoundCode, nil, "")}, // if subject not found should return err
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		versions, err := cli.Versions(c.subject)
		mustEqual(t, versions, ([]int)(nil))
		mustEqual(t, err, c.expected)
	}

	testsSuccess := []testItem{
		{testSubject, expected, mockHttpSuccess(nil, expected)},
	}

	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		versions, err := cli.Versions(c.subject)
		mustEqual(t, versions, c.expected)
		mustEqual(t, err, nil)
	}
}

func TestClient_DeleteSubject(t *testing.T) {
	expected := []string{"1", "2", "3"}

	type testItem struct {
		subject     string
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{"", errRequired("subject"), nil}, // if subject is empty should return err
		{testSubject, ResourceError{ErrorCode: subjectNotFoundCode}, mockHttpError(http.StatusNotFound, subjectNotFoundCode, nil, "")}, // if subject not found should return err
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		versions, err := cli.DeleteSubject(c.subject)
		mustEqual(t, versions, ([]string)(nil))
		mustEqual(t, err, c.expected)
	}

	testsSuccess := []testItem{
		{testSubject, expected, mockHttpSuccess(nil, expected)},
	}

	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		versions, err := cli.DeleteSubject(c.subject)
		mustEqual(t, err, nil)
		mustEqual(t, versions, c.expected)
	}
}

const (
	validSchema   string = "{\"namespace\": \"example.avro\",\"type\": \"record\", \"name\": \"user\",\"fields\":[{\"name\": \"name\",\"type\": \"string\"},{ \"name\": \"favorite_number\",\"type\": \"int\"}]}"
	invalidSchema string = "loremipsum"
	testSubject   string = "testsubject"
)

func TestClient_IsRegistered(t *testing.T) {

	reqBody := schemaOnlyJSON{validSchema}
	expectedRespBody := Schema{
		Schema:  validSchema,
		Subject: testSubject,
		Version: 1,
		ID:      1,
	}

	type testItem struct {
		subject     string
		schema      string
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{testSubject, invalidSchema, ResourceError{ErrorCode: invalidAvroSchema}, mockHttpError(http.StatusUnprocessableEntity, invalidAvroSchema, invalidSchema, "")}, // if schema is invalid should return err
		{testSubject, "", errRequired("schema"), nil},  // if schema is empty should return err
		{"", validSchema, errRequired("subject"), nil}, // if subject is empty should return err
	}
	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		_, _, err := cli.IsRegistered(c.subject, c.schema)
		mustNotNil(t, err)

		if c.expected != struct{}{} {
			mustEqual(t, err, c.expected)
		} else {
			mustNotNil(t, err)
		}
	}

	testsSuccess := []testItem{
		{testSubject, validSchema, expectedRespBody, mockHttpSuccess(reqBody, expectedRespBody)},               // should return isRegistered true
		{testSubject, validSchema, false, mockHttpError(http.StatusNotFound, schemaNotFoundCode, reqBody, "")}, // should return isRegistered false
	}
	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		isRegistered, sc, _ := cli.IsRegistered(c.subject, c.schema)
		mustNotNil(t, sc)

		if val, ok := c.expected.(bool); ok {
			mustEqual(t, isRegistered, val)
		} else {
			mustEqual(t, sc, c.expected)
		}
	}
}

func TestClient_RegisterNewSchema(t *testing.T) {

	reqBody := schemaOnlyJSON{validSchema}
	expectedRespBody := Schema{
		Schema:  validSchema,
		Subject: testSubject,
		Version: 1,
		ID:      1,
	}

	type testItem struct {
		subject     string
		schema      string
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{"", validSchema, errRequired("subject"), nil}, // if subject is not found should return err
		{testSubject, "", errRequired("schema"), nil},  // if schema is not found should return err
		{testSubject, invalidSchema, ResourceError{ErrorCode: invalidAvroSchema}, mockHttpError(http.StatusUnprocessableEntity, invalidAvroSchema, invalidSchema, "")}, // if schema is invalid should return err
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		_, err := cli.RegisterNewSchema(c.subject, c.schema)
		mustNotNil(t, err)

		if c.expected != struct{}{} {
			mustEqual(t, err, c.expected)
		} else {
			mustNotNil(t, err)
		}
	}

	testsSuccess := []testItem{
		{testSubject, validSchema, expectedRespBody.ID, mockHttpSuccess(reqBody, expectedRespBody)},
	}

	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		id, err := cli.RegisterNewSchema(c.subject, c.schema)
		mustEqual(t, err, nil)
		mustEqual(t, id, c.expected)
	}
}

func TestClient_GetSchemaById(t *testing.T) {

	type testItem struct {
		id          int
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{0, ResourceError{ErrorCode: schemaNotFoundCode}, mockHttpError(http.StatusNotFound, schemaNotFoundCode, nil, "")}, // if schema is not found should return err
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		sc, err := cli.GetSchemaById(c.id)
		mustEqual(t, err, c.expected)
		mustEqual(t, sc, "")
	}

	schemaOnly := schemaOnlyJSON{validSchema}
	testsSuccess := []testItem{
		{1, validSchema, mockHttpSuccess(nil, schemaOnly)},
	}

	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		sc, err := cli.GetSchemaById(c.id)
		mustEqual(t, err, nil)
		mustEqual(t, sc, validSchema)
	}
}

func TestClient_GetSchemaByVersion(t *testing.T) {

	expectedRespBody := Schema{
		Schema:  validSchema,
		Subject: testSubject,
		Version: 1,
		ID:      1,
	}

	type testItem struct {
		subject     string
		version     string
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{"", "1", errRequired("subject"), nil},         // if subject is empty should return err
		{testSubject, "", errRequired("version"), nil}, // if version is empty should return err
		{testSubject, "abc", struct{}{}, nil},          // if version is invalid should return err
		{testSubject, "1", ResourceError{ErrorCode: subjectNotFoundCode}, mockHttpError(404, subjectNotFoundCode, nil, "")}, // if subject is not found should return err
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		sc, err := cli.GetSchemaByVersion(c.subject, c.version)

		mustEqual(t, sc, (*Schema)(nil))
		if c.expected != struct{}{} {
			mustEqual(t, err, c.expected)
		} else {
			mustNotNil(t, err)
		}
	}

	testsSuccess := []testItem{
		{testSubject, "1", expectedRespBody, mockHttpSuccess(nil, expectedRespBody)},
	}

	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		sc, err := cli.GetSchemaByVersion(c.subject, c.version)

		mustEqual(t, err, nil)
		mustEqual(t, *sc, c.expected)
	}
}

func TestClient_GetLatestSchema(t *testing.T) {

	type testItem struct {
		subject     string
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{"", errRequired("subject"), nil}, // if subject is empty should return err
		{testSubject, ResourceError{ErrorCode: schemaNotFoundCode}, mockHttpError(404, schemaNotFoundCode, nil, "")}, // if subject is not found should return error
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		sc, err := cli.GetLatestSchema(c.subject)
		mustEqual(t, sc, (*Schema)(nil))
		mustEqual(t, err, c.expected)
	}
}

func TestClient_IsSchemaCompatible(t *testing.T) {

	type testItem struct {
		subject     string
		schema      string
		version     int
		expected    interface{}
		mockHandler doFn
	}

	testsError := []testItem{
		{"", validSchema, 1, errRequired("subject"), nil},    // if subject is empty should return err
		{testSubject, "", 1, errRequired("avroSchema"), nil}, // if schema is empty should return err
		{testSubject, validSchema, 0, struct{}{}, nil},       // if version is invalid should return err
		{testSubject, validSchema, -1, struct{}{}, nil},      // if version is invalid should return err
		{testSubject, invalidSchema, 1, ResourceError{ErrorCode: invalidAvroSchema}, mockHttpError(http.StatusUnprocessableEntity, invalidAvroSchema, invalidSchema, "")}, // if schema is invalid should return err
	}

	for _, c := range testsError {
		cli := Client{client: c.mockHandler}
		is, err := cli.IsSchemaCompatible(c.subject, c.schema, c.version)
		mustEqual(t, is, false)
		if c.expected != struct{}{} {
			mustEqual(t, err, c.expected)
		} else {
			mustNotNil(t, err)
		}
	}

	reqBody := schemaOnlyJSON{validSchema}

	testsSuccess := []testItem{
		{testSubject, validSchema, 1, true, mockHttpSuccess(reqBody, isCompatibleJSON{IsCompatible: true})},   // if schema is compatible should return true
		{testSubject, validSchema, 1, false, mockHttpSuccess(reqBody, isCompatibleJSON{IsCompatible: false})}, // if schema is incompatible should return false
	}

	for _, c := range testsSuccess {
		cli := Client{client: c.mockHandler}
		is, err := cli.IsSchemaCompatible(c.subject, c.schema, c.version)
		mustEqual(t, err, nil)
		mustEqual(t, is, c.expected)
	}
}
