# ghreport

Application to check Github for Pull Requests, that are not Drafts, in repos the user cares about.

## Usage
The program is configured via a YAML file located at `$HOME/.config/ghreport/config.yaml`.

* token: Should be set to a Github API Token with access to the repos you are checking
    * Set permissions for token to repo - full control of private repositories, enable SSO if your repos require it
    * ![Github Personal Access Token Permissions](https://github.com/Jmainguy/ghreport/blob/main/docs/permissions.png?raw=true)
* One of subscribedRepos, autoDiscover.organizations, or autoDiscover.users must be set if you wish to have any results. You can set all three if you wish.
* topic limits what is returned from organizations and users to just that topic, this is an optional field.

Here's an example configuration:

```yaml
autoDiscover:
  organizations:
    - name: your_organization_name
      topic: topic_to_watch
  users:
    - name: your_username
      topic: topic_to_watch
subscribedRepos:
  - Jmainguy/ghreport
  - Jmainguy/bible
  - Jmainguy/ghReview
  - Jmainguy/bak
token: e0e9eac4e84446df6f3db180d07bfb222e91234
defaultOutput: singleline   # <--- NEW: set default output format here (default, singleline, json)
```

## Output Format Configuration

You can control the output format in two ways:

1. **Command-line flag**:  
   Use `-output=<format>` to override the output format for a single run.

2. **Config file option**:  
   Set `defaultOutput` in your `config.yaml` to specify the default output format for all runs.  
   If the `-output` flag is not provided, the value from `defaultOutput` will be used.

Supported formats are:
- `default` (multi-line, human-readable)
- `singleline` (one line per PR)
- `json` (one JSON object per PR)

**Precedence:**  
If both are set, the command-line flag takes priority over the config file.

Running the program

```
ghreport [-output=<format>]
```

### Output Options

The `-output` flag controls the format of the program's output. The following formats are supported:

- **default** (or omitted): Multi-line, human-readable output with fields on separate lines.
- **singleline**: All information for each pull request is printed on a single line.
- **json**: Each pull request is output as a single JSON object per line.

#### Example usage:

```
ghreport -output=default      # Multi-line output (default if omitted)
ghreport -output=singleline   # Single-line output
ghreport -output=json         # JSON output, one object per line
```

### Sample output with default format.

```
jmainguy@fedora:~/Github/ghreport$ ./ghreport 
https://github.com/Jmainguy/statuscode/pull/32
  author: renovate
  Age: 3 days 
  reviewDecision: ❌
  mergeable ✅
https://github.com/Jmainguy/statuscode/pull/33
  author: renovate
  Age: 3 days 
  reviewDecision: ✅
  mergeable ✅
https://github.com/Standouthost/Multicraft/pull/9
  author: TheWebGamer
  Age: 3321 days 
  reviewDecision: ✅
  mergeable ❌
https://github.com/Standouthost/Multicraft/pull/28
  author: ungarscool1
  Age: 2700 days 
  reviewDecision: 😅
  mergeable ✅
```

### Sample output with single-line format

```
jmainguy@fedora:~/Github/ghreport$ ./ghreport -output singleline
https://github.com/Jmainguy/coastie-operator/pull/2 author: dependabot Age: 409 days reviewDecision: 🔍 mergeable: ✅
https://github.com/Jmainguy/coastie-operator/pull/3 author: dependabot Age: 384 days reviewDecision: 🔍 mergeable: ✅
https://github.com/Jmainguy/coastie-operator/pull/4 author: renovate Age: 354 days reviewDecision: 🔍 mergeable: ✅
https://github.com/Standouthost/Multicraft/pull/9 author: TheWebGamer Age: 3321 days reviewDecision: ✅ mergeable: ❌
```


### Sample output with JSON format

```
jmainguy@fedora:~/Github/ghreport$ ./ghreport -output json
{"url":"https://github.com/Jmainguy/k8sCapcity/pull/50","author":"github-actions","age":"24 days","review_decision":"😅","mergeable":"✅"}
{"url":"https://github.com/Jmainguy/k8sDrainReport/pull/27","author":"github-actions","age":"8 days","review_decision":"😅","mergeable":"✅"}
{"url":"https://github.com/Standouthost/Multicraft/pull/9","author":"TheWebGamer","age":"3321 days","review_decision":"✅","mergeable":"❌"}
```

### Fun with jq and JSON
#### Select on repo name
```
jmainguy@fedora:~/Github/ghreport$ ./ghreport -output=json | jq '. | select(.url | contains("k8sCapcity"))'
{
  "url": "https://github.com/Jmainguy/k8sCapcity/pull/50",
  "author": "github-actions",
  "age": "24 days",
  "review_decision": "😅",
  "mergeable": "✅"
}
```
#### Select on author
```
jmainguy@fedora:~/Github/ghreport$ ./ghreport -output=json | jq '. | select(.author == "github-actions")'
{
  "url": "https://github.com/Jmainguy/k8sDrainReport/pull/27",
  "author": "github-actions",
  "age": "8 days",
  "review_decision": "😅",
  "mergeable": "✅"
}
{
  "url": "https://github.com/Jmainguy/k8sCapcity/pull/50",
  "author": "github-actions",
  "age": "24 days",
  "review_decision": "😅",
  "mergeable": "✅"
}
```

## Releases
We currently build releases for RPM, DEB, macOS, and Windows.

Grab Release from [The Releases Page](https://github.com/Jmainguy/ghreport/releases)

## Build it yourself
```/bin/bash
export GO111MODULE=on
go build
```
