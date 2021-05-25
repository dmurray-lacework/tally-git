# tally-git

A golang tool for viewing issues and pull requests on Open Source Repositiories managed by Tech Ally Team.

```
+--------------------------------+--------+---------------+
|           REPOSITORY           | ISSUES | PULL REQUESTS |
+--------------------------------+--------+---------------+
| terraform-provider-lacework    |      2 |             1 |
| chef-lacework                  |      0 |             0 |
| go-sdk                         |     29 |             4 |
| terraform-provisioning         |     11 |             0 |
| circleci-orb-lacework          |      0 |             0 |
| terraform-gcp-config           |      1 |             0 |
| terraform-azure-config         |      0 |             0 |
| terraform-aws-config           |      0 |             0 |
| terraform-aws-cloudtrail       |      3 |             3 |
| terraform-gcp-service-account  |      0 |             0 |
| terraform-aws-iam-role         |      0 |             0 |
| terraform-gcp-audit-log        |      0 |             0 |
| terraform-azure-ad-application |      1 |             1 |
| terraform-azure-activity-log   |      0 |             0 |
| homebrew-tap                   |      0 |             0 |
| terraform-kubernetes-agent     |      1 |             1 |
| terraform-aws-ssm-agent        |      0 |             0 |
+--------------------------------+--------+---------------+
```

## Usage

To successfully request data from the Github Api a personal Github token is required. See docs [Here](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)

Set your token in `example-config.yml` and rename to `config.yml`

From project root run `go run main.go`
### Project Structure
```
tally-git/

    ├── LICENSE
    ├── README.md
    ├── api/
    ├── config/
    ├── config.yml
    ├── go.mod
    ├── go.sum
    ├── internal/
    └── main.go
```