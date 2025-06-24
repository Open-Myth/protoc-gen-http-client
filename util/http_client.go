package util

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MinErrStatusCode is the minimum status code that is considered an error.
const MinErrStatusCode = 400

// EncodeHTTPRequest encodes a request to an HTTP request.
func EncodeHTTPRequest[R protoreflect.ProtoMessage](ctx context.Context, path, method string, req R) (*http.Request, error) {
	jsonBytes, err := jsoniter.Marshal(req)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(jsonBytes)
	httpReq, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	injectRequestPathValue(httpReq, req)

	httpReq.Header.Add("Content-Type", "application/json")

	return httpReq, nil
}

// injectRequestPathValue replaces the placeholders in the path of an HTTP request with the actual values.
func injectRequestPathValue(req *http.Request, val any) {
	v := reflect.ValueOf(val).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		jsonTag, _, ok := strings.Cut(field.Tag.Get("json"), ",")
		if !ok {
			continue
		}
		req.URL.Path = strings.ReplaceAll(req.URL.Path, fmt.Sprintf("{%s}", jsonTag), fmt.Sprint(v.Field(i).Interface()))
	}
}

// DecodeHTTPResponse decodes the response body of an HTTP request and returns the decoded result and an error if any.
func DecodeHTTPResponse[R any](resp *http.Response) (*R, error) {
	if resp.StatusCode >= MinErrStatusCode {
		var errResp ErrorResponse
		if err := jsoniter.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}

		return nil, status.Error(codes.Code(resp.StatusCode), errResp.Message)
	}

	// TODO: not implement for body status code
	// Example: {"code": 1, "message": "some error"} and status_code = 200

	var result R
	if err := jsoniter.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ErrorResponse is returned by the server when an error occurs.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
