package types

import (
	"errors"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v1"
)

func TestPolicyParse1(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("fixtures/policy-parse1-1.yaml")
	if err != nil {
		panic(err)
	}

	var resource Policy

	err = yaml.Unmarshal(yamlFile, &resource)
	if err != nil {
		t.Error(err)
	}

	if resource.ResourceRules != nil {
		t.Error(errors.New("rules should be nil"))
	}
}

func TestPolicyParse2(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("fixtures/policy-parse1-2.yaml")
	if err != nil {
		panic(err)
	}

	var p Policy

	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		t.Error(err)
	}

	if p.ResourceRules == nil {
		t.Error(errors.New("rules should not be nil"))
	}

	if p.ResourceRules.Issues != nil {
		t.Error(errors.New("resource rules for issues should be nil"))
	}
}

func TestPolicyParse3(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("fixtures/policy-parse1-3.yaml")
	if err != nil {
		t.Error(err)
	}

	var p Policy

	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		t.Error(err)
	}

	if p.ResourceRules == nil {
		t.Error(errors.New("rules should not be nil"))
	}

	if p.ResourceRules.Issues == nil {
		t.Error(errors.New("resource rules for issues should not be nil"))
	}

	if len(p.ResourceRules.Issues.Rules) != 1 {
		t.Error(errors.New("there should be a single rule defined"))
	}

	rule := p.ResourceRules.Issues.Rules[0]

	expected := &gitlab.ListGroupIssuesOptions{
		State:  gitlab.String("opened"),
		Labels: gitlab.Labels{"none"},
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
		},
		OrderBy: gitlab.String("created_at"),
		Sort:    gitlab.String("desc"),
	}

	actual := rule.GenerateGroupIssueOptions()

	if !reflect.DeepEqual(actual, expected) {
		t.Log(expected)
		t.Log(actual)
		t.Error(errors.New("options are not equal"))
	}
}
