# ghreport

Application to check Github for Pull Requests, that are not Drafts, in repos the user cares about.

## Usage
The program takes no arguments, and is configured via ENV variables. 

* ghreportToken: Should be set to a Github API Token with access to the repos you are checking
    * Set permissions for token to repo - full control of private repositories, enable SSO if your repos require it
    * ![Github Personal Access Token Permissions](https://github.com/Jmainguy/ghreport/blob/main/docs/permissions.png?raw=true)
* subscribedRepos: Should be set to a space delimmited list of Github Repos you want to check

Example configuration in ~/.bashrc
```
export ghreportToken=e0e9eac4e84446df6f3db180d07bfb222e91234
export subscribedRepos="Jmainguy/ghreport Jmainguy/bible Jmainguy/ghReview Jmainguy/bak"
```

Additionally if you have a long list of repos to watch you can use this format when setting the environment variable:
```
export subscribedRepos="\
somesite/aebot \
somesite/ansible-okta-aws-auth \
somesite/blahblah"
```

Running the progam
```
ghreport
```

Sample output

```
[jmainguy@jmainguy-7410 ghreport]$ ghreport 
https://github.com/Jmainguy/bible/pull/1
https://github.com/Jmainguy/bak/pull/123
```

## Linux / macOS homebrew install

```/bin/bash
brew install jmainguy/tap/ghreport
```

## Releases
We currently build releases for RPM, DEB, macOS, and Windows.

Grab Release from [The Releases Page](https://github.com/Jmainguy/ghreport/releases)

## Build it yourself
```/bin/bash
export GO111MODULE=on
go build
```
