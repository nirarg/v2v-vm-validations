package persistent

// This package provides a public API bridge to the internal persistent package.

import (
	"github.com/nirarg/v2v-vm-validations/internal/persistent"
)

// Re-export persistent types
type (
	Inspector   = persistent.Inspector
	Credentials = persistent.Credentials
	CacheKey    = persistent.CacheKey
	DB          = persistent.DB
)

// Re-export constructor functions
var (
	NewInspector = persistent.NewInspector
)
