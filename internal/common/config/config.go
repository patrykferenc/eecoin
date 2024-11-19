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
