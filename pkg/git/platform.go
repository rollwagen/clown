package git

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"

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

// Platform defines all common methods for a Git platform such as GitHub or Bitbucket.
type Platform interface {
	IsAuthenticated() bool
	CloneReposForGroup(groupName string, repoProgress chan<- string) error
	ListGroups() []string
}

func New(t Type, hostname, authToken string) Platform {
	if t == Gitlab {
		p := newGitlabPlatform(hostname, authToken)

		return &p
	}

	return nil
}

type platform struct {
	hostName  string
	authToken string
}

// Gitlab Platform implementation.
type gitlab struct {
	platform
	client *gogitlab.Client
}

func (g *gitlab) IsAuthenticated() bool {
	t, _, _ := g.client.PersonalAccessTokens.ListPersonalAccessTokens(&gogitlab.ListPersonalAccessTokensOptions{})

	return len(t) > 0
}

func (g *gitlab) CloneReposForGroup(groupName string, repoProgress chan<- string) error {
	err := os.Mkdir(groupName, 0o750)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating group directory '%s': %s", groupName, err)
		os.Exit(1)
	}

	log.Debug().Msgf("ListGroups with groupNmae='%s'", groupName)
	groups, response, _ := g.client.Groups.ListGroups(
		&gogitlab.ListGroupsOptions{
			AllAvailable: truePointer(),
			Search:       &groupName,
		})

	if response.TotalItems == 0 {
		msg := fmt.Sprintf("Expected a search result for '%s', but got none\n", groupName)
		log.Error().Msg(msg)
		_, _ = fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}

	if len(groups) > 1 {
		for i, g := range groups {
			log.Debug().Msgf("Found group[%d]=%s", i, g.Name)
		}
	}

	groupID := groups[0].ID // todo: validate if name matches for idx=0

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

func (g *gitlab) ListGroups() []string {
	groups := make([]string, 0, 100)

	availableGroups, response, _ := g.client.Groups.ListGroups(
		&gogitlab.ListGroupsOptions{
			AllAvailable: truePointer(),
			ListOptions: gogitlab.ListOptions{
				PerPage: 100,
			},
		},
	)
	log.Debug().Msgf("Response items=%d pages=%d", response.TotalItems, response.TotalPages)
	for _, group := range availableGroups {
		groups = append(groups, group.Name)
	}

	sort.Slice(groups, func(i, j int) bool {
		return strings.ToLower(groups[i]) < strings.ToLower(groups[j])
	})

	log.Debug().Msgf("Processed %d groups from gitlab", len(groups))

	return groups
}

func newGitlabPlatform(hostName, authToken string) gitlab {
	if authToken == "" || hostName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please define token with env variable GITLAB_TOKEN=glpat-... "+
			"and gitlab url with GITLAB_URL=gitlab.company.com, or configure in ~/.clown")

		os.Exit(1)
	}

	baseURL := "https://" + hostName + "/api/v4"

	client, err := gogitlab.NewClient(authToken, gogitlab.WithBaseURL(baseURL))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating gitlab client: %s\n", err)
		os.Exit(1)
	}

	return gitlab{
		platform: platform{
			hostName:  hostName,
			authToken: authToken,
		},
		client: client,
	}
}

func truePointer() *bool {
	boolTrue := true
	return &boolTrue
}
