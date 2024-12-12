package cmd

import (
	"flag"

	"github.com/hilmiikhsan/multifinance-service/db/seeds"
	"github.com/hilmiikhsan/multifinance-service/internal/adapter"
	"github.com/rs/zerolog/log"
)

func RunSeed(cmd *flag.FlagSet, args []string) {
	var (
		table = cmd.String("table", "", "seed to run")
		total = cmd.Int("total", 1, "total of records to seed")
	)

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	adapter.Adapters.Sync(
		adapter.WithMultifinanceMySQL(),
	)

	defer func() {
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
	}()

	seeds.Execute(adapter.Adapters.MultifinanceMysql, *table, *total)
}
