package main

import (
	"fmt"
	"log/slog"
	"os"

	"kubeManage/backend/internal/config"
	"kubeManage/backend/internal/infra"
	"kubeManage/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %v\n", err)
		os.Exit(1)
	}

	closeLogger, err := infra.SetupLogger(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logger failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := closeLogger(); cerr != nil {
			fmt.Fprintf(os.Stderr, "close logger failed: %v\n", cerr)
		}
	}()

	store, err := infra.NewStore(cfg)
	if err != nil {
		slog.Error("init store failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := store.Close(); cerr != nil {
			slog.Warn("close store failed", "error", cerr)
		}
	}()

	r := server.NewRouter(store, cfg.K8sAdapterMode, cfg.SecretKey)

	slog.Info("kubeManage backend start", "listenAddr", cfg.ListenAddr, "configFile", cfg.ConfigFile)
	if err := r.Run(cfg.ListenAddr); err != nil {
		slog.Error("backend stopped with error", "error", err)
		os.Exit(1)
	}
}
