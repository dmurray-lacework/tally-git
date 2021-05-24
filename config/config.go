package config

type Configuration struct {
	Github GithubConfig
}

type GithubConfig struct {
	Apikey string
}
