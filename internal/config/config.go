package config

var C Config

type Mirrors struct {
	Mirrors  []string `json:"mirrors"`
	Platform string   `json:"platform"`
}

type Config struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Server   string    `json:"server"`
	Maximum  int       `json:"maximum"`
	Output   string    `json:"output"`
	Mirrors  []Mirrors `json:"mirrors"`
}
