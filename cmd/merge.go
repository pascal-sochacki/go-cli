package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func init() {
	rootCmd.AddCommand(mergeCmd)
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge all merge request with the label merge",
	Run: func(cmd *cobra.Command, args []string) {

		url := viper.GetString("gitlab.url")
		token := viper.GetString("gitlab.token")
		if url == "" || token == "" {
			log.Fatal("please run init")
		}
		git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		state := "opened"
		requests, _, err := git.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
			State:  &state,
			Labels: &gitlab.LabelOptions{"merge"},
		})
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		for _, v := range requests {
			_, _, err := git.MergeRequests.AcceptMergeRequest(v.ProjectID, v.IID, &gitlab.AcceptMergeRequestOptions{})
			if err != nil {
				log.Fatalf("Failed to merge %v", err)
				continue
			}
			log.Printf("Merged titel:%s in project id:%d\n", v.Title, v.ProjectID)
		}
	},
}
