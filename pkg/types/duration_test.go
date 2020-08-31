package types

import (
	"fmt"
	"testing"

	"github.com/ghodss/yaml"
)

type DurationTest struct {
	Period Duration `json:"period" yaml:"period"`
}

func TestDurationStruct(t *testing.T) {
	y := []byte(`period: "5h"`)

	var p DurationTest

	if err := yaml.Unmarshal(y, &p); err != nil {
		t.Error(err)
	}

	if p.Period.String() != "5h0m0s" {
		t.Error(fmt.Errorf("duration does not match expected (actual: %s)", p.Period.String()))
	}
}
