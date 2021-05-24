package api

import (
	"fmt"
	"strconv"

	"github.com/dmurray-lacework/tally-git/internal/array"
)

type GithubService struct {
	client *Client
}

var TECH_ALLY_REPOS []string = []string{"go-sdk", "terraform-", "chef-lacework", "homebrew-tap", "circleci-orb"}

func (gs *GithubService) GetRepos() []Repo {
	response := &[]GithubRepositoryResponse{}
	apiPath := fmt.Sprintf(REPOS_URL)
	openSourceRepos := []Repo{}

	gs.client.RequestDecoder("GET", apiPath, nil, response)

	for _, r := range *response {
		repo := Repo{Name: r.Name, Issues: r.OpenIssues, Private: r.Private}
		if array.ContainsStr(TECH_ALLY_REPOS, repo.Name) && !repo.Private {
			// Get Pull Request details for repo
			repo.PullRequests = gs.GetPRCount(repo.Name)
			openSourceRepos = append(openSourceRepos, repo)
		}
	}

	return openSourceRepos
}

func (gs *GithubService) GetPRCount(repoName string) int {
	var openPRs int
	response := &[]PullRequestResponse{}
	apiPath := fmt.Sprintf(PR_URL, repoName)
	gs.client.RequestDecoder("GET", apiPath, nil, response)

	for _, r := range *response {
		if r.State == "open" {
			openPRs++
		}
	}

	return openPRs
}

type Repo struct {
	Name         string
	Issues       int
	PullRequests int
	Private      bool
}

type GithubRepositoryResponse struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	OpenIssues int    `json:"open_issues_count"`
	Private    bool   `json:"private"`
}

type PullRequestResponse struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *Repo) ToArray() []string {
	array := []string{r.Name, strconv.Itoa(r.Issues), strconv.Itoa(r.PullRequests)}
	return array
}
