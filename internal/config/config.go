package config

import (
	"flag"
	"os"
)

type ServerConfiguration struct {
	RunAddr           string
	AccuralSystemAddr string
	DatabaseDSN       string
}

func LoadServerConfiguration() *ServerConfiguration {
	cfg := &ServerConfiguration{}

	if cfg.RunAddr = os.Getenv("RUN_ADDRESS"); cfg.RunAddr == "" {
		flag.StringVar(&cfg.RunAddr, "a", ":8080", "Server address")
	}

	if cfg.AccuralSystemAddr = os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); cfg.AccuralSystemAddr == "" {
		flag.StringVar(&cfg.AccuralSystemAddr, "r", "", "Accural system address")
	}

	if cfg.DatabaseDSN = os.Getenv("DATABASE_URI"); cfg.DatabaseDSN == "" {
		flag.StringVar(&cfg.DatabaseDSN, "d", "", "")
	}

	flag.Parse()

	return cfg
}
