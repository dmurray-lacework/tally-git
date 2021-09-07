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
			repo.Release = gs.GetLatestVersion(repo.Name)
			if repo.Release.Published != "" {
				repo.Commits = gs.GetLatestCommits(repo.Name, repo.Release.Published)
			}
			openSourceRepos = append(openSourceRepos, repo)
		}
	}
	return openSourceRepos
}

func (gs *GithubService) GetPullRequests(repoName string) []PullRequestResponse {
	response := &[]PullRequestResponse{}
	pullRequests := []PullRequestResponse{}
	apiPath := fmt.Sprintf(PR_URL, repoName)
	gs.client.RequestDecoder("GET", apiPath, nil, response)

	for _, r := range *response {
		pr := PullRequestResponse(r)
		pr.Reviews = append(pr.Reviews, gs.GetReviews(repoName, r.Number)...)
		pullRequests = append(pullRequests, pr)
	}

	return pullRequests
}

func (gs *GithubService) GetLatestVersion(repoName string) ReleaseResponse {
	response := &ReleaseResponse{}
	apiPath := fmt.Sprintf(LATEST_RELEASES_URL, repoName)
	gs.client.RequestDecoder("GET", apiPath, nil, response)
	return *response
}

func (gs *GithubService) GetLatestCommits(repoName string, published string) []CommitsResponse {
	response := &[]CommitsResponse{}
	apiPath := fmt.Sprintf(LATEST_COMMITS_URL, repoName, published)
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
			issue.Repo = repoName
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

func (gs *GithubService) GetReviews(repoName string, pullNumber int) []ReviewResponse {
	response := &[]ReviewResponse{}
	apiPath := fmt.Sprintf(REVIEWS_URL, repoName, pullNumber)
	gs.client.RequestDecoder("GET", apiPath, nil, response)
	return *response
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
	Release            ReleaseResponse       `json:"release"`
	Commits            []CommitsResponse     `json:"commits"`
}

type GithubRepositoryResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	OpenIssues  int    `json:"open_issues_count"`
	Private     bool   `json:"private"`
	Language    string `json:"language"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type PullRequestResponse struct {
	Id                int              `json:"id"`
	Name              string           `json:"title"`
	Number            int              `json:"number"`
	Body              string           `json:"repo"`
	Assignee          string           `json:"assignee,omitempty"`
	Url               string           `json:"html_url"`
	State             string           `json:"state"`
	AuthorAssociation string           `json:"author_association"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
	User              UserResponse     `json:"user"`
	Head              HeadResponse     `json:"head"`
	Reviews           []ReviewResponse `json:"reviews,omitempty"`
}

type ReviewResponse struct {
	Id                int          `json:"id"`
	State             string       `json:"state"`
	Body              string       `json:"body"`
	Url               string       `json:"html_url"`
	AuthorAssociation string       `json:"author_association"`
	SubmittedAt       string       `json:"submitted_at"`
	User              UserResponse `json:"user"`
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

type HeadResponse struct {
	Id   int          `json:"id"`
	Repo RepoResponse `json:"repo"`
}

type PullLinksResponse struct {
	Url      string `json:"url"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
}

type CommitsResponse struct {
	Sha     string `json:"sha"`
	HtmlUrl string `json:"html_url"`
	Commit  Commit `json:"commit"`
}

type Commit struct {
	Message      string `json:"message"`
	Author       Author `json:"author"`
	CommentCount int    `json:"comment_count"`
}
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type ReleaseResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag_name"`
	Published string `json:"published_at"`
	Url       string `json:"html_url"`
	Body      string `json:"body"`
}

func (r *RepoResponse) ToArray() []string {
	array := []string{r.Name, strconv.Itoa(r.Issues), strconv.Itoa(r.PullRequests)}
	return array
}
