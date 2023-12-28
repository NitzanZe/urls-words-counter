package config

type Config struct {
	Logger  LoggerConfig
	General GeneralConfiguration
}

type LoggerConfig struct {
	LogFormatting       string `env:"LOG_FORMATTING" default:"console"` //Modes: json, console (console is for development needs for nicer look on the console)
	LogPath             string `env:"LOG_PATH" default:"/var/logs"`
	LogEnableStdOutput  bool   `env:"LOG_ENABLE_STD_OUTPUT" default:"true"`
	LogEnableFileOutput bool   `env:"LOG_ENABLE_FILE_OUTPUT" default:"true"`
	LogLevel            string `env:"LOG_LEVEL" default:"info"`
}

type GeneralConfiguration struct {
	MaximumWorkersRequestsPerSecond int `env:"MAXIMUM_WORKERS_REQUESTS_SECONDS" default:"10"`
	WordsFileFullPath               string
	UrlsFileFullPath                string
	GetTopNWords                    int `env:"GET_TOP_N_WORDS" default:"10"`
	WorkersMultiplier               int `env:"WORKERS_MULTIPLIER" default:"1"` // The number of workers will be determined by the CPU cores. Setting this to more than one will allow more go-routines for each core
}
