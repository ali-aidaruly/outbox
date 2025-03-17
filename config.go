package outbox

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	dbType     string //  // Database type (e.g., DBPostgres, DBMySQL)
	hardDelete bool   // If true, messages are permanently deleted after processing.

	pollBatchSize int           // Number of messages fetched per polling cycle.
	pollInterval  time.Duration // Interval between polling cycles.
	leaseDuration time.Duration // Duration for which a message is locked before being retried.
}

// Option defines a functional option for configuring the outbox worker.
type Option func(*Config)

// WithDBType sets the database type.
func WithDBType(dbType string) Option {
	return func(cfg *Config) {
		cfg.dbType = dbType
	}
}

// WithHardDelete enables or disables hard deletion of processed messages.
func WithHardDelete(enabled bool) Option {
	return func(cfg *Config) {
		cfg.hardDelete = enabled
	}
}

// WithPollInterval sets the interval at which the worker polls for messages.
func WithPollInterval(interval time.Duration) Option {
	return func(cfg *Config) {
		cfg.pollInterval = interval
	}
}

// WithPollBatchSize sets the number of messages fetched per polling cycle.
func WithPollBatchSize(batchSize int) Option {
	return func(cfg *Config) {
		cfg.pollBatchSize = batchSize
	}
}

// WithLeaseDuration sets the lease duration for locking messages before retry.
func WithLeaseDuration(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.leaseDuration = d
	}
}

// WithDebugLogs enables debug-level logging for the Outbox library.
// It sets the global logger to use the debug level, allowing detailed logs.
func WithDebugLogs() Option {
	return func(_ *Config) {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
}
