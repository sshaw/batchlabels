# Batch Labels

Add or remove labels in batches to/from GitHub issues and pull requests.

## Installation

Download the binary for your platform:

* [Linux x86-64](https://github.com/sshaw/batchlabels/releases/download/v0.0.1/batchlabels-linux-x86-64)
* [Mac](https://github.com/sshaw/batchlabels/releases/download/v0.0.1/batchlabels-mac)
* Windows TODO

Otherwise, [install Go](https://golang.org/dl/), and run `go get github.com/sshaw/batchlabels/...`
And put `$GOPATH/bin` (assuming `GOPATH` has one path) in your `PATH`.


## Usage

    batchlabels [h] [-a token] command label repo [commandN labelN repoN ...]
    Add or remove labels in batches to/from GitHub issues and pull requests.

    Options
    -a --auth token  repository auth token, defaults to the BATCHLABELS_AUTH_TOKEN environment var
    -h --help        print this message
    --hacktoberfest  add "hacktoberfest" labels to all open issues in the given repo
    -v --version     print the version

    command must be add or remove

    label can be one of: label, label#color, issue:label#color or issue1,issue2:label1#color,label2#color
    color is the hex color for the label.
    If label contains no issues it will be added or removed to/from every open issue in its repo.

    repo must be given in username/reponame format.

## [Hacktoberfest](https://hacktoberfest.digitalocean.com/)

Use the `--hacktoberfest` option to add a `#ff9a56` `hacktoberfest` label to all your open issues and pull requests.

This works with both `add` and `delete` commands and can be used in conjunction with normal label arguments.

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
