// Package v2 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package v2

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// NumberRecordsParams defines parameters for NumberRecords.
type NumberRecordsParams struct {
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// MessageRecords request
	MessageRecords(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// MessageRecord request
	MessageRecord(ctx context.Context, messageID uint64, reqEditors ...RequestEditorFn) (*http.Response, error)

	// NumberRecords request
	NumberRecords(ctx context.Context, params *NumberRecordsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) MessageRecords(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMessageRecordsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) MessageRecord(ctx context.Context, messageID uint64, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMessageRecordRequest(c.Server, messageID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) NumberRecords(ctx context.Context, params *NumberRecordsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNumberRecordsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewMessageRecordsRequest generates requests for MessageRecords
func NewMessageRecordsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/message/records")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewMessageRecordRequest generates requests for MessageRecord
func NewMessageRecordRequest(server string, messageID uint64) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "message_id", runtime.ParamLocationPath, messageID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/message/records/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewNumberRecordsRequest generates requests for NumberRecords
func NewNumberRecordsRequest(server string, params *NumberRecordsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/numbers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// MessageRecordsWithResponse request
	MessageRecordsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MessageRecordsResponse, error)

	// MessageRecordWithResponse request
	MessageRecordWithResponse(ctx context.Context, messageID uint64, reqEditors ...RequestEditorFn) (*MessageRecordResponse, error)

	// NumberRecordsWithResponse request
	NumberRecordsWithResponse(ctx context.Context, params *NumberRecordsParams, reqEditors ...RequestEditorFn) (*NumberRecordsResponse, error)
}

type MessageRecordsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]struct {
		// Body The text content of the message.
		Body string `json:"body"`

		// Direction The direction of the message.
		Direction N200Direction `json:"direction"`

		// MessageId The internal ID of the message.
		MessageID uint64 `json:"message_id"`

		// Timestamp The time the message was sent or received.
		Timestamp time.Time `json:"timestamp"`
	}
}

// Status returns HTTPResponse.Status
func (r MessageRecordsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MessageRecordsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type MessageRecordResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Body The text content of the message.
		Body string `json:"body"`

		// Direction The direction of the message.
		Direction N200Direction `json:"direction"`

		// MessageId The internal ID of the message.
		MessageID uint64 `json:"message_id"`

		// Timestamp The time the message was sent or received.
		Timestamp time.Time `json:"timestamp"`
	}
}

// Status returns HTTPResponse.Status
func (r MessageRecordResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MessageRecordResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type NumberRecordsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]struct {
		// Number A number.
		Number string `json:"number"`
	}
}

// Status returns HTTPResponse.Status
func (r NumberRecordsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NumberRecordsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// MessageRecordsWithResponse request returning *MessageRecordsResponse
func (c *ClientWithResponses) MessageRecordsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MessageRecordsResponse, error) {
	rsp, err := c.MessageRecords(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMessageRecordsResponse(rsp)
}

// MessageRecordWithResponse request returning *MessageRecordResponse
func (c *ClientWithResponses) MessageRecordWithResponse(ctx context.Context, messageID uint64, reqEditors ...RequestEditorFn) (*MessageRecordResponse, error) {
	rsp, err := c.MessageRecord(ctx, messageID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMessageRecordResponse(rsp)
}

// NumberRecordsWithResponse request returning *NumberRecordsResponse
func (c *ClientWithResponses) NumberRecordsWithResponse(ctx context.Context, params *NumberRecordsParams, reqEditors ...RequestEditorFn) (*NumberRecordsResponse, error) {
	rsp, err := c.NumberRecords(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseNumberRecordsResponse(rsp)
}

// ParseMessageRecordsResponse parses an HTTP response from a MessageRecordsWithResponse call
func ParseMessageRecordsResponse(rsp *http.Response) (*MessageRecordsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MessageRecordsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []struct {
			// Body The text content of the message.
			Body string `json:"body"`

			// Direction The direction of the message.
			Direction N200Direction `json:"direction"`

			// MessageId The internal ID of the message.
			MessageID uint64 `json:"message_id"`

			// Timestamp The time the message was sent or received.
			Timestamp time.Time `json:"timestamp"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseMessageRecordResponse parses an HTTP response from a MessageRecordWithResponse call
func ParseMessageRecordResponse(rsp *http.Response) (*MessageRecordResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MessageRecordResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Body The text content of the message.
			Body string `json:"body"`

			// Direction The direction of the message.
			Direction N200Direction `json:"direction"`

			// MessageId The internal ID of the message.
			MessageID uint64 `json:"message_id"`

			// Timestamp The time the message was sent or received.
			Timestamp time.Time `json:"timestamp"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseNumberRecordsResponse parses an HTTP response from a NumberRecordsWithResponse call
func ParseNumberRecordsResponse(rsp *http.Response) (*NumberRecordsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &NumberRecordsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []struct {
			// Number A number.
			Number string `json:"number"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xWS2/jRgz+KwO2R1nyCz3o1mKx6B6aQ3ZPTY1iItH2LDyPcKhsDEP/vZiHLfmRZtPu",
	"TR5xyI/f95HyARqrnTVo2EN9AELvrPGYfzSW2vDUWMNoODxK53aqkaysqb56a8IZvkjtdhgeH227hxp+",
	"x93OQgGtImxCLNRgO95YZTZQgEbv5Qb/Vi3Uy3kBrDR6ltpBDfPpfDGZzSez5ZfZsp4t6sX0T+gL8M0W",
	"tQw1fiZcQw0/VQP2Kr31Vcbc930BLfqGlMvl73NrYk1Wi6cOaS/WlkTGItLNMpRKj/49nT8cwHT6EQlq",
	"mM0XS+hX74Ps34M5lcqQfQn9iZ8L4RxZh8QqCZrEOVwU+bJFwfjCIvcq7FrwFo/ElFAA7x1CDZ4p6Nef",
	"6Xor3en1jVxoOg31AyjTWJ3scHLG6kapsVdu1VKGkYzciU8fblRbW9KSoQZl+Jfl0Eq4tUGCAl4mGzsx",
	"UofTP9LFTx+O5zm6S7f7M6fe5FFpHEMQ36QXPpIa5UL1jO0ZrlYyTsK1a5qjE586RdgGwkZEjHEUSdax",
	"JgON9vErNnzhacWo/fcO0imVJJL75DRsOlK8/xxCk7N+Q0lIv3a8jT6Lvz4eW7ROPnWhv5g75EoBQ8db",
	"Zpf8r8zaXlN7j45s2zXoY+scJu788BnJp9h5OS2nAbd1aKRTUMOinJYLKMBJ3ka4VeayGrGyQb4u/BG5",
	"2V6sCF/+ZSCmp7gMPrWDc+5zvuJ8k86n09f4PsWN9sBAcVwsY3IfgFC2dUbkYRX3TKe1pP2reANpcuPD",
	"7YwTVqHKJQ3VYfBY/wYn8nJxvkVK5J+kRkbysa/3zrIKYUFDKCAP7NlMDMPC1OF4/f7oLTCO1qfoIMV/",
	"lf0HqH6lyKuyp8/HW67PUaW4R9/tWEgKmbkjg61QRkhB0rRWC0st0i3972KGYSYu9I96xo/aIOhOacUw",
	"1k4ro3T4ZMyuZOv/D+PfN2hHqm4yfmRoRPTdKf7fs6d8SM9HMjra5UXo66rK/yvKxupKOlU9z6O7cpVL",
	"vT5zt14LtqK14pvirbgbgGVe0wn0xZt3jyYbXT66p1/1/wQAAP//FspQsjEKAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
