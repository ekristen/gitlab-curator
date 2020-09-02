package types

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

// GenerateListProjectIssuesOptions --
func (rule *Rule) GenerateListProjectIssuesOptions() *gitlab.ListProjectIssuesOptions {
	options := &gitlab.ListProjectIssuesOptions{
		ListOptions: rule.GenerateLimits(),
	}

	options.State = &rule.Conditions.State
	options.Labels = rule.Conditions.Labels
	options.NotLabels = rule.Conditions.ForbiddenLabels

	if rule.Conditions.Date != nil {
		dur := -rule.Conditions.Date.Duration.Duration
		now := time.Now().UTC()
		olderThan := now.Add(dur)

		if rule.Conditions.Date.Condition == "created_before" {
			options.CreatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "created_after" {
			options.CreatedAfter = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_before" {
			options.UpdatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_after" {
			options.UpdatedAfter = &olderThan
		}
	}

	if rule.Limits != nil && rule.Limits.MostRecent != nil {
		options.OrderBy = gitlab.String("created_at")
		options.Sort = gitlab.String("desc")
	}

	return options
}

// GenerateGroupIssueOptions --
func (rule *Rule) GenerateGroupIssueOptions() *gitlab.ListGroupIssuesOptions {
	options := &gitlab.ListGroupIssuesOptions{
		ListOptions: rule.GenerateLimits(),
	}

	options.State = &rule.Conditions.State
	options.Labels = rule.Conditions.Labels
	options.NotLabels = rule.Conditions.ForbiddenLabels

	if rule.Conditions.MissingLabels != nil {
		options.NotLabels = append(options.NotLabels, rule.Conditions.MissingLabels...)
	}

	if rule.Conditions.Milestone != "" {
		options.Milestone = &rule.Conditions.Milestone
	}

	if rule.Conditions.Date != nil {
		dur := -rule.Conditions.Date.Duration.Duration
		now := time.Now().UTC()
		olderThan := now.Add(dur)

		if rule.Conditions.Date.Condition == "created_before" {
			options.CreatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "created_after" {
			options.CreatedAfter = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_before" {
			options.UpdatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_after" {
			options.UpdatedAfter = &olderThan
		}
	}

	if rule.Limits != nil && rule.Limits.MostRecent != nil {
		options.OrderBy = gitlab.String("created_at")
		options.Sort = gitlab.String("desc")
	}

	return options
}

// GenerateLimits --
func (rule *Rule) GenerateLimits() gitlab.ListOptions {
	listOptions := gitlab.ListOptions{}
	if rule.Limits != nil {
		if rule.Limits.PerPage != nil {
			listOptions.PerPage = *rule.Limits.PerPage
		}

		if rule.Limits.MostRecent != nil {
			listOptions.PerPage = *rule.Limits.MostRecent
		}
	}
	return listOptions
}

// GenerateGroupMergeRequestsOptions --
func (rule *Rule) GenerateGroupMergeRequestsOptions() *gitlab.ListGroupMergeRequestsOptions {
	options := &gitlab.ListGroupMergeRequestsOptions{
		ListOptions: rule.GenerateLimits(),
	}

	labels := gitlab.Labels(rule.Conditions.Labels)
	notLabels := gitlab.Labels(rule.Conditions.ForbiddenLabels)

	options.State = &rule.Conditions.State
	options.Labels = &labels
	options.NotLabels = &notLabels

	if rule.Conditions.Date != nil {
		dur := -rule.Conditions.Date.Duration.Duration
		now := time.Now().UTC()
		olderThan := now.Add(dur)

		if rule.Conditions.Date.Condition == "created_before" {
			options.CreatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "created_after" {
			options.CreatedAfter = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_before" {
			options.UpdatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_after" {
			options.UpdatedAfter = &olderThan
		}
	}

	if rule.Limits != nil && rule.Limits.MostRecent != nil {
		options.OrderBy = gitlab.String("created_at")
		options.Sort = gitlab.String("desc")
	}

	return options
}

// GenerateProjectMergeRequestsOptions --
func (rule *Rule) GenerateProjectMergeRequestsOptions() *gitlab.ListProjectMergeRequestsOptions {
	options := &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: rule.GenerateLimits(),
	}

	labels := gitlab.Labels(rule.Conditions.Labels)
	notLabels := gitlab.Labels(rule.Conditions.ForbiddenLabels)

	options.State = &rule.Conditions.State
	options.Labels = labels
	options.NotLabels = notLabels

	if rule.Conditions.Date != nil {
		dur := -rule.Conditions.Date.Duration.Duration
		now := time.Now().UTC()
		olderThan := now.Add(dur)

		if rule.Conditions.Date.Condition == "created_before" {
			options.CreatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "created_after" {
			options.CreatedAfter = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_before" {
			options.UpdatedBefore = &olderThan
		} else if rule.Conditions.Date.Condition == "updated_after" {
			options.UpdatedAfter = &olderThan
		}
	}

	if rule.Limits != nil && rule.Limits.MostRecent != nil {
		options.OrderBy = gitlab.String("created_at")
		options.Sort = gitlab.String("desc")
	}

	return options
}

// GenerateListGroupMilestonesOptions --
func (rule *Rule) GenerateListGroupMilestonesOptions() *gitlab.ListGroupMilestonesOptions {
	options := &gitlab.ListGroupMilestonesOptions{
		ListOptions: rule.GenerateLimits(),
	}

	if rule.Conditions != nil {
		if rule.Conditions.State != "" {
			options.State = rule.Conditions.State
		}
	}

	return options
}
