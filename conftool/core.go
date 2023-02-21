package conftool

type confType string

const (
	Yaml confType = "yaml"
	Json confType = "json"
	Toml confType = "toml"
)
