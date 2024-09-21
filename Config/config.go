package Config

type AppConfig struct {
	Host string
	Port string
}

type Configs struct {
	App AppConfig
}

const Secret = "secret"
