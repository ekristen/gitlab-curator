package types

import "github.com/sirupsen/logrus"

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
	Issues        *ResourcePolicy `json:"issues,omitempty" yaml:"issues,omitempty"`
	MergeRequests *ResourcePolicy `json:"merge_requests,omitempty" yaml:"merge_requests,omitempty"`
}

// Process --
func (r *ResourceRules) Process(opts *Options) error {
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

			/*
				if sourceType == "group" {
					options := r.GenerateGroupMergeRequestsOptions()

					mergeRequests, _, err := git.MergeRequests.ListGroupMergeRequests(sourceID, options)
					if err != nil {
						return err
					}

					logrus.WithField("count", len(mergeRequests)).Info("found merge requests")

					//fmt.Println(mergeRequests)
				} else if sourceType == "project" {
					options := r.GenerateProjectMergeRequestsOptions()

					mergeRequests, _, err := git.MergeRequests.ListProjectMergeRequests(sourceID, options)
					if err != nil {
						return err
					}

					for _, mr := range mergeRequests {
						c, _ := r.CommentMR(mr)
						c, _ = c.CreateComment(git, dryrun)

					}

					logrus.WithField("count", len(mergeRequests)).Info("found merge requests")
				}
			*/
		}
	}

	return nil
}
