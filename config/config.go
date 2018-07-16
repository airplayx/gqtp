package config

var (
	Url     string
	FileUrl string
	Suffix  string
)

func init() {
	Url = "http://www.gqtp.com/"
	FileUrl = "http://images.gqtp.com/"
	Suffix = ".html"
}
