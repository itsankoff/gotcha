package server

type Config struct {
	ListenHost       string
	FileServerHost   string
	FileServerPath   string
	FileServerFolder string
}

func NewConfig() *Config {
	return &Config{}
}
