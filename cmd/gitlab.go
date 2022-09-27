package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var owned bool

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Clone all projects (repos) in a group",
	Run: func(cmd *cobra.Command, args []string) {
		cloneGitlabGroup()
	},
}

type GitlabGroup struct {
	Name string
	ID   int
}

func cloneGitlabGroup() {
	home, _ := os.UserHomeDir()
	_ = godotenv.Load(home + "/.clown")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	gitlabHost := os.Getenv("GITLAB_HOST")

	if gitlabToken == "" {
		_, _ = fmt.Fprintln(os.Stderr, `Please define token with env variable GITLAB_TOKEN=glpat-... and gitlab url with
GITLAB_URL=https://gitlab.company.com, or configure in ~/.clown`)
		os.Exit(1)
	}

	baseURL := "https://" + gitlabHost + "/api/v4"
	gitlabClient, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(baseURL))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating gitlab client: %s\n", err)
		os.Exit(1)
	}

	var groups []GitlabGroup

	boolTrue := bool(true)
	availableGroups, _, _ := gitlabClient.Groups.ListGroups(
		&gitlab.ListGroupsOptions{
			AllAvailable: &boolTrue,
			Owned:        &owned,
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
		},
	)
	for _, g := range availableGroups {
		groups = append(groups, GitlabGroup{g.Name, g.ID})
	}

	sort.Slice(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].Name) > strings.ToLower(groups[j].Name)
	})

	idx, err := fuzzyfinder.Find(groups, func(i int) string {
		return groups[i].Name
	})
	if err != nil {
		os.Exit(1)
	}
	groupToClone := groups[idx]

	err = os.Mkdir(groupToClone.Name, 0o750)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating group directory '%s': %s", groupToClone.Name, err)
		os.Exit(1)
	}

	s := spinner.New(spinner.CharSets[43], 60*time.Millisecond)
	_ = s.Color("cyan")
	cyan := color.New(color.FgCyan).SprintFunc()
	s.Start()

	projects, _, _ := gitlabClient.Groups.ListGroupProjects(groupToClone.ID, &gitlab.ListGroupProjectsOptions{})

	for _, p := range projects {
		s.Suffix = fmt.Sprintf(" Cloning project %s to folder %s/ ...", cyan(p.Name), groupToClone.Name)
		path := groupToClone.Name + string(os.PathSeparator) + p.Name
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: p.SSHURLToRepo,
		})
		if err != nil {
			s.Stop()
			_, _ = fmt.Fprintf(os.Stderr, "Error cloning project %s:\n%s\nAborting.\n", cyan(p.Name), err)
		}
	}
	s.FinalMSG = fmt.Sprintf("Finished cloning %d projects to folder %s.\n", len(projects), groupToClone.Name)
	s.Stop()
}

func init() {
	rootCmd.AddCommand(gitlabCmd)
	gitlabCmd.Flags().BoolVarP(&owned, "owned", "o", false, "Limit to groups explicitly owned by the current user")
}
