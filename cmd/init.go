package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/xanzy/go-gitlab"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init the client",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configPath := path.Join(home, ".config", "go-cli", "config.yaml")

		fmt.Print("Enter gitlab root url [e.g. https://gitlab.com]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		url, err := url.Parse(input[:len(input)-1])
		if err != nil {
			log.Fatal(err.Error())
		}

		// https://gitlab.com/-/profile/personal_access_tokens?scopes=api,write_repository
		fmt.Printf("Enter token gitlab [click here to create https://%s/-/profile/personal_access_tokens?scopes=api]: ", url.Hostname())
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		token := input[:len(input)-1]

		url.Path = "/api/v4"
		url.Scheme = "https"
		git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url.String()))
		_, response, _ := git.Groups.ListGroups(&gitlab.ListGroupsOptions{})
		if response.Status != "200 OK" {
			log.Fatalf("was testing if gitlab api can be called but received not ok while requesting groups")
		}
		fmt.Printf("Enter gitlab repository where renovate is configured [e.g. https://%s/renovate]: ", url.Hostname())
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		renovate := input[:len(input)-1]

		viper.Set("gitlab.token", token)
		viper.Set("gitlab.url", url.String())
		viper.Set("gitlab.renovate.repository", renovate)
		err = viper.WriteConfigAs(configPath)
		if err != nil {
			fmt.Printf(err.Error())
		}
	},
}
