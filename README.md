# ghreport

Application to check Github for Pull Requests, that are not Drafts, in repos the user cares about.

## Usage
The program takes no arguments, and is configured via ENV variables. 

* ghreportToken: Should be set to a Github API Token with access to the repos you are checking
* subscribedRepos: Should be set to a space delimmited list of Github Repos you want to check

Example configuration in ~/.bashrc
```
export ghreportToken=e0e9eac4e84446df6f3db180d07bfb222e91234
export subscribedRepos="Jmainguy/ghreport Jmainguy/bible Jmainguy/ghReview Jmainguy/bak"
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

## Releases
Grab Release from [The Releases Page](https://github.com/Jmainguy/ghreport/releases)

## Build
```/bin/bash
export GO111MODULE=on
go build
```
