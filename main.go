package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/vefaas"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
)

func main() {
	// Start your vefaas function =D.
	vefaas.Start(handler)
}

// Define your handler function.
func handler(ctx context.Context, r *events.HTTPRequest) (*events.EventResponse, error) {
	fmt.Printf("received new request: %s %s, request id: %s\n", r.HTTPMethod, r.Path, vefaascontext.RequestIdFromContext(ctx))
	fmt.Printf("debug request: header:%v, body:%s\n", r.Headers, r.Body)
	ret := make(map[string]interface{})
	ret["content"] = "Hello veFaaS from louyuting-v1!"
	ret["http_method"] = r.HTTPMethod
	ret["http_path"] = r.Path
	query := make(map[string]interface{})
	header := make(map[string]interface{})
	for k, v := range r.QueryStringParameters {
		query[k] = v
	}
	for k, v := range r.Headers {
		header[k] = v
	}
	ret["http_query"] = query
	ret["http_header"] = header
	ret["http_body"] = string(r.Body)
	ret["now_time"] = time.Now().Format(time.RFC3339)
	retBody, _ := json.Marshal(ret)
	return &events.EventResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: retBody,
	}, nil
}
