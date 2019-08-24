package config

// Config stores the Appilcation Configuration
type Config struct {
	Hostname    string
	Environment string
	Server      struct {
		Port int
	}
	Log struct {
		Level string // one of debug, info, warn, error, fatal
	}
	Application struct {
		Name string
	}
}
