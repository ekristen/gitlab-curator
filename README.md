# GitLab Curator

The purpose of the GitLab Curator is to allow the triage team, product team, etc to define simple rules via yaml to help triage issues, merge requests, and epics (future)

## Credit / Attribution

I borrowed this idea from GitLab. They have a number of triage repos and scripts that is all ruby and not really setup as a standalone tool and while I didn't borrow any code, except a bit of yaml, the idea was their's first.

## Usage

At the core this is just a CLI tool that takes a source type (group or project) and a source id (group id or project id) and executes a policy against it. The policy is made up of rules that generally involve getting a list of issues or merge requests that meet specific conditions then take actions. The actions range from creating a rollup/summary issue, to adding labels or commenting on the issue.

```bash
gitlab-curator --token $GITLAB_TOKEN --source-type group --source-id 1234567 --file policy-file.yaml --dry-run
```

**Note:** with dry-run the logs will output the actions the tool *WOULD* had taken should the dry-run not have been specified.

### Usage in CI

This was designed to be used by a scheduled system or Gitlab CI directly. Basically on whatever schedule you want specific policies to run you simple call the CLI with the correct CLI options and that's it.

```yaml
unlabelled-issues:
  stage: hygeine
  image: docker.io/ekristen/gitlab-curator:master
  script:
    - gitlab-curator --token $TRIAGE_GITLAB_TOKEN --source-type group --source-id 1234567890 --file unlabelled-issues-summary.yaml
  rules:
    - if: '$UNLABELLED_ISSUES == "true"'

```

### Help

```
NAME:
   gitlab-curator - gitlab-curator

USAGE:
   gitlab-curator [global options] command [command options] [arguments...]

VERSION:
   dirty

AUTHOR:
   Erik Kristensen <erik@erikkristensen.com>

COMMANDS:
   run      run a policy file
   version  print version
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

#### Run Help

```
NAME:
   gitlab-curator run - run a policy file

USAGE:
   gitlab-curator run [command options] [arguments...]

OPTIONS:
   --token value        GitLab Token [$TOKEN]
   --source-type value  Source Type [$SOURCE_TYPE]
   --source-id value    Source ID [$SOURCE_ID]
   --file value         File [$FILE]
   --dry-run            Dry Run (default: false) [$DRY_RUN]
   --log-level value    Log Level (default: "info") [$LOG_LEVEL]
   --help, -h           show help (default: false)
```

## How it Works

This tool takes a very simple input in YAML to define rules to be run against `issues` or `merge requests`. The rules have conditions and actions.

The following YAML will find the most recent 50 issues that are opened without any labels, however there are no actions defined for this setup.

```yaml
resource_rules:
  issues:
    rules:
      - name: Collate latest unlabelled issues
        conditions:
          state: opened
          labels:
            - none
        limits:
          most_recent: 50
```

We can modify the YAML to have a summary action

```yaml
resource_rules:
  issues:
    rules:
      - name: Collate latest unlabelled issues
        conditions:
          state: opened
          labels:
            - none
        limits:
          most_recent: 50
        actions:
          summarize:
            destination: ekristen/triage
            title: |
              {{ now | date "2006-01-02" }} Newly created unlabelled issues requiring initial triage
            summary: |
              Hi Triage Team,

              Here is a list of the latest issues without labels in the project.

              {{ range .Issues }}
              - [ ] {{ .References.Full }} {{ .Title }}
              {{- end }}
```

Now this will find the 50 most recent issues without labels and generate a summary issue that will get opened in the project `ekristen/triage`. The title and summary fields are both golang templates that have the entire list of issues available to them. In the summary field, the template ranges (for loops) over the issues and generates a checklist list of issues.

## Templating System

This tool makes use of golang's templating system. This allows for easy interpolation of structured data into a template. These templates can then be rendered for the body of an issue or comment.

### Writing a Template

Templates are passed a structure of data, but it varies depending on the target. For the summary action, the data passed is `Issues` so it will be available with `{{ .Issues }}` when building the template. For all other actions they are acted upon a single issue or single merge request, so their data is `Issue` or `MergeRequest` and available via `{{ .Issue }}` or `{{ .MergeRequest }}`.

## Actions

- Label (issues or merge requests)
- Comment (issues or merge requests)
- Summarize (issues)
- State (issues)

### Label

This action is pretty straight forward, this will add labels to whatever issues or merge requests meet the search criteria.

### Comment

This action will leave a comment on the issue or merge request that meet the search criteria.

### Summarize

This action will take a search of issues or merge requests and create a summary issue that can be assigned to someone else.

#### State

This action will either open or close issues that meet the search criteria

## Ideas

### Report / Summary Ideas

- No labels, generate report
- Issues missing kind, area labels -> comment to author to add (if they have permissions to do so)
- Issues with labels not in milestone -> generate report
- Issues with milestone not in epic -> generate report
- Merge Requests without milestone -> comment to author to add (if they have permissions to do so)
