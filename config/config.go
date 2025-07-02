package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DSN                  string `env:"DSN,required"`
	CsvFile              string `env:"CSV_FILE,required"`
	DbMinConns           int    `env:"DB_MIN_CONNS" envDefault:"0"`
	DbMaxConns           int    `env:"DB_MAX_CONNS" envDefault:"30"`
	WorkerPoolGoroutines int    `env:"WORKER_POOL_GOROUTINES" envDefault:"16"`
	BatchSize            int    `env:"BATCH_SIZE" envDefault:"75"`
	CronSchedule         string `env:"CRON_SCHEDULE" envDefault:"@every 24h"`
	MigrationsDir        string `env:"MIGRATIONS_DIR" envDefault:"/root/migrations"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
