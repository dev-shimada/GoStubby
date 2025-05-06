package main

import (
	"net/url"

	"github.com/dev-shimada/gostubby/internal/domain/model"
)

// Export functions for testing
var (
	ExportPathMatcher = func(e model.Endpoint, rawPath, path string) (bool, map[string]string) {
		return pathMatcher(e, rawPath, path)
	}
	ExportQueryMatcher = func(e model.Endpoint, q url.Values) bool {
		return queryMatcher(e, q)
	}
	// Export handle function for testing
	Handle = handle
)

// Export types for testing
type (
	Endpoint = model.Endpoint
	Request  = model.Request
	Response = model.Response
	Matcher  = model.Matcher
)
