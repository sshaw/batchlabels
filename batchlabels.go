package main

import (
	"fmt"
	"flag"
	"os"
	"regexp"
	"strings"
	"strconv"
	"context"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

const (
	argListSep = ","
	argIdLabelSep = ":"
	argColorSep = "#"
	allIssues = "__ALL__"

	version = "v0.0.1"
	userAgent = "Batch Labels " + version
)

//    If label contains no issues it will be added to every open issue in its repo
// or remove
const usage = `batchlabels [-a auth] command label repo [commandN labelN repoN ...]
Add/remove labels in batches to/from GitHub issues and pull requests.

Options
-a --auth   repository auth token, defaults to the BATCHLABELS_AUTH_TOKEN environment var

command must be add

color is the hex color for the issue
label can be one of: label, label#color, issue:label#color or issue1,issue2:label1#color,label2#color
repo must be given in username/reponame format
`

// Regexp to match repository: username/reponame
var repoRegexp = regexp.MustCompile(`^([^/]+)/([^/]+)$`)

type Label struct {
	color string
	name string
}

type Issue struct {
	id string
	labels []Label
}

type Repo struct {
	owner string
	name string
	issues []Issue
}

func (r Repo) String() string {
	var b strings.Builder
	b.WriteString(r.owner)

	if r.name != "" {
		b.WriteString("/")
		b.WriteString(r.name)
	}

	return b.String()
}


func githubClient(auth string) *github.Client {
	gh := github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth})))

	gh.UserAgent = fmt.Sprintf("%s - (%s)", userAgent, gh.UserAgent)
	return gh

}

func addLabels(gh *github.Client, repos []Repo)  {
	for _, repo := range(repos) {
		for _, issue := range(repo.issues) {
			if issue.id != allIssues {
				addLabelToIssue(gh, repo, issue)
			} else {
				// Not Yet!
				// var ids []string
				// for id := range(ids) {
				// 	issue.id = id // string
				// 	addLabelToIssue(issue)
				// }
			}
		}
	}
}

func addLabelToIssue(gh *github.Client, repo Repo, issue Issue) {
	id, err := strconv.Atoi(issue.id)
	if err != nil {
		panic(err)
	}

	var labels []string
	ctx := context.Background()

	for _, labelCfg := range(issue.labels) {
		label := &github.Label{Color: &labelCfg.color, Name: &labelCfg.name}

		_, res, err := gh.Issues.CreateLabel(ctx, repo.owner, repo.name, label)
		// Assume 422 means it already exists
		if err != nil && res.StatusCode != 422 {
			panic(fmt.Errorf("Cannot create label for %s: %s\n", repo, err))
		}

		labels = append(labels, labelCfg.name)
	}

	_, _, err = gh.Issues.AddLabelsToIssue(ctx, repo.owner, repo.name, id, labels)
	if err != nil {
		panic(fmt.Errorf("Cannot to add labels to %s: %s\n", repo, err))
	}
}

func removeLabels(gh *github.Client, repos []Repo)  {
	fmt.Printf("Not removing labels: %+v\n", repos)
}

// parse user-supplied arguments and create Repo list
func buildRepoList(argv []string) []Repo {
	var repos []Repo
	issues := make(map[string][]Label)

	i := 0
	for i < len(argv) {
		parts := repoRegexp.FindStringSubmatch(argv[i])

		// We have a repo
		if parts != nil {
			repo := Repo{owner: parts[1], name: parts[2]}
			for id, labels := range(issues) {
				repo.issues = append(repo.issues, Issue{id: id, labels: labels})
			}

			repos = append(repos, repo)

			i++
			continue
		}


		var labels []Label
		// Parse issue ids and labels
		tags := strings.SplitN(argv[i], argIdLabelSep, 2)

		if len(tags) == 1 {
			labels = createLabels(tags[0])
			issues[allIssues] = labels
		} else {
			ids := strings.Split(tags[0], argListSep)
			labels = createLabels(tags[1])

			for _, id := range(ids) {
				issues[id] = labels
			}
		}

		i++
	}

	return repos
}

func createLabels(s string) []Label {
	var labels []Label
	labelCfg := strings.Split(s, argListSep)

	for _, cfg := range(labelCfg) {
		nameColor := strings.SplitN(cfg, argColorSep, 2)
		label := Label{name: nameColor[0]}
		if len(nameColor) > 1 {
			label.color = nameColor[1]
		}

		labels = append(labels, label)
	}

	return labels
}

func main() {
	var auth string

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	flag.StringVar(&auth, "a", os.Getenv("BATCHLABELS_AUTH_TOKEN"), "")
	flag.StringVar(&auth, "auth", os.Getenv("BATCHLABELS_AUTH_TOKEN"), "")
	flag.Parse()

	// At minimum we need: command label repo
	if len(flag.Args()) < 3 {
		flag.Usage();
	}


	command := flag.Arg(0)
	repos := buildRepoList(flag.Args()[1:])

	if len(repos) == 0 {
		fmt.Fprintln(os.Stderr, "No repository given")
		os.Exit(2)
	}


	gh := githubClient(auth)

	switch command {
	case "add":
		addLabels(gh, repos)
	// case "remove":
	// 	removeLabels(gh, repos)
	default:
		flag.Usage()
	}
}
