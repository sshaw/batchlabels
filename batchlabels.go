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

const usage = `batchlabels [h] [-a token] command label repo [commandN labelN repoN ...]
Add or remove labels in batches to/from GitHub issues and pull requests.

Options
-a --auth token     repository auth token, defaults to the BATCHLABELS_AUTH_TOKEN environment var
-h --help           print this message
--hacktoberfest     add "hacktoberfest" labels to all open issues in the given repo
-v --version        print the version

command must be add or remove.

label can be one of: label, label#color, issue:label#color or issue1,issue2:label1#color,label2#color
color is the hex color for the label.
If label contains no issues it will be added or removed to/from every open issue in its repo.

repo must be given in username/reponame format.

`

// Regexp to match repository: username/reponame
var repoRegexp = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
var hacktoberfestIssue = Issue{ ID: allIssues, Labels: []Label{ Label{Color: "ff9a56", Name:"hacktoberfest"} } }

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

func (r Repo) String() string {
	var b strings.Builder
	b.WriteString(r.Owner)

	if r.Name != "" {
		b.WriteString("/")
		b.WriteString(r.Name)
	}

	return b.String()
}

func GitHubClient(auth string) *github.Client {
	gh := github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth})))

	gh.UserAgent = fmt.Sprintf("%s - (%s)", userAgent, gh.UserAgent)
	return gh

}

func AddLabels(gh *github.Client, repos []Repo) error {
	for _, repo := range(repos) {
		for _, issue := range(repo.Issues) {
			if issue.ID != allIssues {
				err := AddLabelsToIssue(gh, repo, issue)
				if err != nil {
					return err
				}
			} else {
				issues, err := ListOpenIssues(gh, repo)
				if err != nil {
					return fmt.Errorf("Cannot find open issues for %s: %s", repo, err)
				}

				for _, i := range(issues) {
					issue.ID = strconv.Itoa(*i.Number)
					err := AddLabelsToIssue(gh, repo, issue)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func AddLabelsToIssue(gh *github.Client, repo Repo, issue Issue) error {
	errorFormat := "Cannot add labels to %s: %s"

	id, err := strconv.Atoi(issue.ID)
	if err != nil {
		return fmt.Errorf(errorFormat, repo, err)
	}

	var labels []string
	ctx := context.Background()

	for _, labelCfg := range(issue.Labels) {
		label := &github.Label{Color: &labelCfg.Color, Name: &labelCfg.Name}

		_, res, err := gh.Issues.CreateLabel(ctx, repo.Owner, repo.Name, label)
		// Assume 422 means it already exists
		if err != nil && res.StatusCode != 422 {
			return fmt.Errorf("Cannot create label for %s: %s\n", repo, err)
		}

		labels = append(labels, labelCfg.Name)
	}

	_, _, err = gh.Issues.AddLabelsToIssue(ctx, repo.Owner, repo.Name, id, labels)
	if err != nil {
		return fmt.Errorf(errorFormat, repo, err)
	}

	return nil
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

// BuildRepoList parse user-supplied arguments and create Repo list containing issues and labels
func BuildRepoList(argv []string) []Repo {
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
			labels = CreateLabels(tags[0])
			issues[allIssues] = append(issues[allIssues], labels...)
		} else {
			ids := strings.Split(tags[0], argListSep)
			labels = CreateLabels(tags[1])

			for _, id := range(ids) {
				issues[id] = append(issues[id], labels...)
			}
		}

		i++
	}

	return repos
}

func CreateLabels(s string) []Label {
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
func ExitFailure(error string, code int) {
	fmt.Fprintln(os.Stderr, error)
	os.Exit(code)
}

func main() {
	var auth string
	var showHelp, showVersion, hacktoberfest bool

	flag.Usage = func() {
		ExitFailure(usage, 2)
	}

	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&hacktoberfest, "hacktoberfest", false, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
	flag.StringVar(&auth, "a", os.Getenv("BATCHLABELS_AUTH_TOKEN"), "")
	flag.StringVar(&auth, "auth", os.Getenv("BATCHLABELS_AUTH_TOKEN"), "")
	flag.Parse()

	if showVersion {
		fmt.Println(version);
		os.Exit(0)
	}

	argv := flag.Args()

	// We need: command label repo or, if hacktoberfest: command repo
	if showHelp || !hacktoberfest && len(argv) < 3 || hacktoberfest && len(argv) < 2 {
		flag.Usage();
	}

	command := argv[0]
	repos := BuildRepoList(argv[1:])

	if len(repos) == 0 {
		ExitFailure("No repository given", 2)
	}

	if hacktoberfest {
		for i := 0; i < len(repos); i++ {
			repos[i].Issues = append(repos[i].Issues, hacktoberfestIssue)
		}
	}

	gh := GitHubClient(auth)

	switch command {
	case commandAdd:
		fmt.Println("Adding labels...")
		err := AddLabels(gh, repos)
		if err != nil {
			ExitFailure("Failed to add labels: " + err.Error(), 3)
		} else {
			fmt.Println("Labels successfully added")
		}
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
