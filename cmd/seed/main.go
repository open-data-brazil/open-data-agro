// Package main seeds catalog.dataset_registry from the YAML registry.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/db"
)

func main() {
	mvpOnly := flag.Bool("mvp", false, "Seed only MVP datasets (conab.estimativa-graos, conab.serie-historica-graos)")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.LoadFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	reg, err := catalog.LoadDefaultRegistry()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	repo, err := db.NewRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer repo.Close()

	mvp := map[string]struct{}{
		"conab.estimativa-graos":      {},
		"conab.serie-historica-graos": {},
	}

	var seeded int
	for _, entry := range reg.Entries() {
		if *mvpOnly {
			if _, ok := mvp[entry.DatasetID.String()]; !ok {
				continue
			}
		}
		if err := repo.UpsertDatasetRegistry(ctx, entry); err != nil {
			fmt.Fprintf(os.Stderr, "seed %s: %v\n", entry.DatasetID, err)
			os.Exit(1)
		}
		seeded++
	}

	fmt.Printf("seeded %d dataset(s) into catalog.dataset_registry (portal=%s)\n",
		seeded, catalog.CONABSourcePortalURL)
}
