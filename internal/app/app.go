package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanmu42/gzip"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strconv"

	"test-zero-agency/internal/api/handler"
	"test-zero-agency/internal/api/middleware"
	"test-zero-agency/internal/app/config"
	"test-zero-agency/internal/entity"
	"test-zero-agency/internal/getaway"
	"test-zero-agency/internal/repository"
	"test-zero-agency/internal/service"
)

func Run(ctx context.Context, cfg *config.Config) {
	closer := newCloser()
	logger := newLogger()
	router := chi.NewRouter()
	db := newDataBase(ctx, logger, cfg)

	server := newServer(cfg, router)
	closer.Add(server.Shutdown)
	defer db.Close()

	err := db.Ping(ctx)
	if err != nil {
		panic(err.Error())
	}

	// Getaway
	enrichment := getaway.NewEnrichment(cfg)

	// Repository
	peopleRepository := repository.NewPeopleRepository(db)

	// Service
	peopleService := service.NewPeopleService(cfg, peopleRepository, enrichment)

	// API
	middleware := middleware.NewMiddleware(logger)
	handler.RegisterPeopleHandlers(router, peopleService, logger, middleware)

	go func() {
		logger.Info("starting the server", zap.String("Address", server.Addr))
		logger.DPanic("ListenAndServe", zap.Any("Error", server.ListenAndServe()))
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := closer.Close(shutdownCtx); err != nil {
		logger.Error("Close err", zap.Error(err))
	}
}

func newServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Handler:        gzip.DefaultHandler().WrapHandler(router),
		Addr:           ":" + strconv.Itoa(cfg.Port),
		WriteTimeout:   cfg.WriteTimeout,
		ReadTimeout:    cfg.ReadTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: 1 << 20,
	}
}

func newDataBase(ctx context.Context, logger *zap.Logger, cfg *config.Config) *pgxpool.Pool {

	conn := "postgres" + "://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Address + "/" + cfg.DBName + cfg.Params

	logger.Info("start establishing a connection to the database", zap.String("Connect", conn))

	db, err := pgxpool.New(ctx,
		conn)
	if err != nil {
		panic(err.Error())
	}
	//pgLogger := zapadapter.NewLogger(logger)
	//config.ConnConfig.Logger = pgLogger

	err = db.Ping(ctx)
	if err != nil {
		panic(err.Error())
	}

	logger.Info("successful connection to the database")

	return db
}

func newLogger() *zap.Logger {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger.Info("Zap Logger", zap.String("Level", logger.Level().String()))

	return logger
}

func newCloser() *entity.Closer {
	return &entity.Closer{}
}
