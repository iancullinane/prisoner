package lambdahttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// Handler wraps h in a function matching the shape lambda.StartHandlerFunc wants.
func Handler(h http.Handler) func(context.Context, events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	return func(ctx context.Context, ev events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		req, err := newRequest(ctx, ev)
		if err != nil {
			return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest}, nil
		}

		rec := &recorder{header: make(http.Header), status: http.StatusOK}
		h.ServeHTTP(rec, req)
		return rec.response(), nil
	}
}

func newRequest(ctx context.Context, ev events.LambdaFunctionURLRequest) (*http.Request, error) {
	body := []byte(ev.Body)
	if ev.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(ev.Body)
		if err != nil {
			return nil, err
		}
		body = decoded
	}

	uri := ev.RawPath
	if ev.RawQueryString != "" {
		uri += "?" + ev.RawQueryString
	}

	req, err := http.NewRequestWithContext(ctx, ev.RequestContext.HTTP.Method, uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for k, v := range ev.Headers {
		req.Header.Set(k, v)
	}
	// Function URLs split cookies out of Headers; net/http wants them rejoined.
	if len(ev.Cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(ev.Cookies, "; "))
	}
	req.RequestURI = uri
	req.Host = req.Header.Get("Host")
	req.RemoteAddr = ev.RequestContext.HTTP.SourceIP
	return req, nil
}

// Satisfy the http.ResponseWriter interface for ServeHttp
// recorder is an http.ResponseWriter that buffers the handler's output instead
// of writing it to a connection.
type recorder struct {
	header      http.Header
	buf         bytes.Buffer
	status      int
	wroteHeader bool
}

func (r *recorder) Header() http.Header { return r.header }

func (r *recorder) WriteHeader(status int) {
	if !r.wroteHeader {
		r.status = status
		r.wroteHeader = true
	}
}

func (r *recorder) Write(b []byte) (int, error) {
	r.WriteHeader(http.StatusOK)
	return r.buf.Write(b)
}

func (r *recorder) response() events.LambdaFunctionURLResponse {
	resp := events.LambdaFunctionURLResponse{
		StatusCode: r.status,
		Headers:    make(map[string]string, len(r.header)),
		Cookies:    r.header.Values("Set-Cookie"),
		Body:       r.buf.String(),
	}
	for k, v := range r.header {
		if k == "Set-Cookie" {
			continue
		}
		resp.Headers[k] = strings.Join(v, ",")
	}
	return resp
}
