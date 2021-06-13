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

func (gs *GithubService) GetRepos() []RepoResponse {
	response := &[]GithubRepositoryResponse{}
	apiPath := fmt.Sprintf(REPOS_URL)
	openSourceRepos := []RepoResponse{}

	gs.client.RequestDecoder("GET", apiPath, nil, response)

	for _, r := range *response {
		repo := RepoResponse{Name: r.Name,
			Issues:      r.OpenIssues,
			Private:     r.Private,
			Language:    r.Language,
			Description: r.Description}
		if array.ContainsStr(TECH_ALLY_REPOS, repo.Name) && !repo.Private {
			// Get Pull Request details for repo
			repo.PullRequests = gs.GetPRCount(repo.Name)
			repo.IssueDetails = gs.GetIssues(repo.Name)
			repo.PullRequestDetails = gs.GetPullRequests(repo.Name)
			openSourceRepos = append(openSourceRepos, repo)
		}
	}
	return openSourceRepos
}

func (gs *GithubService) GetPullRequests(repoName string) []PullRequestResponse {
	response := &[]PullRequestResponse{}
	apiPath := fmt.Sprintf(PR_URL, repoName)
	gs.client.RequestDecoder("GET", apiPath, nil, response)
	return *response
}

func (gs *GithubService) GetIssues(repoName string) []IssueResponse {
	response := &[]IssueResponse{}
	issues := &[]IssueResponse{}
	apiPath := fmt.Sprintf(ISSUES_URL, repoName)
	gs.client.RequestDecoder("GET", apiPath, nil, response)
	for _, issue := range *response {
		if issue.PRLinks == (PullLinksResponse{}) {
			*issues = append(*issues, issue)
		}
	}
	return *issues
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

type RepoResponse struct {
	Name               string                `json:"name"`
	Issues             int                   `json:"issues"`
	PullRequests       int                   `json:"pull_requests"`
	Private            bool                  `json:"private"`
	Language           string                `json:"language"`
	Url                string                `json:"url"`
	Description        string                `json:"description"`
	IssueDetails       []IssueResponse       `json:"issue_details"`
	PullRequestDetails []PullRequestResponse `json:"pull_request_details"`
}

type GithubRepositoryResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	OpenIssues  int    `json:"open_issues_count"`
	Private     bool   `json:"private"`
	Language    string `json:"language"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type PullRequestResponse struct {
	Id        int    `json:"id"`
	Url       string `json:"html_url"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IssueResponse struct {
	Id        int               `json:"id"`
	Repo      string            `json:"repo"`
	Url       string            `json:"html_url"`
	State     string            `json:"state"`
	Number    int               `json:"number"`
	Title     string            `json:"title"`
	Labels    []LabelResponse   `json:"labels"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	PRLinks   PullLinksResponse `json:"pull_request"`
}

type LabelResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type PullLinksResponse struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
}

func (r *RepoResponse) ToArray() []string {
	array := []string{r.Name, strconv.Itoa(r.Issues), strconv.Itoa(r.PullRequests)}
	return array
}
