package main

import (
	"log"
	"os"

	"github.com/dmurray-lacework/tally-git/api"
	"github.com/dmurray-lacework/tally-git/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

var configuration config.Configuration

func main() {
	client, err := api.NewClient(configuration.Github.Apikey)
	if err != nil {
		os.Exit(1)
	}

	repos := client.Github.GetRepos()
	outputTable([]string{"Repository", "Issues", "Pull Requests"}, tableFormat(repos))
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func tableFormat(repos []api.RepoResponse) [][]string {
	data := [][]string{}
	for _, r := range repos {
		data = append(data, r.ToArray())
	}
	return data
}

func outputTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
