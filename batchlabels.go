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

	commandAdd = "add"
	commandRemove = "remove"

	version = "v0.0.1"
	userAgent = "Batch Labels " + version
)

//    If label contains no issues it will be added to every open issue in its repo
// or remove
const usage = `batchlabels [-a auth] command label repo [commandN labelN repoN ...]
Add or remove labels in batches to/from GitHub issues and pull requests.

Options
-a --auth   repository auth token, defaults to the BATCHLABELS_AUTH_TOKEN environment var

command must be add or remove.

color is the hex color for the issue.
label can be one of: label, label#color, issue:label#color or issue1,issue2:label1#color,label2#color
repo must be given in username/reponame format.

If no issues are given the labels are added or removed to/from all open issues.
`

// Regexp to match repository: username/reponame
var repoRegexp = regexp.MustCompile(`^([^/]+)/([^/]+)$`)

// Issue Label
type Label struct {
	Color string
	Name string
}

// Repository Issue
type Issue struct {
	ID string
	Labels []Label
}

// Repository
type Repo struct {
	Owner string
	Name string
	Issues []Issue
}

// func (r Repo) String() string {
// 	var b strings.Builder
// 	b.WriteString(r.Owner)

// 	if r.Name != "" {
// 		b.WriteString("/")
// 		b.WriteString(r.Name)
// 	}

// 	return b.String()
// }


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
		for _, issue := range(repo.Issues) {
			if issue.ID != allIssues {
				addLabelToIssue(gh, repo, issue)
			} else {
				issues, err := ListOpenIssues(gh, repo)
				if err != nil {
					panic(err)
				}

				for _, i := range(issues) {
					issue.ID = strconv.Itoa(*i.Number)
					addLabelToIssue(gh, repo, issue)
				}
			}
		}
	}
}

func addLabelToIssue(gh *github.Client, repo Repo, issue Issue) {
	id, err := strconv.Atoi(issue.ID)
	if err != nil {
		panic(err)
	}

	var labels []string
	ctx := context.Background()

	for _, labelCfg := range(issue.Labels) {
		label := &github.Label{Color: &labelCfg.Color, Name: &labelCfg.Name}

		_, res, err := gh.Issues.CreateLabel(ctx, repo.Owner, repo.Name, label)
		// Assume 422 means it already exists
		if err != nil && res.StatusCode != 422 {
			panic(fmt.Errorf("Cannot create label for %s: %s\n", repo, err))
		}

		labels = append(labels, labelCfg.Name)
	}

	_, _, err = gh.Issues.AddLabelsToIssue(ctx, repo.Owner, repo.Name, id, labels)
	if err != nil {
		panic(fmt.Errorf("Cannot to add labels to %s: %s\n", repo, err))
	}
}

func RemoveLabels(gh *github.Client, repos []Repo) error {
	for _, repo := range(repos) {
		for _, issue := range(repo.Issues) {
			if issue.ID != allIssues {
				err := RemoveLabelsFromIssue(gh, repo, issue)
				if err != nil {
					return err
				}
			} else {
				issues, err := ListOpenIssues(gh, repo)
				if err != nil {
					return err
				}

				for _, i := range(issues) {
					issue.ID = strconv.Itoa(*i.Number)
					err := RemoveLabelsFromIssue(gh, repo, issue)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func RemoveLabelsFromIssue(gh *github.Client, repo Repo, issue Issue) error {
	errorFormat := "Cannot remove label for %s: %s"

	id, err := strconv.Atoi(issue.ID)
	if err != nil {
		return fmt.Errorf(errorFormat, repo, err)
	}

	ctx := context.Background()
	for _, label := range(issue.Labels) {
		res, err := gh.Issues.RemoveLabelForIssue(ctx, repo.Owner, repo.Name, id, label.Name)
		if err != nil && res.StatusCode != 404 {
			return fmt.Errorf(errorFormat, repo, err)
		}
	}

	return nil
}

func ListOpenIssues(gh *github.Client, repo Repo) ([]*github.Issue, error)  {
	ctx := context.Background()
	issues, _, err := gh.Issues.ListByRepo(ctx, repo.Owner, repo.Name, &github.IssueListByRepoOptions{State: "open"})

	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve open issues for %s: %s", repo, err)
	}

	return issues, nil
}
// parse user-supplied arguments and create Repo list containing issues and labels
func buildRepoList(argv []string) []Repo {
	var repos []Repo
	issues := make(map[string][]Label)

	i := 0
	for i < len(argv) {
		parts := repoRegexp.FindStringSubmatch(argv[i])

		// We have a repo
		if parts != nil {
			repo := Repo{Owner: parts[1], Name: parts[2]}
			for id, labels := range(issues) {
				repo.Issues = append(repo.Issues, Issue{ID: id, Labels: labels})
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
			issues[allIssues] = append(issues[allIssues], labels...)
		} else {
			ids := strings.Split(tags[0], argListSep)
			labels = createLabels(tags[1])

			for _, id := range(ids) {
				issues[id] = append(issues[id], labels...)
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
		label := Label{Name: nameColor[0]}
		if len(nameColor) > 1 {
			label.Color = nameColor[1]
		}

		labels = append(labels, label)
	}

	return labels
}

// ExitFailure print error to stderr and Exit() with code.
func ExitFailure(error string, code int)  {
	fmt.Fprintln(os.Stderr, error)
	os.Exit(code)
}

func main() {
	var auth string
	var help bool

	flag.Usage = func() {
		ExitFailure(usage, 2)
	}

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.StringVar(&auth, "a", os.Getenv("BATCHLABELS_AUTH_TOKEN"), "")
	flag.StringVar(&auth, "auth", os.Getenv("BATCHLABELS_AUTH_TOKEN"), "")
	flag.Parse()

	// At minimum we need: command label repo
	if help || len(flag.Args()) < 3 {
		flag.Usage();
	}

	command := flag.Arg(0)
	repos := buildRepoList(flag.Args()[1:])

	if len(repos) == 0 {
		ExitFailure("No repository given", 2)
	}

	gh := githubClient(auth)

	switch command {
	case commandAdd:
		addLabels(gh, repos)
	case commandRemove:
		fmt.Println("Removing labels...")
		err := RemoveLabels(gh, repos)
		if err != nil {
			ExitFailure("Failed to remove labels: " + err.Error(), 3)
		} else {
			fmt.Println("Labels successfully removed")
		}
	default:
		flag.Usage()
	}
}
