package main

import (
	"github.com/dev-shimada/gostubby/internal/domain/model"
)

// Export functions for testing
var (
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
