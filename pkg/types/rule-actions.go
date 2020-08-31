package types

import (
	"bufio"
	"bytes"
	"html/template"
	"reflect"

	"github.com/Masterminds/sprig"
	"github.com/leekchan/gtf"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// Summary --
type Summary struct {
	Options *gitlab.CreateIssueOptions
	Project string
}

// Summarize --
func (rule *Rule) Summarize(opts *Options, issues []*gitlab.Issue) error {
	if rule.Actions.Summarize == nil {
		return nil
	}

	logrus.WithField("component", "summarize").Info("building summary")

	var titleBuf bytes.Buffer
	titleOut := bufio.NewWriter(&titleBuf)

	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)

	titleTpl, err := template.
		New("template").
		Funcs(sprig.FuncMap()).
		Funcs(gtf.GtfFuncMap).
		Parse(rule.Actions.Summarize.Title)
	if err != nil {
		return err
	}

	tpl, err := template.New("template").Funcs(sprig.FuncMap()).Funcs(gtf.GtfFuncMap).Parse(rule.Actions.Summarize.Summary)
	if err != nil {
		return err
	}

	data := struct {
		Issues []*gitlab.Issue
	}{
		Issues: issues,
	}

	if err := tpl.Execute(out, data); err != nil {
		return err
	}
	if err := titleTpl.Execute(titleOut, data); err != nil {
		return err
	}

	if err := titleOut.Flush(); err != nil {
		return err
	}
	if err := out.Flush(); err != nil {
		return err
	}

	logrus.Debug(buf.String())

	summary := &Summary{
		Options: &gitlab.CreateIssueOptions{
			Title:       gitlab.String(titleBuf.String()),
			Description: gitlab.String(buf.String()),
		},
		Project: rule.Actions.Summarize.Destination,
	}

	if opts.dryRun {
		logrus.WithField("destination", summary.Project).Warn("WOULD have created a summary issue")
		return nil
	}

	if _, _, err := opts.client.Issues.CreateIssue(summary.Project, summary.Options); err != nil {
		return err
	}

	return nil
}

// Label --
func (rule *Rule) Label(opts *Options, entity interface{}) error {
	log := logrus.WithField("action", "label")

	if rule.Actions.Labels == nil {
		log.Debug("no label action defined")
		return nil
	}

	switch reflect.TypeOf(entity).String() {
	case "*gitlab.Issue":
		issue := entity.(*gitlab.Issue)

		log = log.
			WithField("author", issue.Author.Username).
			WithField("id", issue.ID).
			WithField("title", issue.Title).
			WithField("project", issue.ProjectID).
			WithField("created", issue.CreatedAt)

		rule.labelIssue(opts, issue, log)
		break
	case "*gitlab.MergeRequest":
		mergeRequest := entity.(*gitlab.MergeRequest)

		log = log.
			WithField("author", mergeRequest.Author.Username).
			WithField("id", mergeRequest.ID).
			WithField("title", mergeRequest.Title).
			WithField("project", mergeRequest.ProjectID).
			WithField("created", mergeRequest.CreatedAt)

		rule.labelMergeRequest(opts, mergeRequest, log)
		break
	}

	return nil
}

// Comment --
func (rule *Rule) Comment(opts *Options, entity interface{}) error {
	if rule.Actions.Comment == nil {
		logrus.Debug("no comment action defined")
		return nil
	}

	switch reflect.TypeOf(entity).String() {
	case "*gitlab.Issue":
		issue := entity.(*gitlab.Issue)

		log := logrus.
			WithField("author", issue.Author.Username).
			WithField("title", issue.Title).
			WithField("created", issue.CreatedAt)

		ok, err := rule.checkMemberPermissions(opts, issue.ProjectID, issue.Author.ID)
		if err != nil {
			return err
		}
		if !ok {
			log.Warn("user does not have permissions to modify labels, will not ask them to")
			return nil
		}

		// TODO: if user doesn't have permission, @ the right team members to come take a look

		rule.commentIssue(opts, issue, log)
		break
	case "*gitlab.MergeRequest":
		mergeRequest := entity.(*gitlab.MergeRequest)

		log := logrus.
			WithField("author", mergeRequest.Author.Username).
			WithField("id", mergeRequest.ID).
			WithField("title", mergeRequest.Title).
			WithField("project", mergeRequest.ProjectID)

		rule.commentMergeRequest(opts, mergeRequest, log)
		break
	}

	return nil
}

func (rule *Rule) checkMemberPermissions(opts *Options, projectID, memberID int) (bool, error) {
	if rule.Conditions.Author != nil && rule.Conditions.Author.CanLabel {
		member, _, err := opts.client.ProjectMembers.GetAllProjectMember(projectID, memberID)
		if err != nil {
			return false, err
		}

		if gitlab.ReporterPermissions <= member.AccessLevel {
			return true, nil
		}

		return false, nil
	}

	return true, nil
}

func (rule *Rule) labelIssue(opts *Options, issue *gitlab.Issue, log *logrus.Entry) error {
	log = log.WithField("component", "label-issue")

	addLabels := gitlab.Labels(rule.Actions.Labels)

	if opts.dryRun == true {
		log.WithField("addLabels", addLabels).Warn("WOULD have added label to issue")
		return nil
	}

	_, _, err := opts.client.Issues.UpdateIssue(issue.ProjectID, issue.ID, &gitlab.UpdateIssueOptions{
		AddLabels: &addLabels,
	})
	if err != nil {
		log.WithError(err).Error("unable to update issue")
		return err
	}

	return nil
}

func (rule *Rule) labelMergeRequest(opts *Options, mergeRequest *gitlab.MergeRequest, log *logrus.Entry) error {
	logrus.WithField("component", "label-merge-request").Warn("not implemented yet")
	return nil
}

func (rule *Rule) commentIssue(opts *Options, issue *gitlab.Issue, log *logrus.Entry) error {
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)

	tpl, err := template.New("comment_template").Funcs(sprig.FuncMap()).Funcs(gtf.GtfFuncMap).Parse(*rule.Actions.Comment)
	if err != nil {
		log.WithError(err).Error("unable to compile template")
		return err
	}

	data := struct {
		Issue *gitlab.Issue
	}{
		Issue: issue,
	}

	if err := tpl.Execute(out, data); err != nil {
		log.WithError(err).Error("unable to execute template")
		return err
	}

	if err := out.Flush(); err != nil {
		log.WithError(err).Error("unable to flush io writer")
		return err
	}

	if opts.dryRun == true {
		log.Warn("WOULD have created note on issue")
		return nil
	}

	note, _, err := opts.client.Notes.CreateIssueNote(issue.ProjectID, issue.ID, &gitlab.CreateIssueNoteOptions{
		Body: gitlab.String(buf.String()),
	})
	if err != nil {
		log.WithError(err).Error("unable to create issue note")
		return err
	}

	log.WithField("note", note.ID).WithField("issue", issue.IID).WithField("project", issue.ProjectID).Info("note created")

	return nil
}

func (rule *Rule) commentMergeRequest(opts *Options, issue *gitlab.MergeRequest, log *logrus.Entry) error {
	log.Warn("not implemented")
	return nil
}

// State will change the state of an issue
func (rule *Rule) State(opts *Options, issue *gitlab.Issue) error {
	log := logrus.
		WithField("component", "state").
		WithField("issue", issue.ID).
		WithField("project", issue.ProjectID)

	if rule.Actions.State == nil {
		log.WithField("skipped", true).Debug("skipped, no state action defined")
		return nil
	}

	compareState := "unknown"
	newState := "unknown"
	if *rule.Actions.State == "close" || *rule.Actions.State == "closed" {
		compareState = "closed"
		newState = "opened"
	} else if *rule.Actions.State == "open" || *rule.Actions.State == "opened" {
		compareState = "opened"
		newState = "closed"
	}

	if *rule.Actions.State != compareState {
		if opts.dryRun == true {
			log.WithField("old", issue.State).WithField("new", newState).Warn("WOULD have updated issue state")
			return nil
		}

		_, _, err := opts.client.Issues.UpdateIssue(issue.ProjectID, issue.ID, &gitlab.UpdateIssueOptions{
			StateEvent: gitlab.String(newState),
		})
		if err != nil {
			return nil
		}
	}

	return nil
}
