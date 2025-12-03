package inspection

// This package provides a public API bridge to the internal inspection package.
// It re-exports types and functions to allow external usage temporarily.
// Once vm-deep-inspection-demo is fully migrated, this can be removed and
// inspection can remain internal-only.

import (
	"github.com/nirarg/v2v-vm-validations/internal/inspection"
)

// Re-export inspector types
type (
	VirtInspector    = inspection.VirtInspector
	VirtV2vInspector = inspection.VirtV2vInspector
	NBDKitSession    = inspection.NBDKitSession
	V2VSession       = inspection.V2VSession
)

// Re-export constructor functions
var (
	NewVirtInspector    = inspection.NewVirtInspector
	NewVirtV2vInspector = inspection.NewVirtV2vInspector
)

// Re-export constants
const (
	UseVirtV2VOpen = inspection.UseVirtV2VOpen
)
