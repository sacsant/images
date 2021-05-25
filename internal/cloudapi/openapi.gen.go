// Package cloudapi provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package cloudapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// AWSUploadRequestOptions defines model for AWSUploadRequestOptions.
type AWSUploadRequestOptions struct {
	Ec2    AWSUploadRequestOptionsEc2 `json:"ec2"`
	Region string                     `json:"region"`
	S3     AWSUploadRequestOptionsS3  `json:"s3"`
}

// AWSUploadRequestOptionsEc2 defines model for AWSUploadRequestOptionsEc2.
type AWSUploadRequestOptionsEc2 struct {
	AccessKeyId       string    `json:"access_key_id"`
	SecretAccessKey   string    `json:"secret_access_key"`
	ShareWithAccounts *[]string `json:"share_with_accounts,omitempty"`
	SnapshotName      *string   `json:"snapshot_name,omitempty"`
}

// AWSUploadRequestOptionsS3 defines model for AWSUploadRequestOptionsS3.
type AWSUploadRequestOptionsS3 struct {
	AccessKeyId     string `json:"access_key_id"`
	Bucket          string `json:"bucket"`
	SecretAccessKey string `json:"secret_access_key"`
}

// AWSUploadStatus defines model for AWSUploadStatus.
type AWSUploadStatus struct {
	Ami    string `json:"ami"`
	Region string `json:"region"`
}

// AzureUploadRequestOptions defines model for AzureUploadRequestOptions.
type AzureUploadRequestOptions struct {

	// Name of the uploaded image. It must be unique in the given resource group.
	// If name is omitted from the request, a random one based on a UUID is
	// generated.
	ImageName *string `json:"image_name,omitempty"`

	// Location where the image should be uploaded and registered. This link explain
	// how to list all locations:
	// https://docs.microsoft.com/en-us/cli/azure/account?view=azure-cli-latest#az_account_list_locations'
	Location string `json:"location"`

	// Name of the resource group where the image should be uploaded.
	ResourceGroup string `json:"resource_group"`

	// ID of subscription where the image should be uploaded.
	SubscriptionId string `json:"subscription_id"`

	// ID of the tenant where the image should be uploaded. This link explains how
	// to find it in the Azure Portal:
	// https://docs.microsoft.com/en-us/azure/active-directory/fundamentals/active-directory-how-to-find-tenant
	TenantId string `json:"tenant_id"`
}

// AzureUploadStatus defines model for AzureUploadStatus.
type AzureUploadStatus struct {
	ImageName string `json:"image_name"`
}

// ComposeRequest defines model for ComposeRequest.
type ComposeRequest struct {
	Customizations *Customizations `json:"customizations,omitempty"`
	Distribution   string          `json:"distribution"`
	ImageRequests  []ImageRequest  `json:"image_requests"`
}

// ComposeResult defines model for ComposeResult.
type ComposeResult struct {
	Id string `json:"id"`
}

// ComposeStatus defines model for ComposeStatus.
type ComposeStatus struct {
	ImageStatus ImageStatus `json:"image_status"`
}

// Customizations defines model for Customizations.
type Customizations struct {
	Packages     *[]string     `json:"packages,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`
}

// GCPUploadRequestOptions defines model for GCPUploadRequestOptions.
type GCPUploadRequestOptions struct {

	// Name of an existing STANDARD Storage class Bucket.
	Bucket string `json:"bucket"`

	// The name to use for the imported and shared Compute Engine image.
	// The image name must be unique within the GCP project, which is used
	// for the OS image upload and import. If not specified a random
	// 'composer-api-<uuid>' string is used as the image name.
	ImageName *string `json:"image_name,omitempty"`

	// The GCP region where the OS image will be imported to and shared from.
	// The value must be a valid GCP location. See https://cloud.google.com/storage/docs/locations.
	// If not specified, the multi-region location closest to the source
	// (source Storage Bucket location) is chosen automatically.
	Region *string `json:"region,omitempty"`

	// List of valid Google accounts to share the imported Compute Engine image with.
	// Each string must contain a specifier of the account type. Valid formats are:
	//   - 'user:{emailid}': An email address that represents a specific
	//     Google account. For example, 'alice@example.com'.
	//   - 'serviceAccount:{emailid}': An email address that represents a
	//     service account. For example, 'my-other-app@appspot.gserviceaccount.com'.
	//   - 'group:{emailid}': An email address that represents a Google group.
	//     For example, 'admins@example.com'.
	//   - 'domain:{domain}': The G Suite domain (primary) that represents all
	//     the users of that domain. For example, 'google.com' or 'example.com'.
	// If not specified, the imported Compute Engine image is not shared with any
	// account.
	ShareWithAccounts *[]string `json:"share_with_accounts,omitempty"`
}

// GCPUploadStatus defines model for GCPUploadStatus.
type GCPUploadStatus struct {
	ImageName string `json:"image_name"`
	ProjectId string `json:"project_id"`
}

// ImageRequest defines model for ImageRequest.
type ImageRequest struct {
	Architecture  string        `json:"architecture"`
	ImageType     string        `json:"image_type"`
	Ostree        *OSTree       `json:"ostree,omitempty"`
	Repositories  []Repository  `json:"repositories"`
	UploadRequest UploadRequest `json:"upload_request"`
}

// ImageStatus defines model for ImageStatus.
type ImageStatus struct {
	Status       ImageStatusValue `json:"status"`
	UploadStatus *UploadStatus    `json:"upload_status,omitempty"`
}

// ImageStatusValue defines model for ImageStatusValue.
type ImageStatusValue string

// List of ImageStatusValue
const (
	ImageStatusValue_building    ImageStatusValue = "building"
	ImageStatusValue_failure     ImageStatusValue = "failure"
	ImageStatusValue_pending     ImageStatusValue = "pending"
	ImageStatusValue_registering ImageStatusValue = "registering"
	ImageStatusValue_success     ImageStatusValue = "success"
	ImageStatusValue_uploading   ImageStatusValue = "uploading"
)

// OSTree defines model for OSTree.
type OSTree struct {
	Ref *string `json:"ref,omitempty"`
	Url *string `json:"url,omitempty"`
}

// Repository defines model for Repository.
type Repository struct {
	Baseurl    *string `json:"baseurl,omitempty"`
	Metalink   *string `json:"metalink,omitempty"`
	Mirrorlist *string `json:"mirrorlist,omitempty"`
	Rhsm       bool    `json:"rhsm"`
}

// Subscription defines model for Subscription.
type Subscription struct {
	ActivationKey string `json:"activation-key"`
	BaseUrl       string `json:"base-url"`
	Insights      bool   `json:"insights"`
	Organization  int    `json:"organization"`
	ServerUrl     string `json:"server-url"`
}

// UploadRequest defines model for UploadRequest.
type UploadRequest struct {
	Options interface{} `json:"options"`
	Type    UploadTypes `json:"type"`
}

// UploadStatus defines model for UploadStatus.
type UploadStatus struct {
	Options interface{} `json:"options"`
	Status  string      `json:"status"`
	Type    UploadTypes `json:"type"`
}

// UploadTypes defines model for UploadTypes.
type UploadTypes string

// List of UploadTypes
const (
	UploadTypes_aws   UploadTypes = "aws"
	UploadTypes_azure UploadTypes = "azure"
	UploadTypes_gcp   UploadTypes = "gcp"
)

// Version defines model for Version.
type Version struct {
	Version string `json:"version"`
}

// ComposeJSONBody defines parameters for Compose.
type ComposeJSONBody ComposeRequest

// ComposeRequestBody defines body for Compose for application/json ContentType.
type ComposeJSONRequestBody ComposeJSONBody

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
	// https://api.deepmap.com for example.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
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
		client.Client = http.DefaultClient
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
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Compose request  with any body
	ComposeWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	Compose(ctx context.Context, body ComposeJSONRequestBody) (*http.Response, error)

	// ComposeStatus request
	ComposeStatus(ctx context.Context, id string) (*http.Response, error)

	// GetOpenapiJson request
	GetOpenapiJson(ctx context.Context) (*http.Response, error)

	// GetVersion request
	GetVersion(ctx context.Context) (*http.Response, error)
}

func (c *Client) ComposeWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewComposeRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) Compose(ctx context.Context, body ComposeJSONRequestBody) (*http.Response, error) {
	req, err := NewComposeRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) ComposeStatus(ctx context.Context, id string) (*http.Response, error) {
	req, err := NewComposeStatusRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetOpenapiJson(ctx context.Context) (*http.Response, error) {
	req, err := NewGetOpenapiJsonRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetVersion(ctx context.Context) (*http.Response, error) {
	req, err := NewGetVersionRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewComposeRequest calls the generic Compose builder with application/json body
func NewComposeRequest(server string, body ComposeJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewComposeRequestWithBody(server, "application/json", bodyReader)
}

// NewComposeRequestWithBody generates requests for Compose with any type of body
func NewComposeRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/compose")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewComposeStatusRequest generates requests for ComposeStatus
func NewComposeStatusRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "id", id)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/compose/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetOpenapiJsonRequest generates requests for GetOpenapiJson
func NewGetOpenapiJsonRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/openapi.json")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetVersionRequest generates requests for GetVersion
func NewGetVersionRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/version")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
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
	// Compose request  with any body
	ComposeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*ComposeResponse, error)

	ComposeWithResponse(ctx context.Context, body ComposeJSONRequestBody) (*ComposeResponse, error)

	// ComposeStatus request
	ComposeStatusWithResponse(ctx context.Context, id string) (*ComposeStatusResponse, error)

	// GetOpenapiJson request
	GetOpenapiJsonWithResponse(ctx context.Context) (*GetOpenapiJsonResponse, error)

	// GetVersion request
	GetVersionWithResponse(ctx context.Context) (*GetVersionResponse, error)
}

type ComposeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *ComposeResult
}

// Status returns HTTPResponse.Status
func (r ComposeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ComposeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ComposeStatusResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ComposeStatus
}

// Status returns HTTPResponse.Status
func (r ComposeStatusResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ComposeStatusResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetOpenapiJsonResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetOpenapiJsonResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetOpenapiJsonResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetVersionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Version
}

// Status returns HTTPResponse.Status
func (r GetVersionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetVersionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ComposeWithBodyWithResponse request with arbitrary body returning *ComposeResponse
func (c *ClientWithResponses) ComposeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*ComposeResponse, error) {
	rsp, err := c.ComposeWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseComposeResponse(rsp)
}

func (c *ClientWithResponses) ComposeWithResponse(ctx context.Context, body ComposeJSONRequestBody) (*ComposeResponse, error) {
	rsp, err := c.Compose(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseComposeResponse(rsp)
}

// ComposeStatusWithResponse request returning *ComposeStatusResponse
func (c *ClientWithResponses) ComposeStatusWithResponse(ctx context.Context, id string) (*ComposeStatusResponse, error) {
	rsp, err := c.ComposeStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParseComposeStatusResponse(rsp)
}

// GetOpenapiJsonWithResponse request returning *GetOpenapiJsonResponse
func (c *ClientWithResponses) GetOpenapiJsonWithResponse(ctx context.Context) (*GetOpenapiJsonResponse, error) {
	rsp, err := c.GetOpenapiJson(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetOpenapiJsonResponse(rsp)
}

// GetVersionWithResponse request returning *GetVersionResponse
func (c *ClientWithResponses) GetVersionWithResponse(ctx context.Context) (*GetVersionResponse, error) {
	rsp, err := c.GetVersion(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetVersionResponse(rsp)
}

// ParseComposeResponse parses an HTTP response from a ComposeWithResponse call
func ParseComposeResponse(rsp *http.Response) (*ComposeResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ComposeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest ComposeResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseComposeStatusResponse parses an HTTP response from a ComposeStatusWithResponse call
func ParseComposeStatusResponse(rsp *http.Response) (*ComposeStatusResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ComposeStatusResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ComposeStatus
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetOpenapiJsonResponse parses an HTTP response from a GetOpenapiJsonWithResponse call
func ParseGetOpenapiJsonResponse(rsp *http.Response) (*GetOpenapiJsonResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetOpenapiJsonResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetVersionResponse parses an HTTP response from a GetVersionWithResponse call
func ParseGetVersionResponse(rsp *http.Response) (*GetVersionResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetVersionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Version
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create compose
	// (POST /compose)
	Compose(w http.ResponseWriter, r *http.Request)
	// The status of a compose
	// (GET /compose/{id})
	ComposeStatus(w http.ResponseWriter, r *http.Request, id string)
	// get the openapi json specification
	// (GET /openapi.json)
	GetOpenapiJson(w http.ResponseWriter, r *http.Request)
	// get the service version
	// (GET /version)
	GetVersion(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Compose operation middleware
func (siw *ServerInterfaceWrapper) Compose(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.Compose(w, r.WithContext(ctx))
}

// ComposeStatus operation middleware
func (siw *ServerInterfaceWrapper) ComposeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	siw.Handler.ComposeStatus(w, r.WithContext(ctx), id)
}

// GetOpenapiJson operation middleware
func (siw *ServerInterfaceWrapper) GetOpenapiJson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.GetOpenapiJson(w, r.WithContext(ctx))
}

// GetVersion operation middleware
func (siw *ServerInterfaceWrapper) GetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.GetVersion(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerFromMux(si, chi.NewRouter())
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	r.Group(func(r chi.Router) {
		r.Post("/compose", wrapper.Compose)
	})
	r.Group(func(r chi.Router) {
		r.Get("/compose/{id}", wrapper.ComposeStatus)
	})
	r.Group(func(r chi.Router) {
		r.Get("/openapi.json", wrapper.GetOpenapiJson)
	})
	r.Group(func(r chi.Router) {
		r.Get("/version", wrapper.GetVersion)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"qs3mPwEAAP//wUQxDOsmAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
