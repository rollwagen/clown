package gitlab

import (
	"fmt"
	"os"
	"time"

	"github.com/rollwagen/clown/internal/prompter"
	clown "github.com/rollwagen/clown/pkg"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var owned bool

func NewCmd() *cobra.Command {
	return gitlabCmd
}

// GitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "CloneReposForGroup all projects (repos) in a group",
	Run: func(cmd *cobra.Command, args []string) {
		cloneGitlabGroup()
	},
}

func cloneGitlabGroup() {
	homeDir, _ := os.UserHomeDir()

	_ = godotenv.Load(homeDir + "/.clown")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	gitlabHost := os.Getenv("GITLAB_HOST")

	gitlabPlatform := clown.New(clown.Gitlab, gitlabHost, gitlabToken)
	groups := gitlabPlatform.ListGroups()
	p := prompter.New()
	idx, err := p.FuzzySelect("Select group to clone", "", groups)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Aborted. Exiting.")
		if err != nil {
			return
		}
		os.Exit(1)
	}

	progress := spinner.New(spinner.CharSets[43], 70*time.Millisecond)
	_ = progress.Color("cyan")
	cyan := color.New(color.FgCyan).SprintFunc()
	progress.Start()
	pollProgress := func(ch <-chan string) {
		for projectName := range ch {
			progress.Suffix = fmt.Sprintf(" Cloning project %s to folder %s/ ...", cyan(projectName), groups[idx])
		}
	}
	repoProgressChannel := make(chan string)
	go pollProgress(repoProgressChannel)
	_ = gitlabPlatform.CloneReposForGroup(groups[idx], repoProgressChannel)
	progress.FinalMSG = fmt.Sprintf("Finished cloning projects to folder %s.\n", groups[idx])
	progress.Stop()

	os.Exit(0)
}

func init() {
	gitlabCmd.Flags().BoolVarP(&owned, "owned", "o", false, "Limit to groups explicitly owned by the current user")
}
