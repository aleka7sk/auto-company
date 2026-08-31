package assets

import "embed"

// FS contains the deterministic project templates, profiles, and integration registry
// shipped with Auto Company. Keeping them embedded makes `go install` distributions
// behave the same as a source checkout.
//
//go:embed templates/*.tmpl profiles/*.json registry/*.json
var FS embed.FS
