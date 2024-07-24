package cmd

import (
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func init() {
	rootCmd.AddCommand(renovateCmd)
}

var renovateCmd = &cobra.Command{
	Use:   "renovate",
	Short: "Trigger renovate gitlab pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		gitlabUrl := viper.GetString("gitlab.url")
		token := viper.GetString("gitlab.token")
		if gitlabUrl == "" || token == "" {
			log.Fatal("please run init")
		}
		renovateRepo := viper.GetString("gitlab.renovate.repository")
		if renovateRepo == "" {
			log.Fatal("please run init and cofigure the renovate repository")
		}
		git, err := gitlab.NewClient(token, gitlab.WithBaseURL(gitlabUrl))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		renovateUrl, err := url.Parse(renovateRepo)
		if err != nil {
			log.Fatalf("Failed to parse renovate url: %s", renovateUrl)
		}
		project, _, err := git.Projects.GetProject(renovateUrl.Path[1:], &gitlab.GetProjectOptions{})
		if err != nil {
			log.Fatalf("Failed to look up renovate repository")
		}
		ref := "main"
		renovate_autodiscover_filter := "RENOVATE_AUTODISCOVER_FILTER"
		value := "/gibt/es/nicht"
		variables := []*gitlab.PipelineVariableOptions{}
		variables = append(variables, &gitlab.PipelineVariableOptions{
			Key:   &renovate_autodiscover_filter,
			Value: &value,
		})
		pipeline, _, err := git.Pipelines.CreatePipeline(project.ID, &gitlab.CreatePipelineOptions{
			Ref:       &ref,
			Variables: &variables,
		})
		if err != nil {
			log.Fatalf("Failed to trigger pipeline: %s", err.Error())
		}

		fmt.Println(pipeline.WebURL)
	},
}
