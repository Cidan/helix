// Package helix provides an client for the Twitch Helix API.
package helix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultAPIBaseURL is the base URL for composing API requests.
	DefaultAPIBaseURL = "https://api.twitch.tv/helix"

	// AuthBaseURL is the base URL for composing authentication requests.
	AuthBaseURL = "https://id.twitch.tv/oauth2"
)

type Option func(*Options) error

type httpDo func(*http.Request) (*http.Response, error)

type Client struct {
	mu           sync.RWMutex
	opts         *Options
	lastResponse *Response
}

type Options struct {
	clientID        string
	clientSecret    string
	appAccessToken  string
	userAccessToken string
	userAgent       string
	redirectURI     string
	httpDo          httpDo
	rateLimitFunc   RateLimitFunc
	apiBaseURL      string
	extensionOpts   ExtensionOptions
}

type ExtensionOptions struct {
	OwnerUserID    string
	Secret         string
	SignedJWTToken string
}

// DateRange is a generic struct used by various responses.
type DateRange struct {
	StartedAt Time `json:"started_at"`
	EndedAt   Time `json:"ended_at"`
}

type RateLimitFunc func(*Response) error

type ResponseCommon struct {
	StatusCode   int
	Header       http.Header
	Error        string `json:"error"`
	ErrorStatus  int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func (rc *ResponseCommon) convertHeaderToInt(str string) int {
	i, _ := strconv.Atoi(str)

	return i
}

// GetRateLimit returns the "RateLimit-Limit" header as an int.
func (rc *ResponseCommon) GetRateLimit() int {
	return rc.convertHeaderToInt(rc.Header.Get("RateLimit-Limit"))
}

// GetRateLimitRemaining returns the "RateLimit-Remaining" header as an int.
func (rc *ResponseCommon) GetRateLimitRemaining() int {
	return rc.convertHeaderToInt(rc.Header.Get("RateLimit-Remaining"))
}

// GetRateLimitReset returns the "RateLimit-Reset" header as an int.
func (rc *ResponseCommon) GetRateLimitReset() int {
	return rc.convertHeaderToInt(rc.Header.Get("RateLimit-Reset"))
}

type Response struct {
	ResponseCommon
	Data interface{}
}

// HydrateResponseCommon copies the content of the source response's ResponseCommon to the supplied ResponseCommon argument
func (r *Response) HydrateResponseCommon(rc *ResponseCommon) {
	rc.StatusCode = r.ResponseCommon.StatusCode
	rc.Header = r.ResponseCommon.Header
	rc.Error = r.ResponseCommon.Error
	rc.ErrorStatus = r.ResponseCommon.ErrorStatus
	rc.ErrorMessage = r.ResponseCommon.ErrorMessage
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

func defaultClientOpts() *Options {
	return &Options{
		userAgent:   "Twitch Helix API Go library github.com/nicklaw5/helix",
		redirectURI: DefaultAPIBaseURL,
		httpDo:      http.DefaultClient.Do,
		apiBaseURL:  DefaultAPIBaseURL,
	}
}

func WithHttpDo(f httpDo) Option {
	return func(o *Options) error {
		o.httpDo = f
		return nil
	}
}

func WithClientID(v string) Option {
	return func(o *Options) error {
		o.clientID = v
		return nil
	}
}

func WithClientSecret(v string) Option {
	return func(o *Options) error {
		o.clientSecret = v
		return nil
	}
}

func WithRedirectURI(v string) Option {
	return func(o *Options) error {
		o.redirectURI = v
		return nil
	}
}

func WithUserAccessToken(v string) Option {
	return func(o *Options) error {
		o.userAccessToken = v
		return nil
	}
}

func WithAppAccessToken(v string) Option {
	return func(o *Options) error {
		o.appAccessToken = v
		return nil
	}
}

func WithUserAgent(v string) Option {
	return func(o *Options) error {
		o.userAgent = v
		return nil
	}
}

func WithExtensionOptions(v ExtensionOptions) Option {
	return func(o *Options) error {
		o.extensionOpts = v
		return nil
	}
}

func WithRateLimitFunc(v RateLimitFunc) Option {
	return func(o *Options) error {
		o.rateLimitFunc = v
		return nil
	}
}

func WithAPIBaseURL(v string) Option {
	return func(o *Options) error {
		o.apiBaseURL = v
		return nil
	}
}

// NewClient returns a new Twitch Helix API client. It returns an
// if clientID is an empty string. It is concurrency safe.
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	c := &Client{
		opts: defaultClientOpts(),
	}

	for _, opt := range opts {
		if err := opt(c.opts); err != nil {
			return nil, err
		}
	}

	if c.opts.clientID == "" {
		return nil, errors.New("A client ID was not provided but is required")
	}

	return c, nil
}

func (c *Client) get(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodGet, path, respData, reqData, false, opts)
}

func (c *Client) post(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodPost, path, respData, reqData, false, opts)
}

func (c *Client) put(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodPut, path, respData, reqData, false, opts)
}

func (c *Client) delete(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodDelete, path, respData, reqData, false, opts)
}

func (c *Client) patchAsJSON(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodPatch, path, respData, reqData, true, opts)
}

func (c *Client) postAsJSON(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodPost, path, respData, reqData, true, opts)
}

func (c *Client) putAsJSON(ctx context.Context, path string, respData, reqData interface{}, opts []Option) (*Response, error) {
	return c.sendRequest(ctx, http.MethodPut, path, respData, reqData, true, opts)
}

func (c *Client) sendRequest(ctx context.Context, method, path string, respData, reqData interface{}, hasJSONBody bool, opts []Option) (*Response, error) {
	resp := &Response{}
	if respData != nil {
		resp.Data = respData
	}

	req, err := c.newRequest(method, path, reqData, hasJSONBody)
	if err != nil {
		return nil, err
	}

	err = c.doRequest(req.WithContext(ctx), resp, opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func buildQueryString(req *http.Request, v interface{}) (string, error) {
	isNil, err := isZero(v)
	if err != nil {
		return "", err
	}

	if isNil {
		return "", nil
	}

	query := req.URL.Query()
	vType := reflect.TypeOf(v).Elem()
	vValue := reflect.ValueOf(v).Elem()

	for i := 0; i < vType.NumField(); i++ {
		var defaultValue string

		field := vType.Field(i)
		tag := field.Tag.Get("query")

		// Get the default value from the struct tag
		if strings.Contains(tag, ",") {
			tagSlice := strings.Split(tag, ",")

			tag = tagSlice[0]
			defaultValue = tagSlice[1]
		}

		if field.Type.Kind() == reflect.Slice {
			// Attach any slices as query params
			fieldVal := vValue.Field(i)
			for j := 0; j < fieldVal.Len(); j++ {
				query.Add(tag, fmt.Sprintf("%v", fieldVal.Index(j)))
			}
		} else if isDatetimeTagField(tag) {
			// Get and correctly format datetime fields, and attach them query params
			dateStr := fmt.Sprintf("%v", vValue.Field(i))

			if strings.Contains(dateStr, " m=") {
				datetimeSplit := strings.Split(dateStr, " m=")
				dateStr = datetimeSplit[0]
			}

			date, err := time.Parse(requestDateTimeFormat, dateStr)
			if err != nil {
				return "", err
			}

			// Determine if the date has been set. If it has we'll add it to the query.
			if !date.IsZero() {
				query.Add(tag, date.Format(time.RFC3339))
			}
		} else {
			// Add any scalar values as query params
			fieldVal := fmt.Sprintf("%v", vValue.Field(i))

			// If no value was set by the user, use the default
			// value specified in the struct tag.
			if fieldVal == "" || fieldVal == "0" {
				if defaultValue == "" {
					continue
				}

				fieldVal = defaultValue
			}

			query.Add(tag, fieldVal)
		}
	}

	return query.Encode(), nil
}

func isZero(v interface{}) (bool, error) {
	t := reflect.TypeOf(v)
	if !t.Comparable() {
		return false, fmt.Errorf("type is not comparable: %v", t)
	}
	return v == reflect.Zero(t).Interface(), nil
}

func (c *Client) newRequest(method, path string, data interface{}, hasJSONBody bool) (*http.Request, error) {
	url := c.getBaseURL(path) + path

	if hasJSONBody {
		return c.newJSONRequest(method, url, data)
	}

	return c.newStandardRequest(method, url, data)
}

func (c *Client) newStandardRequest(method, url string, data interface{}) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return req, nil
	}

	query, err := buildQueryString(req, data)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query

	return req, nil
}

func (c *Client) newJSONRequest(method, url string, data interface{}) (*http.Request, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	query, err := buildQueryString(req, data)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) getBaseURL(path string) string {
	for _, authPath := range authPaths {
		if strings.Contains(path, authPath) {
			return AuthBaseURL
		}
	}

	return c.opts.apiBaseURL
}

func (c *Client) mergeOptions(opts []Option) (*Options, error) {
	o := *c.opts
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}
	return &o, nil
}

func (c *Client) doRequest(req *http.Request, resp *Response, opts []Option) error {
	o, err := c.mergeOptions(opts)
	if err != nil {
		return err
	}

	c.setRequestHeaders(req, o)

	rateLimitFunc := c.opts.rateLimitFunc

	for {
		if c.lastResponse != nil && rateLimitFunc != nil {
			err := rateLimitFunc(c.lastResponse)
			if err != nil {
				return err
			}
		}

		response, err := c.opts.httpDo(req)
		if err != nil {
			return fmt.Errorf("Failed to execute API request: %s", err.Error())
		}
		defer response.Body.Close()

		resp.Header = response.Header

		setResponseStatusCode(resp, "StatusCode", response.StatusCode)

		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		// Only attempt to decode the response if we have a response we can handle
		if len(bodyBytes) > 0 && resp.StatusCode < http.StatusInternalServerError {
			if resp.Data != nil && resp.StatusCode < http.StatusBadRequest {
				// Successful request
				err = json.Unmarshal(bodyBytes, &resp.Data)
			} else {
				// Failed request
				err = json.Unmarshal(bodyBytes, &resp)
			}

			if err != nil {
				return fmt.Errorf("Failed to decode API response: %s", err.Error())
			}
		}

		if rateLimitFunc == nil {
			break
		} else {
			c.mu.Lock()
			c.lastResponse = resp
			c.mu.Unlock()

			if rateLimitFunc != nil &&
				c.lastResponse.StatusCode == http.StatusTooManyRequests {
				// Rate limit exceeded, retry to send request after
				// applying rate limiter callback
				continue
			}

			break
		}
	}

	return nil
}

func (c *Client) setRequestHeaders(req *http.Request, opts *Options) {
	req.Header.Set("Client-ID", opts.clientID)

	if opts.userAgent != "" {
		req.Header.Set("User-Agent", opts.userAgent)
	}

	var bearerToken string
	if opts.appAccessToken != "" {
		bearerToken = opts.appAccessToken
	}
	if opts.userAccessToken != "" {
		bearerToken = opts.userAccessToken
	}
	if opts.extensionOpts.SignedJWTToken != "" {
		bearerToken = opts.extensionOpts.SignedJWTToken
	}

	authType := "Bearer"
	// Token validation requires different type of Auth
	if req.URL.String() == AuthBaseURL+authPaths["validate"] {
		authType = "OAuth"
	}

	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", authType, bearerToken))
	}
}

func setResponseStatusCode(v interface{}, fieldName string, code int) {
	s := reflect.ValueOf(v).Elem()
	field := s.FieldByName(fieldName)
	field.SetInt(int64(code))
}

// GetAppAccessToken returns the current app access token.
func (c *Client) GetAppAccessToken() string {
	return c.opts.appAccessToken
}

func (c *Client) SetAppAccessToken(accessToken string) {
	c.opts.appAccessToken = accessToken
}

// GetUserAccessToken returns the current user access token.
func (c *Client) GetUserAccessToken() string {
	return c.opts.userAccessToken
}

func (c *Client) SetUserAccessToken(accessToken string) {
	c.opts.userAccessToken = accessToken
}

// GetAppAccessToken returns the current app access token.
func (c *Client) GetExtensionSignedJWTToken() string {
	return c.opts.extensionOpts.SignedJWTToken
}

func (c *Client) SetExtensionSignedJWTToken(jwt string) {
	c.opts.extensionOpts.SignedJWTToken = jwt
}

func (c *Client) SetUserAgent(userAgent string) {
	c.opts.userAgent = userAgent
}

func (c *Client) SetRedirectURI(uri string) {
	c.opts.redirectURI = uri
}
