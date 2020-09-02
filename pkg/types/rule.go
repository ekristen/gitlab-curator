package types

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// ProcessIssues --
func (r *Rule) ProcessIssues(opts *Options) error {
	var err error
	var issues []*gitlab.Issue

	if opts.sourceType == "group" {
		options := r.GenerateGroupIssueOptions()
		issues, _, err = opts.client.Issues.ListGroupIssues(opts.sourceID, options)
	} else if opts.sourceType == "project" {
		options := r.GenerateListProjectIssuesOptions()
		issues, _, err = opts.client.Issues.ListProjectIssues(opts.sourceID, options)
	}
	if err != nil {
		return err
	}

	logrus.WithField("name", r.Name).WithField("count", len(issues)).Info("found issues")

	if err := r.Summarize(opts, issues); err != nil {
		return err
	}

	for _, i := range issues {
		if err := r.Label(opts, i); err != nil {
			return err
		}
		if err := r.Comment(opts, i); err != nil {
			return err
		}
		if err := r.State(opts, i); err != nil {
			return err
		}
	}

	return nil
}

// ProcessMergeRequests --
func (r *Rule) ProcessMergeRequests(opts *Options) error {
	var err error
	var mergeRequests []*gitlab.MergeRequest

	if opts.sourceType == "group" {
		options := &gitlab.ListGroupMergeRequestsOptions{}
		mergeRequests, _, err = opts.client.MergeRequests.ListGroupMergeRequests(opts.sourceID, options)
	} else if opts.sourceType == "project" {
		options := &gitlab.ListProjectMergeRequestsOptions{}
		mergeRequests, _, err = opts.client.MergeRequests.ListProjectMergeRequests(opts.sourceID, options)
	}
	if err != nil {
		return err
	}

	for _, mr := range mergeRequests {
		if err := r.Comment(opts, mr); err != nil {
			return err
		}
	}

	return nil
}

// ProcessMilestones --
func (r *Rule) ProcessMilestones(opts *Options) error {
	log := logrus.WithField("component", "process-milestones").WithField("source-type", opts.sourceType)

	if opts.sourceType == "group" {
		options := r.GenerateListGroupMilestonesOptions()
		milestones, _, err := opts.client.GroupMilestones.ListGroupMilestones(opts.sourceID, options)
		if err != nil {
			return err
		}

		log.WithField("count", len(milestones)).Info("found milestones")

		milestones, err = r.FilterGroupMilestones(opts, milestones, log)
		if err != nil {
			return err
		}

		log.WithField("count", len(milestones)).Info("milestones after filtering")

		for _, milestone := range milestones {
			if err := r.State(opts, milestone); err != nil {
				return err
			}
		}
	} else if opts.sourceType == "project" {
		options := &gitlab.ListMilestonesOptions{}
		milestones, _, err := opts.client.Milestones.ListMilestones(opts.sourceID, options)
		if err != nil {
			return err
		}

		for _, milestone := range milestones {
			if err := r.State(opts, milestone); err != nil {
				return err
			}
		}
	}

	return nil
}

// FilterGroupMilestones --
func (r *Rule) FilterGroupMilestones(opts *Options, milestones []*gitlab.GroupMilestone, log *logrus.Entry) ([]*gitlab.GroupMilestone, error) {
	if r.Filters == nil {
		return milestones, nil
	}

	log.Info("filtering milestones")

	var err error
	var filteredMilestones = milestones

	for _, filter := range r.Filters {
		filteredMilestones, err = filter.GroupMilestones(opts, filteredMilestones, log)
		if err != nil {
			return nil, err
		}
	}

	return filteredMilestones, nil
}
