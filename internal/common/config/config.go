package config

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Peers       Peers       `yaml:"peers"`
	Log         Log         `yaml:"log"`
	Persistence Persistence `yaml:"persistence"`
}

type Peers struct {
	FilePath           string        `yaml:"filePath" env:"PEERS_FILE_PATH" env-default:"/etc/eecoin/peers"`
	PingDuration       time.Duration `yaml:"pingDuration" env:"PEERS_PING_DURATION" env-default:"5s"`
	UpdateFileDuration time.Duration `yaml:"updateFileDuration" env:"PEERS_UPDATE_FILE_DURATION" env-default:"1m"`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}

type Persistence struct {
	ChainFilePath      string        `yaml:"chainPath" env:"CHAIN_FILE_PATH" env-default:"/etc/eecoin/chain"`
	UpdateFileDuration time.Duration `yaml:"updateFileDuration" env:"CHAIN_UPDATE_FILE_DURATION" env-default:"1m"`
	SelfKey            string        `yaml:"selfKey" env:"SELF_KEY" env-default:"3059301306072a8648ce3d020106082a8648ce3d03010703420004fd957c299f6532aa445fc33f3fc87a7e9d5b8e32e0e9faaf8e38f706afdb6751a127cefe9e07fcca442e1053956fefdcb3bd8b412e7aade982638a3792890ed0"`
}

func (l *Log) LevelIfSet() (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(l.Level))
	return level, err
}

func Read(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	return &cfg, err
}
