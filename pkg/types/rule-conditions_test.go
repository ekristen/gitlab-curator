package types

import (
	"errors"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/xanzy/go-gitlab"
)

func TestCondition1(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("fixtures/condition-created-older.yaml")
	if err != nil {
		panic(err)
	}

	var p Policy

	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		t.Error(err)
	}

	rule := p.ResourceRules.Issues.Rules[0]

	dur := 8760 * time.Hour
	now := time.Now().UTC()
	bef := now.Add(-dur)

	expected := &gitlab.ListGroupIssuesOptions{
		State:         gitlab.String("opened"),
		NotLabels:     gitlab.Labels{"automation/close"},
		CreatedBefore: &bef,
	}

	actual := rule.GenerateGroupIssueOptions()
	expected.CreatedBefore = actual.CreatedBefore

	if !reflect.DeepEqual(actual, expected) {
		t.Log(expected)
		t.Log(actual)
		t.Error(errors.New("options are not equal"))
	}
}
