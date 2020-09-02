package types

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// Filter --
type Filter struct {
	Name       string          `json:"name" yaml:"name"`
	Relation   string          `json:"relation" yaml:"relation" default:"self"`
	Conditions *RuleConditions `json:"conditions" yaml:"conditions"`
	Action     string          `json:"action" yaml:"action" default:"include"`
}

// GroupMilestones --
func (f *Filter) GroupMilestones(opts *Options, milestones []*gitlab.GroupMilestone, log *logrus.Entry) ([]*gitlab.GroupMilestone, error) {
	if f.Relation == "" {
		f.Relation = "self"
	}
	if f.Action == "" {
		f.Action = "include"
	}

	log = log.
		WithField("relation", f.Relation).
		WithField("action", f.Action).
		WithField("name", f.Name)

	log.WithField("prefilter_count", len(milestones)).Debug("prefilter count milestones")

	var filteredMilestones []*gitlab.GroupMilestone

	for _, milestone := range milestones {
		matched := false

		log = log.WithField("milestone", milestone.Title)

		if f.Relation == "self" {
			log.Debug("relation: self - called")
			if f.Conditions != nil && f.Conditions.Expired == milestone.Expired {
				log.Debug("expired matched")
				matched = true
			}
		} else if f.Relation == "assigned_issues" {
			issues, _, err := opts.client.GroupMilestones.GetGroupMilestoneIssues(milestone.GroupID, milestone.ID, &gitlab.GetGroupMilestoneIssuesOptions{})
			if err != nil {
				return milestones, err
			}

			if len(issues) > 0 {
				log.WithField("issue_count", len(issues)).Debug("found issues for milestone")

				for _, i := range issues {
					if f.Conditions != nil && f.Conditions.State == i.State {
						log.Debug("issue state matched")
						//filteredMilestones = append(filteredMilestones, milestone)
						matched = true
					}
				}
			}
		}

		if f.Action == "include" && matched {
			filteredMilestones = append(filteredMilestones, milestone)
		} else if f.Action == "exclude" && !matched {
			filteredMilestones = append(filteredMilestones, milestone)
		}
	}

	log.WithField("filtered_count", len(filteredMilestones)).Debug("filtered milestones")

	return filteredMilestones, nil
}
