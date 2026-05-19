package testassets

import (
	"fmt"
	"maps"
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/utils"
)

type HTTPCasesRunner struct {
	Cases       []HTTPCase
	baseURL     string
	httpMethod  string
	jwtToken    string
	headers     map[string]string
	queryParams map[string]any
	pathParams  map[string]any
}

func NewHTTPCasesRunner(baseURL, httpMethod, jwtToken string) *HTTPCasesRunner {
	return &HTTPCasesRunner{
		Cases:       []HTTPCase{},
		baseURL:     baseURL,
		httpMethod:  httpMethod,
		jwtToken:    jwtToken,
		headers:     make(map[string]string),
		queryParams: make(map[string]any),
		pathParams:  make(map[string]any),
	}
}

func (h *HTTPCasesRunner) AddHeader(name string, value any) *HTTPCasesRunner {
	h.headers[name] = fmt.Sprint(value)
	return h
}

func (h *HTTPCasesRunner) AddQueryParams(params map[string]any) *HTTPCasesRunner {
	maps.Copy(h.queryParams, params)
	return h
}

func (h *HTTPCasesRunner) AddPathParams(params map[string]any) *HTTPCasesRunner {
	maps.Copy(h.pathParams, params)
	return h
}

func (h *HTTPCasesRunner) newBaseCase(name string, expectedStatusCode int, expectedResponse any) BaseHTTPCase {
	return BaseHTTPCase{
		Name:                 name,
		Path:                 h.baseURL,
		PathParams:           h.pathParams,
		QueryParams:          h.queryParams,
		Method:               h.httpMethod,
		JWTToken:             h.jwtToken,
		Headers:              h.headers,
		ExpectedStatusCode:   expectedStatusCode,
		ExpectedResponseBody: expectedResponse,
	}
}

func (h *HTTPCasesRunner) addCase(tc HTTPCase) *HTTPCasesRunner {
	h.Cases = append(h.Cases, tc)
	return h
}

func (h *HTTPCasesRunner) Run(t *testing.T, url string, ignoredFields []string) {
	t.Helper()
	for _, tc := range h.Cases {
		tc.Run(t, url, ignoredFields)
	}
}

func (h *HTTPCasesRunner) NewRequestWithBody(name string, inputBody any, statusCode int, responseBody any) *HTTPCasesRunner {
	tc := &HTTPWithBodyCase{
		BaseHTTPCase: h.newBaseCase(name, statusCode, responseBody),
		InputBody:    inputBody,
	}
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewOKRequestWithBody(inputBody any) *HTTPCasesRunner {
	tc := &HTTPWithBodyCase{
		BaseHTTPCase: h.newBaseCase("OK", http.StatusOK, inputBody),
		InputBody:    inputBody,
	}
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewRequestWithDataManipulation(name string, inputBody any, manipulation DataManipulation, statusCode int, responseBody any) *HTTPCasesRunner {
	tc := &HTTPWithBodyCase{
		BaseHTTPCase:     h.newBaseCase(name, statusCode, responseBody),
		InputBody:        inputBody,
		DataManipulation: manipulation,
	}
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewUnmarshalErrorRequest() *HTTPCasesRunner {
	tc := &HTTPWithBodyCase{
		BaseHTTPCase: h.newBaseCase("UnmarshalError", http.StatusBadRequest, ExpectedErrorUnmarshal),
		InputBody:    GarbageString,
	}
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewNilBodyRequest() *HTTPCasesRunner {
	tc := &HTTPWithBodyCase{
		BaseHTTPCase: h.newBaseCase("NilBody", http.StatusBadRequest, ExpectedResponseEmptyBody),
		InputBody:    nil,
	}
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewBadJWTRequest() *HTTPCasesRunner {
	tc := &HTTPWithBodyCase{
		BaseHTTPCase: h.newBaseCase("BadJWT", http.StatusUnauthorized, ExpectedResponseBadJWT),
		InputBody:    nil,
	}
	tc.JWTToken = EmptyString
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewOKRequestWithoutBody(queryParams, pathParams map[string]any, responseBody any) *HTTPCasesRunner {
	tc := &HTTPWithoutBodyCase{
		BaseHTTPCase: h.newBaseCase("OK", http.StatusOK, responseBody),
	}
	tc.QueryParams = utils.MergeMaps(h.queryParams, queryParams)
	tc.PathParams = utils.MergeMaps(h.pathParams, pathParams)
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewRequestWithoutBody(name string, queryParams, pathParams map[string]any, statusCode int, responseBody any) *HTTPCasesRunner {
	tc := &HTTPWithoutBodyCase{
		BaseHTTPCase: h.newBaseCase(name, statusCode, responseBody),
	}
	tc.QueryParams = utils.MergeMaps(h.queryParams, queryParams)
	tc.PathParams = utils.MergeMaps(h.pathParams, pathParams)
	return h.addCase(tc)
}

func (h *HTTPCasesRunner) NewNoAccessRequest(queryParams, pathParams map[string]any) *HTTPCasesRunner {
	tc := &HTTPWithoutBodyCase{
		BaseHTTPCase: h.newBaseCase("NoAccess", http.StatusForbidden, ExpectedResponseHaveNoAccess),
	}
	tc.QueryParams = utils.MergeMaps(h.queryParams, queryParams)
	tc.PathParams = utils.MergeMaps(h.pathParams, pathParams)
	tc.JWTToken = EmptyToken
	return h.addCase(tc)
}
