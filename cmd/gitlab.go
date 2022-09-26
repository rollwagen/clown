/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"os"
	"sort"
	"strings"
)

// gitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains ...`,
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

	listAll := bool(true)
	availableGroups, _, _ := gitlabClient.Groups.ListGroups(
		&gitlab.ListGroupsOptions{
			AllAvailable: &listAll,
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
	groupToClone := groups[idx]

	err = os.Mkdir(groupToClone.Name, 0750)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating group directory '%s': %s", groupToClone.Name, err)
		os.Exit(1)
	}
	projects, _, _ := gitlabClient.Groups.ListGroupProjects(groupToClone.ID, &gitlab.ListGroupProjectsOptions{})
	for _, p := range projects {
		repoURL := "git@" + gitlabHost + ":" + groupToClone.Name + "/" + p.Name + ".git"
		fmt.Println(repoURL)
		//_, err = git.PlainClone(groupToClone.Name, false, &git.CloneOptions{
		path := groupToClone.Name + string(os.PathSeparator) + p.Name
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error cloning project %s: %s", p.Name, err)
		}
	}
}

func init() {
	rootCmd.AddCommand(gitlabCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitlabCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitlabCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}