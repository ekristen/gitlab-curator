package types

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v1"
)

// Policy --
type Policy struct {
	ResourceRules *ResourceRules `json:"resource_rules" yaml:"resource_rules"`
}

// ResourcePolicy --
type ResourcePolicy struct {
	Rules []Rule `json:"rules" yaml:"rules"`
}

// ResourceRules --
type ResourceRules struct {
	Epics         *ResourcePolicy `json:"epics,omitempty" yaml:"epics,omitempty"`
	Issues        *ResourcePolicy `json:"issues,omitempty" yaml:"issues,omitempty"`
	MergeRequests *ResourcePolicy `json:"merge_requests,omitempty" yaml:"merge_requests,omitempty"`
	Milestones    *ResourcePolicy `json:"milestones,omitempty" yaml:"milestones,omitempty"`
}

// Process --
func (r *ResourceRules) Process(opts *Options) error {
	if r.Milestones != nil {
		logrus.WithField("rule_count", len(r.Milestones.Rules)).Info("executing milestones based policy rules")
		for _, mr := range r.Milestones.Rules {
			logrus.WithField("name", mr.Name).WithField("type", opts.sourceType).Info("executing rule")

			generatedIssuesRules, err := mr.ProcessMilestones(opts)
			if err != nil {
				return err
			}

			b, _ := yaml.Marshal(generatedIssuesRules)
			logrus.Debug(string(b))

			if generatedIssuesRules != nil {
				logrus.Debug("found generated issue rules, appending ...")

				if r.Issues == nil {
					r.Issues = &ResourcePolicy{
						Rules: []Rule{},
					}
				}
				r.Issues.Rules = append(r.Issues.Rules, generatedIssuesRules...)
			}
		}
	}

	if r.Issues != nil {
		logrus.WithField("rule_count", len(r.Issues.Rules)).Info("executing issues based policy rules")

		for _, r := range r.Issues.Rules {
			logrus.WithField("name", r.Name).WithField("type", opts.sourceType).Info("executing rule")

			if err := r.ProcessIssues(opts); err != nil {
				return err
			}
		}
	}

	if r.MergeRequests != nil {
		logrus.WithField("rule_count", len(r.MergeRequests.Rules)).Info("executing merge requests based policy rules")
		for _, r := range r.MergeRequests.Rules {
			logrus.WithField("name", r.Name).WithField("type", opts.sourceType).Info("executing rule")

			if err := r.ProcessMergeRequests(opts); err != nil {
				return err
			}
		}
	}

	return nil
}
