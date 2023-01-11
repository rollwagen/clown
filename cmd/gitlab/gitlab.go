package gitlab

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"

	"github.com/rollwagen/clown/pkg/config"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/rollwagen/clown/internal/prompter"
	"github.com/rollwagen/clown/pkg/git"
	"github.com/spf13/cobra"
)

var owned bool

func NewCmd() *cobra.Command {
	return gitlabCmd
}

// GitlabCmd represents the gitlab command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Clone all projects (repos) in a group",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		c := config.New()

		// user interaction i.e. let user select the group to clone
		gitlabPlatform := git.New(git.Gitlab, c.Host, c.AuthToken)
		groups := gitlabPlatform.ListGroups()
		log.Debug().Msgf("Group list length len(groups)=%d", len(groups))
		p := prompter.New()
		idx, err := p.Select("Select group to clone", "", groups)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Aborted. Exiting.")
			if err != nil {
				return
			}
			os.Exit(1)
		}
		log.Debug().Msgf("Selected group idx=%d groups[idx]=%s", idx, groups[idx])

		// actual cloning of all repos in a group
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
	},
}

func init() {
	gitlabCmd.Flags().BoolVarP(&owned, "owned", "o", false, "DEPRECATED - Limit to groups explicitly owned by the current user")
}
