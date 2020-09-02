module github.com/ekristen/gitlab-curator

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.1.1 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/leekchan/gtf v0.0.0-20190214083521-5fba33c5b00b
	github.com/mitchellh/copystructure v1.0.0
	github.com/sirupsen/logrus v1.6.0
	github.com/urfave/cli/v2 v2.2.0
	github.com/xanzy/go-gitlab v0.36.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
)

replace github.com/xanzy/go-gitlab => github.com/ekristen/go-gitlab v0.36.1-0.20200902004707-2edc9371c745
