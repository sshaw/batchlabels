# Batch Labels

Add or remove labels in batches to/from GitHub issues and pull requests.

## Installation

Download the binary for your platform:

* [Linux x86-64](https://github.com/sshaw/batchlabels/releases/download/v0.1.0/batchlabels-linux-x86-64)
* [Mac](https://github.com/sshaw/batchlabels/releases/download/v0.1.0/batchlabels-mac)
* [Windows](https://github.com/sshaw/batchlabels/releases/download/v0.1.0/batchlabels.exe)

Otherwise, [install Go](https://golang.org/dl/), and run `go get github.com/sshaw/batchlabels/...`
and put `$GOPATH/bin` (assuming `GOPATH` has one path) in your `PATH`.

An [Auth token from GitHub](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token) with repository access will need to be generated. To use, pass the token with the flag `-a` when using batchlables or store it in `BATCHLABELS_AUTH_TOKEN` environment variable.

## Usage

    batchlabels [hipv] [-a token] [--hacktoberfest] command label repo [repoN ...]
    Add or remove labels in batches to/from GitHub issues and pull requests.

    Options
    -a --auth token    repository auth token, defaults to the BATCHLABELS_AUTH_TOKEN environment var
    -h --help          print this message
    --hacktoberfest    add "hacktoberfest" labels to the given IDs (see label below) or, if none are given,
                       to all open issues (not pull requests) in the given repository
    -i --issues        Only apply labels to issues and not pull requests
    -p --pull-requests Only apply labels to pull requests
    -v --version       print the version

    command must be add or remove.

    label can be one of: label, label#color, issue:label#color or issue1,issue2:labelA#color,labelB#color
    When --hacktoberfest is provided and label is an integer or list of integers they are
    treated as issue IDs for which the hacktoberfest label will be applied.

    color is the hex color for the label.
    If label contains no issues it will be added or removed to/from every open issue in the repo(s).

    repo must be given in username/reponame format.

## [Hacktoberfest](https://hacktoberfest.digitalocean.com/)

Use the `--hacktoberfest` option to add an `hacktoberfest` label with the color `#ff9a56` to the specified issues
or, if none are given, all open issues (not pull requests).

This works with both `add` and `delete` commands and can also be used in conjunction with normal label arguments.
See [Examples](#Examples).

## Examples

Add the labels `foo` and `bar` to all open issues and pull requests in repository `sshaw/git-link`:

    batchlabels add foo bar sshaw/git-link

Remove the labels `foo` and `bar` from all open issues and pull requests in repository `sshaw/git-link`:

    batchlabels remove foo bar sshaw/git-link

Add the labels `foo` and `bar` only to open issues in repository `sshaw/git-link`:

    batchlabels -i add foo bar sshaw/git-link

To add labels having more than one word, use backslash `\` to escape spaces. Example below will add `good first issue` to all open issues in the repository `sshaw/git-link:

    batchlabels -i add good\ first\ issue sshaw/git-link

Add the label `foo` with the color `#de5833` only to open issues in repository `sshaw/git-link`:

    batchlabels -i add foo#de5833 sshaw/git-link

Add the label `foo` with the color `#de5833` to issue `44` in repository `sshaw/git-link`:

    batchlabels add 44:foo#de5833 sshaw/git-link

Add the label `foo` with color `#ccc` and label `bar` with color `#fff` to issues `12` and `27` in repository `sshaw/git-link`:

    batchlabels add 12,27:foo#ccc,bar:#fff sshaw/git-link

Remove the labels `foo` from issues `44` and `12` in repository `sshaw/git-link`:

    batchlabels remove 12,44:foo sshaw/git-link

Add the `hacktoberfest` label with color `#ff9a56` to all open issues in repository `sshaw/git-link`:

    batchlabels --hacktoberfest add sshaw/git-link

Add the `hacktoberfest` label with color `#ff9a56` to pull requests `44` and `12` in repository `sshaw/git-link`:

    batchlabels --hacktoberfest add 44,12 sshaw/git-link

Add the `hacktoberfest` label with color `#ff9a56` to all open issues in repository `sshaw/git-link`.
In addition, add the label `foo` with the color `#de5833` to issue `44`

    batchlabels --hacktoberfest add 44:foo#de5833 sshaw/git-link

## See Also

- [Export Pull Requests](https://github.com/sshaw/export-pull-requests) - Export pull requests and/or issues to a CSV file

## TODO

GitLab! (But first check to make sure they don't already support doing it.)

## Author

Skye Shaw [skye.shaw AT gmail.com]

## License

Released under the MIT License: http://www.opensource.org/licenses/MIT

---

Made by [ScreenStaring](http://screenstaring.com)
