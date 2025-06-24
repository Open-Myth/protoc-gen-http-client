package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/require"
)

func Test_getClientStruct(t *testing.T) {
	gf := jen.NewFile("test")
	got := getClientStruct("Hello")

	expected := `// HTTPClient is a http client for the Hello service
type HelloHTTPClient struct {
	BaseURL      string
	roundTripper http.RoundTripper
} 
`

	buf := bytes.Buffer{}
	err := gf.Add(got).Render(&buf)

	require.NoError(t, err)
	require.Equal(t, expected, buf.String())
}

func Test_getNewClient(t *testing.T) {
	got := getNewClient("Hello")

	expected := `func NewHelloHTTPClient(baseURL string) *HelloHTTPClient {
	return &HelloHTTPClient{
		baseURL:      baseURL,
		roundTripper: httpclient.NewRoundTripper(),
	}
}`

	buf := bytes.Buffer{}
	err := got.Render(&buf)

	require.NoError(t, err)
	require.Equal(t, expected, buf.String())
}

func Test_getMethod(t *testing.T) {

	methodName := "SayHello"
	serviceName := "Greeting"
	reqName := "HelloRequest"
	respName := "HelloResponse"
	path := "/hello"
	httpMethod := "POST"
	got := getMethod(methodName, serviceName, reqName, respName, path, httpMethod)
	buf := bytes.Buffer{}
	err := got.Render(&buf)

	expected := fmt.Sprintf(`// %s is a http call method for the %s service
func (c *%[2]sHTTPClient) %[1]s(ctx context.Context, reqData *%[3]s) (*%[4]s, error) {
	path, err := url.JoinPath(c.BaseURL, "%s")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Errorf("path is not valid: %%w", err).Error())
	}

	reqClient, err := util.EncodeHTTPRequest(ctx, path, "%s", reqData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Errorf("unable to encode http request: %%w", err).Error())
	}
	client := http.Client{}
	resp, err := client.Do(reqClient)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("unable to request: %%w", err).Error())
	}

	return util.DecodeHTTPResponse[%[4]s](resp)
} 
`, methodName, serviceName, reqName, respName, path, httpMethod)
	require.NoError(t, err)
	require.Equal(t, expected, buf.String())
}
