package testassets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/utils"
)

type HTTPCase interface {
	Run(t *testing.T, baseURL string, ignoredFields []string)
}

type DataManipulation func(any) any

type BaseHTTPCase struct {
	Name                 string
	Path                 string
	PathParams           map[string]any
	QueryParams          map[string]any
	Method               string
	JWTToken             string
	Headers              map[string]string
	ExpectedStatusCode   int
	ExpectedResponseBody any
}

type HTTPWithBodyCase struct {
	BaseHTTPCase
	InputBody        any
	DataManipulation DataManipulation
}

type HTTPWithoutBodyCase struct {
	BaseHTTPCase
}

func (b *BaseHTTPCase) buildFinalPath() string {
	path := replacePathParams(b.Path, b.PathParams)
	if len(b.QueryParams) > 0 {
		path += processQuery(b.QueryParams)
	}
	return path
}

func replacePathParams(path string, params map[string]any) string {
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		path = strings.ReplaceAll(path, placeholder, fmt.Sprintf("%v", value))
	}
	return path
}

func processQuery(values map[string]any) string {
	if len(values) == 0 {
		return ""
	}
	var parts []string
	for key, value := range values {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice {
			var partValues []string
			for i := range v.Len() {
				partValues = append(partValues, fmt.Sprintf("%v", v.Index(i).Interface()))
			}
			parts = append(parts, fmt.Sprintf("%s=%s", key, url.QueryEscape(strings.Join(partValues, ","))))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", key, url.QueryEscape(fmt.Sprintf("%v", value))))
		}
	}
	return "?" + strings.Join(parts, "&")
}

func (b *BaseHTTPCase) prepareRequest(baseURL, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(b.Method, baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", b.JWTToken)
	for key, value := range b.Headers {
		req.Header.Add(key, value)
	}
	return req, nil
}

func (b *BaseHTTPCase) executeRequest(t *testing.T, req *http.Request, ignoredFields []string, expectedResponseBody any) {
	t.Helper()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	marshal, err := json.Marshal(expectedResponseBody)
	if err != nil {
		t.Fatal(err)
	}

	normalizedExpected := utils.Normalize(string(marshal), ignoredFields...)
	normalizedActual := utils.Normalize(string(body), ignoredFields...)

	assert.Equal(t, b.ExpectedStatusCode, resp.StatusCode)
	assert.Equal(t, normalizedExpected, normalizedActual)
}

func (b *BaseHTTPCase) submitRequest(t *testing.T, baseURL, path string, body io.Reader, ignoredFields []string, expectedResponseBody any) {
	t.Helper()
	req, err := b.prepareRequest(baseURL, path, body)
	if err != nil {
		t.Fatal(err)
	}
	b.executeRequest(t, req, ignoredFields, expectedResponseBody)
}

func (tc *HTTPWithBodyCase) Run(t *testing.T, baseURL string, ignoredFields []string) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		inputBody := tc.InputBody
		expectedBody := tc.ExpectedResponseBody

		if tc.DataManipulation != nil {
			inputBody = tc.DataManipulation(tc.InputBody)
			if _, ok := tc.ExpectedResponseBody.(string); !ok {
				expectedBody = tc.DataManipulation(tc.ExpectedResponseBody)
			}
		}

		marshal, err := json.Marshal(inputBody)
		if err != nil {
			t.Fatal(err)
		}

		reqBody := bytes.NewBuffer(marshal)
		finalPath := tc.buildFinalPath()
		tc.submitRequest(t, baseURL, finalPath, reqBody, ignoredFields, expectedBody)
	})
}

func (tt *HTTPWithoutBodyCase) Run(t *testing.T, baseURL string, ignoredFields []string) {
	t.Helper()
	t.Run(tt.Name, func(t *testing.T) {
		finalPath := tt.buildFinalPath()
		tt.submitRequest(t, baseURL, finalPath, http.NoBody, ignoredFields, tt.ExpectedResponseBody)
	})
}
