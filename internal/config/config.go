package config

var C Config

type HubMirrors struct {
	Content  []string `json:"mirror"`
	Platform string   `json:"platform"`
}

type Config struct {
	Content    string
	Maximum    int
	Repository string
	Username   string
	Password   string
	Output     string
	HubMirrors *HubMirrors
}
