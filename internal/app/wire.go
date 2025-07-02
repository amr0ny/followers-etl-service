//go:build wireinject
// +build wireinject

package app

import (
	"context"

	"github.com/amr0ny/followers-etl-service/config"
	"github.com/amr0ny/followers-etl-service/internal/common/logger/zaplogger"
	"github.com/amr0ny/followers-etl-service/internal/infrastructure/database"
	"github.com/amr0ny/followers-etl-service/internal/infrastructure/database/pg"
	"github.com/amr0ny/followers-etl-service/internal/infrastructure/fsmanager"
	"github.com/amr0ny/followers-etl-service/internal/infrastructure/workerpool"
	"github.com/amr0ny/followers-etl-service/internal/service"

	"github.com/google/wire"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

func InitializeApp(ctx context.Context, conf *config.Config) (*App, error) {
	wire.Build(
		LoggerSet,
		DBSet,
		FSManagerSet,
		WorkerPoolSet,
		ServiceSet,
		AppSet,
	)
	return &App{}, nil
}

//
// Wire sets
//

var LoggerSet = wire.NewSet(
	provideZapConfig,
	zaplogger.NewZapLogger,
)

var DBSet = wire.NewSet(
	providePgPoolOptions,
	database.NewPgPool,
	pg.NewFollowerRepositoryPG,
	pg.NewHelperRepositoryPG,
)

var FSManagerSet = wire.NewSet(
	provideLocalFSManagerFileName,
	provideLocalFSManagerBatchSize,
	fsmanager.NewLocalFSManager,
)

var WorkerPoolSet = wire.NewSet(
	provideAntsOptions,
	provideAntsErrGroupCapacity,
	workerpool.NewAntsErrGroup,
	wire.Struct(new(workerpool.AntsErrGroup), "*"),
)

var ServiceSet = wire.NewSet(
	service.NewFollowerService,
	service.NewHelperService,
	service.NewFSManagerService,
	service.NewWorkerPoolService,
)

var AppSet = wire.NewSet(
	wire.Struct(new(App), "*"),
)

//
// Providers
//

func provideZapConfig(conf *config.Config) zap.Config {
	return zap.NewDevelopmentConfig() // TODO: заменить на продакшн-конфиг
}

func providePgPoolOptions(conf *config.Config) database.PgPoolOptions {
	return database.PgPoolOptions{
		DSN:      conf.DSN,
		MinConns: int32(conf.DbMinConns),
		MaxConns: int32(conf.DbMaxConns),
	}
}

func provideLocalFSManagerFileName(conf *config.Config) string {
	return conf.CsvFile
}

func provideLocalFSManagerBatchSize(conf *config.Config) fsmanager.BatchSize {
	return fsmanager.BatchSize(conf.BatchSize)
}

func provideAntsErrGroupCapacity(conf *config.Config) workerpool.AntsErrGroupCapacity {
	return workerpool.AntsErrGroupCapacity(conf.WorkerPoolGoroutines)
}

func provideAntsOptions() []ants.Option {
	return []ants.Option{
		ants.WithPreAlloc(true),
	}
}
