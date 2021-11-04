package schemaregistry

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type (
	httpDoer interface {
		Do(req *http.Request) (*http.Response, error)
	}

	IClient interface {
		Subjects() (subjects []string, err error)
		Versions(subject string) (versions []int, err error)
		DeleteSubject(subject string) (versions []string, err error)
		IsRegistered(subject, schema string) (bool, Schema, error)
		RegisterNewSchema(subject string, avroSchema string) (int, error)
		GetSchemaById(id int) (string, error)
		GetSchemaByVersion(subject string, version string) (*Schema, error)
		GetLatestSchema(subject string) (*Schema, error)
		IsSchemaCompatible(subject string, avroSchema string, version int) (bool, error)
		IsLatestSchemaCompatible(subject string, avroSchema string) (bool, error)
	}

	Client struct {
		baseUrl string
		client  httpDoer
	}

	Option func(*Client)
)

func getTransportLayer(httpClient *http.Client, timeout time.Duration) http.RoundTripper {
	if t := httpClient.Transport; t != nil {
		return t
	}

	httpTransport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	if timeout > 0 {
		httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			return net.DialTimeout(network, addr, timeout)
		}
	}

	return httpTransport
}

func usingClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient == nil {
			return
		}

		transport := getTransportLayer(httpClient, 0)
		httpClient.Transport = transport

		c.client = httpClient
	}
}

func NewClient(baseUrl string) (IClient, error) {
	if baseUrl == "" {
		return nil, errRequired("baseUrl")
	}

	if _, err := url.ParseRequestURI(baseUrl); err != nil {
		return nil, err
	}

	c := &Client{baseUrl: baseUrl}
	if c.client == nil {
		httpClient := &http.Client{}
		usingClient(httpClient)(c)
	}

	return c, nil
}

type ResourceError struct {
	ErrorCode int    `json:"error_code"`
	Method    string `json:"method,omitempty"`
	Uri       string `json:"uri,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (err ResourceError) Error() string {
	return fmt.Sprintf("client: (%s: %s) failed with error: %d%s",
		err.Uri, err.Method, err.ErrorCode, err.Message)
}

func newResourceError(errorCode int, method, uri, message string) ResourceError {
	unescapedUri, err := url.QueryUnescape(uri)
	if err != nil {
		panic(err)
	}

	return ResourceError{
		ErrorCode: errorCode,
		Method:    method,
		Uri:       unescapedUri,
		Message:   message,
	}
}

const (
	subjectNotFoundCode = 40401
	schemaNotFoundCode  = 40403
	versionNotFound     = 40402
)

func IsSubjectNotFound(err error) bool {
	if err == nil {
		return false
	}

	if resErr, ok := err.(ResourceError); ok {
		return resErr.ErrorCode == subjectNotFoundCode
	}

	return false
}

func IsSchemaNotFound(err error) bool {
	if err == nil {
		return false
	}

	if resErr, ok := err.(ResourceError); ok {
		return resErr.ErrorCode == schemaNotFoundCode
	}

	return false
}

func IsVersionNotFound(err error) bool {
	if err == nil {
		return false
	}

	if resErr, ok := err.(ResourceError); ok {
		return resErr.ErrorCode == versionNotFound
	}

	return false
}

func isOk(resp *http.Response) bool {
	return !(resp.StatusCode < 200 || resp.StatusCode >= 300)
}

func acquireBuffer(b []byte) *bytes.Buffer {
	if len(b) > 0 {
		return bytes.NewBuffer(b)
	}

	return new(bytes.Buffer)
}

const (
	contentTypeHeaderKey  = "Content-Type"
	contentTypeJSON       = "application/json"
	contentTypeSchemaJSON = "application/vnd.schemaregistry.v1+json"

	acceptHeaderKey          = "Accept"
	acceptEncodingHeaderKey  = "Accept-Encoding"
	contentEncodingHeaderKey = "Content-Encoding"
	gzipEncodingHeaderValue  = "gzip"
)

type gzipReadCloser struct {
	gzipReader io.ReadCloser
	respReader io.ReadCloser
}

func (rc *gzipReadCloser) Close() error {
	if rc.gzipReader != nil {
		defer rc.gzipReader.Close()
	}

	return rc.respReader.Close()
}

func (rc *gzipReadCloser) Read(p []byte) (n int, err error) {
	if rc.gzipReader != nil {
		return rc.gzipReader.Read(p)
	}

	return rc.respReader.Read(p)
}

func (c *Client) acquireResponseBodyStream(resp *http.Response) (io.ReadCloser, error) {
	// check for gzip
	var (
		reader = resp.Body
		err    error
	)

	if encoding := resp.Header.Get(contentEncodingHeaderKey); encoding == gzipEncodingHeaderValue {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("client: failed to read gzip compressed content, trace: %v", err)
		}

		return &gzipReadCloser{
			respReader: resp.Body,
			gzipReader: reader,
		}, nil
	}

	return reader, err
}

func (c *Client) readResponseBody(resp *http.Response) ([]byte, error) {
	reader, err := c.acquireResponseBodyStream(resp)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(reader)
	if err = reader.Close(); err != nil {
		return nil, err
	}

	return body, err
}

func (c *Client) readJSON(resp *http.Response, val interface{}) error {
	b, err := c.readResponseBody(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, val)
}

func (c *Client) do(method, path, contentType string, send []byte) (*http.Response, error) {
	if path[0] == '/' {
		path = path[1:]
	}

	uri := c.baseUrl + "/" + path

	req, err := http.NewRequest(method, uri, acquireBuffer(send))
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set(contentTypeHeaderKey, contentType)
	}

	req.Header.Add(acceptEncodingHeaderKey, gzipEncodingHeaderValue)
	req.Header.Add(acceptHeaderKey, contentTypeJSON+", "+contentTypeSchemaJSON)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if !isOk(resp) {
		defer resp.Body.Close()
		var resError ResourceError
		c.readJSON(resp, &resError)

		return nil, resError
	}

	return resp, nil
}

const (
	subjectsPath = "subjects"
	subjectPath  = subjectsPath + "/%s"
	versionsPath = subjectPath + "/versions"
	versionPath  = versionsPath + "/%s"
	schemaPath   = "schemas/%d"
)

var errRequired = func(field string) error {
	return fmt.Errorf("client: %s is required", field)
}

// Subjects returns list of subjects
func (c *Client) Subjects() (subjects []string, err error) {

	// GET /subjects
	resp, resError := c.do(http.MethodGet, subjectsPath, "", nil)
	if resError != nil {
		err = resError
		return
	}

	err = c.readJSON(resp, &subjects)
	return
}

// Versions returns all versions of a subject
func (c *Client) Versions(subject string) (versions []int, err error) {
	if subject == "" {
		err = errRequired("subject")
		return
	}

	// GET /subjects/{string: subject}/versions
	path := fmt.Sprintf(versionsPath, subject)
	resp, resError := c.do(http.MethodGet, path, "", nil)
	if resError != nil {
		err = resError
		return
	}

	err = c.readJSON(resp, &versions)
	return
}

// DeleteSubject deletes subject and returns deleted versions belong with it
func (c *Client) DeleteSubject(subject string) (versions []string, err error) {
	if subject == "" {
		err = errRequired("subject")
		return
	}

	// DELETE /subjects/{string: subject}
	path := fmt.Sprintf(subjectPath, subject)
	resp, resError := c.do(http.MethodDelete, path, "", nil)
	if resError != nil {
		err = resError
		return
	}

	err = c.readJSON(resp, &versions)
	return
}

type (
	schemaOnlyJSON struct {
		Schema string `json:"schema"`
	}

	idOnlyJSON struct {
		ID int `json:"id"`
	}

	isCompatibleJSON struct {
		IsCompatible bool `json:"is_compatible"`
	}

	Schema struct {
		Schema  string `json:"schema"`
		Subject string `json:"subject"`
		Version int    `json:"version"`
		ID      int    `json:"id,omitempty"`
	}
)

// SchemaLatestVersion only valid string for version, it's the "latest" version string
const SchemaLatestVersion = "latest"

// IsRegistered returns true if the given schema is registered already
func (c *Client) IsRegistered(subject, schema string) (bool, Schema, error) {
	var sc Schema

	if subject == "" {
		return false, sc, errRequired("subject")
	}
	if schema == "" {
		return false, sc, errRequired("schema")
	}

	schemaOnly := schemaOnlyJSON{schema}
	send, err := json.Marshal(schemaOnly)
	if err != nil {
		return false, sc, err
	}

	// POST /subjects/{string: subject}
	path := fmt.Sprintf(subjectPath, subject)
	resp, resErr := c.do(http.MethodPost, path, "", send)
	if resErr != nil {
		// is schema found?
		if IsSchemaNotFound(resErr) {
			return false, sc, nil
		}

		return false, sc, resErr
	}

	if err = c.readJSON(resp, &sc); err != nil {
		return true, sc, err // found but error when unmarshal
	}

	return true, sc, nil
}

// RegisterNewSchema registers a new schema and returns id of it
func (c *Client) RegisterNewSchema(subject string, avroSchema string) (int, error) {
	if subject == "" {
		return 0, errRequired("subject")
	}
	if avroSchema == "" {
		return 0, errRequired("schema")
	}

	schema := schemaOnlyJSON{avroSchema}
	send, err := json.Marshal(schema)
	if err != nil {
		return 0, err
	}

	// POST /subjects/{string: subject}/versions
	path := fmt.Sprintf(versionsPath, subject)
	resp, err := c.do(http.MethodPost, path, contentTypeSchemaJSON, send)
	if err != nil {
		return 0, err
	}

	var idOnly idOnlyJSON
	err = c.readJSON(resp, &idOnly)
	return idOnly.ID, err
}

// GetSchemaById gets schema by id
func (c *Client) GetSchemaById(id int) (string, error) {

	// GET /schemas/{int: id}
	path := fmt.Sprintf(schemaPath, id)
	resp, err := c.do(http.MethodGet, path, "", nil)
	if err != nil {
		return "", err
	}

	var sc schemaOnlyJSON
	if err = c.readJSON(resp, &sc); err != nil {
		return "", err
	}

	return sc.Schema, nil
}

// GetSchemaByVersion gets schema by version number
func (c *Client) GetSchemaByVersion(subject string, version string) (*Schema, error) {
	if subject == "" {
		return nil, errRequired("subject")
	}

	if err := checkSchemaVersionNumber(version); err != nil {
		return nil, err
	}

	// GET /subjects/{string: subject}/versions/{string: version}
	path := fmt.Sprintf(versionPath, subject, version)
	resp, err := c.do(http.MethodGet, path, "", nil)
	if err != nil {
		return nil, err
	}

	var schema Schema
	if err = c.readJSON(resp, &schema); err != nil {
		return nil, err
	}

	return &schema, nil
}

// GetLatestSchema gets the latest schema of subject
func (c *Client) GetLatestSchema(subject string) (*Schema, error) {
	return c.GetSchemaByVersion(subject, SchemaLatestVersion)
}

func checkSchemaVersionNumber(versionNumber interface{}) error {
	if versionNumber == nil {
		return errRequired("versionNumber must be string \"latest\" or int")
	}

	if verStr, ok := versionNumber.(string); ok {
		if verStr != SchemaLatestVersion {
			return fmt.Errorf("client: %v string is not a valid value for versionNumber", versionNumber)
		}
	}

	if verInt, ok := versionNumber.(int); ok {
		if verInt <= 0 || verInt > 2^31-1 {
			return fmt.Errorf("client: %v string is not a valid value for versionNumber", versionNumber)
		}
	}

	return nil
}

func (c *Client) isSchemaCompatibleAtVersion(subject string, avroSchema string, version interface{}) (bool, error) {
	if subject == "" {
		return false, errRequired("subject")
	}
	if avroSchema == "" {
		return false, errRequired("avroSchema")
	}
	if err := checkSchemaVersionNumber(version); err != nil {
		return false, err
	}

	schema := schemaOnlyJSON{avroSchema}

	send, err := json.Marshal(schema)
	if err != nil {
		return false, err
	}

	// POST /compatibility/subjects/{string: subject}/versions/{string: version}
	path := fmt.Sprintf("compatibility/"+versionPath, subject, version)
	resp, err := c.do(http.MethodPost, path, contentTypeSchemaJSON, send)
	if err != nil {
		return false, err
	}

	var isCompatibleJSON isCompatibleJSON
	err = c.readJSON(resp, &isCompatibleJSON)

	return isCompatibleJSON.IsCompatible, err
}

// IsSchemaCompatible is schema is compatible with version
func (c *Client) IsSchemaCompatible(subject string, avroSchema string, version int) (bool, error) {
	return c.isSchemaCompatibleAtVersion(subject, avroSchema, version)
}

// IsLatestSchemaCompatible is schema is compatible with last version
func (c *Client) IsLatestSchemaCompatible(subject string, avroSchema string) (bool, error) {
	return c.isSchemaCompatibleAtVersion(subject, avroSchema, SchemaLatestVersion)
}
