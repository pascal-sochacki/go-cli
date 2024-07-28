package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	renovateCmd.AddCommand(renovateThisCmd)
}

var renovateThisCmd = &cobra.Command{
	Use:   "this",
	Short: "Trigger renovate gitlab pipeline with Auto discvoer filter based on your current working direcory",
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		repo, err := git.PlainOpen(pwd)
		if err != nil {
			log.Fatalf(err.Error())
		}
		remotes, err := repo.Remotes()
		for _, v := range remotes {
			urls := v.Config().URLs
			for _, j := range urls {
				splits := strings.Split(j, ":")
				lastSegment := splits[len(splits)-1]
				println(lastSegment)

			}
		}
	},
}
