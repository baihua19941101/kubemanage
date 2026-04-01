package main

import (
	"log"

	"kubeManage/backend/internal/config"
	"kubeManage/backend/internal/infra"
	"kubeManage/backend/internal/server"
)

func main() {
	cfg := config.Load()
	store, err := infra.NewStore(cfg)
	if err != nil {
		log.Fatalf("init store failed: %v", err)
	}
	defer func() {
		if cerr := store.Close(); cerr != nil {
			log.Printf("close store failed: %v", cerr)
		}
	}()

	r := server.NewRouter(store, cfg.K8sAdapterMode)

	log.Printf("kubeManage backend start on %s", cfg.ListenAddr)
	if err := r.Run(cfg.ListenAddr); err != nil {
		log.Fatalf("backend stopped with error: %v", err)
	}
}
