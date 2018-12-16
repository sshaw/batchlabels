# Batch Labels

Add or remove labels in batches to/from GitHub issues and pull requests.

## Installation

...

## Usage

    batchlabels [-a auth] command label repo [commandN labelN repoN ...]
    Add or remove labels in batches to/from GitHub issues and pull requests.

    Options
    -a --auth   repository auth token, defaults to the BATCHLABELS_AUTH_TOKEN environment var

    command must be add or remove

    color is the hex color for the issue.
    label can be one of: label, label#color, issue:label#color or issue1,issue2:label1#color,label2#color
    repo must be given in username/reponame format.

    If no issues are given the labels are added or removed to/from all open issues.

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
