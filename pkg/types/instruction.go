package types

// Rule --
type Rule struct {
	Name       string          `json:"name" yaml:"name"`
	Conditions *RuleConditions `json:"conditions" yaml:"conditions"`
	Actions    *RuleActions    `json:"actions" yaml:"actions"`
	Limits     *RuleLimits     `json:"limits,omitempty" yaml:"limits,omitempty"`

	// Allow for embedding additional rules, only applies for epics, milestone
	Epics         *ResourcePolicy `json:"epics,omitempty" yaml:"epics,omitempty"`
	Issues        *ResourcePolicy `json:"issues,omitempty" yaml:"issues,omitempty"`
	MergeRequests *ResourcePolicy `json:"merge_requests,omitempty" yaml:"merge_requests,omitempty"`

	Filters []*Filter `json:"filters,omitempty" yaml:"filters,omitempty"`
}

// RuleConditions --
type RuleConditions struct {
	State           string             `json:"state" yaml:"state"`
	Labels          []string           `json:"labels,omitempty" yaml:"labels,omitempty"`
	IgnoredLabels   []string           `json:"ignored_labels,omitempty" yaml:"ignored_labels,omitempty"`
	MissingLabels   []string           `json:"missing_labels,omitempty" yaml:"missing_labels,omitempty"`
	ForbiddenLabels []string           `json:"forbidden_labels,omitempty" yaml:"forbidden_labels,omitempty" comment:"Labels to which the action should not take place against"`
	AuthorMember    *AuthorMember      `json:"author_member,omitempty" yaml:"author_member,omitempty"`
	Author          *Author            `json:"author,omitempty" yaml:"author,omitempty"`
	Date            *RuleConditionDate `json:"date,omitempty" yaml:"date,omitempty"`
	Weight          string             `json:"weight,omitempty" yaml:"weight,omitempty"`
	Milestone       string             `json:"milestone,omitempty" yaml:"milestone,omitempty"`
	Expired         bool               `json:"expired,omitempty" yaml:"expired,omitempty" comment:"Only valid for filters for milestones"`

	// Epic Only
	FixedDates bool `json:"fixed_dates,omitempty" yaml:"fixed_dates,omitempty" resource:"epic" comment:"Only valid for filters for epics"`
}

// RuleConditionDate --
type RuleConditionDate struct {
	Condition string   `json:"condition" yaml:"condition"`
	Duration  Duration `json:"duration" yaml:"duration"`
}

// RuleActions --
type RuleActions struct {
	State       *string              `json:"state,omitempty" yaml:"state,omitempty"`
	Labels      []string             `json:"labels,omitempty" yaml:"labels,omitempty"`
	Unlabel     []string             `json:"unlabel,omitempty" yaml:"unlabel,omitempty"`
	Comment     *string              `json:"comment,omitempty" yaml:"comment,omitempty"`
	CommentType *string              `json:"comment_type,omitempty" yaml:"comment_type,omitempty"`
	Summarize   *RuleActionSummarize `json:"summarize,omitempty" yaml:"summarize,omitempty"`
}

// RuleLimits --
type RuleLimits struct {
	MostRecent *int `json:"most_recent,omitempty" yaml:"most_recent,omitempty"`
	PerPage    *int `json:"per_page,omitempty" yaml:"per_page,omitempty"`
}

// RuleActionSummarize --
type RuleActionSummarize struct {
	Title       string `json:"title" yaml:"title"`
	Item        string `json:"item" yaml:"item"`
	Summary     string `json:"summary" yaml:"summary"`
	Destination string `json:"destination" yaml:"destination"`
}

// AuthorMember --
type AuthorMember struct {
	Source    string `json:"source" yaml:"source"`
	Condition string `json:"condition" yaml:"condition"`
	SourceID  int    `json:"source_id" yaml:"source_id"`
}

// Author --
type Author struct {
	CanLabel bool   `json:"can_label" yaml:"can_label"`
	MemberOf string `json:"member_of,omitempty" yaml:"member_of,omitempty"`
}
