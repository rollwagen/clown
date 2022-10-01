package platform

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	gogitlab "github.com/xanzy/go-gitlab"
)

type Type int

const (
	Gitlab    Type = iota // Gitlab =0
	Github                // Github = 1
	Bitbucket             // Bitbucket = 2
)

func (t Type) String() string {
	return []string{"Gitlab", "Github", "Bitbucket"}[t]
}

// GitPlatform defines all common methods for a Git platform such as GitHub or Bitbucket.
type GitPlatform interface {
	IsAuthenticated() bool
	CloneReposForGroup(groupName string, repoProgress chan<- string) error
	ListGroups() []string
}

func New(t Type, hostname, authToken string) GitPlatform {
	if t == Gitlab {
		p := newGitlabPlatform(hostname, authToken)

		return &p
	}

	return nil
}

type gitPlatform struct {
	hostName  string
	authToken string
}

// Gitlab Platform implementation.
type gitlabPlatform struct {
	gitPlatform
	client *gogitlab.Client
}

func (g *gitlabPlatform) IsAuthenticated() bool {
	t, _, _ := g.client.PersonalAccessTokens.ListPersonalAccessTokens(&gogitlab.ListPersonalAccessTokensOptions{})

	return len(t) > 0
}

func (g *gitlabPlatform) CloneReposForGroup(groupName string, repoProgress chan<- string) error {
	err := os.Mkdir(groupName, 0o750)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating group directory '%s': %s", groupName, err)
		os.Exit(1)
	}

	groups, _, _ := g.client.Groups.ListGroups(&gogitlab.ListGroupsOptions{Search: &groupName})
	if len(groups) > 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Expected one search result for %s, but got %d", groupName, len(groups))
	}

	groupID := groups[0].ID

	projects, _, _ := g.client.Groups.ListGroupProjects(groupID, &gogitlab.ListGroupProjectsOptions{})

	for _, project := range projects {
		repoProgress <- project.Name
		path := groupName + string(os.PathSeparator) + project.Name
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: project.SSHURLToRepo,
		})

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error cloning project %s:\n%s\nAborting.\n", project.Name, err)

			return err
		}
	}

	return nil
}

func (g *gitlabPlatform) ListGroups() []string {
	groups := []string{}

	boolTrue := bool(true)

	availableGroups, _, _ := g.client.Groups.ListGroups(
		&gogitlab.ListGroupsOptions{
			AllAvailable: &boolTrue,
			ListOptions: gogitlab.ListOptions{
				PerPage: 100,
			},
		},
	)
	for _, group := range availableGroups {
		groups = append(groups, group.Name)
	}

	sort.Slice(groups, func(i, j int) bool {
		return strings.ToLower(groups[i]) < strings.ToLower(groups[j])
	})

	return groups
}

func newGitlabPlatform(hostName, authToken string) gitlabPlatform {
	if authToken == "" || hostName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please define token with env variable GITLAB_TOKEN=glpat-... "+
			"and gitlabPlatform url with GITLAB_URL=https://gitlabPlatform.company.com, or configure in ~/.clown")

		os.Exit(1)
	}

	baseURL := "https://" + hostName + "/api/v4"

	client, err := gogitlab.NewClient(authToken, gogitlab.WithBaseURL(baseURL))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating gitlabPlatform client: %s\n", err)
		os.Exit(1)
	}

	return gitlabPlatform{
		gitPlatform: gitPlatform{
			hostName:  hostName,
			authToken: authToken,
		},
		client: client,
	}
}
