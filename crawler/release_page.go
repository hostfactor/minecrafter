package crawler

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
)

type ReleaseParser interface {
	ParseRelease(e *colly.HTMLElement) *Release
}

type Release struct {
	// The semver Parsed Name.
	Version *semver.Version

	// The Value used within the edition. Sometimes the Name and the value used internally can be different,
	Value string

	// The human-readable name for the release. This will often be a semver e.g. 1.18.1.
	Name string

	ArtifactURL string
}
