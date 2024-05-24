package web

import "net/http"

// Rotate returns a new HTTP adapter that intercept outgoing HTTP requests and
// forward them through a proxy created via AWS API Gateway.
func Rotate(r *http.Request) string {

	return ""
}
