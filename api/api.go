package api

const (
	REPOS_URL           = "orgs/lacework/teams/tech-ally/repos?page=1&per_page=60"
	PR_URL              = "repos/lacework/%s/pulls?state=open"
	ISSUES_URL          = "repos/lacework/%s/issues?state=open"
	LATEST_RELEASES_URL = "repos/lacework/%s/releases/latest"
	LATEST_COMMITS_URL  = "repos/lacework/%s/commits?since=%s"
)
