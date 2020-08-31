package types

import "github.com/xanzy/go-gitlab"

// Options --
type Options struct {
	client     *gitlab.Client
	sourceType string
	sourceID   string
	dryRun     bool
}

// NewOptions --
func NewOptions(client *gitlab.Client, sourceType, sourceID string, dryrun bool) *Options {
	return &Options{
		client:     client,
		sourceID:   sourceID,
		sourceType: sourceType,
		dryRun:     dryrun,
	}
}
