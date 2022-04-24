package main

import (
	"context"
	"dto-mapper/dto"
	"dto-mapper/mapper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func main() {
	logger := zerolog.New(os.Stdout)
	ctx := logger.WithContext(context.Background())

	logger = log.Ctx(ctx).With().Timestamp().
		Fields(map[string]interface{}{
			"Class":  "main",
			"Method": "main",
		}).
		Logger()

	input := dto.SourceDto{
		Name: dto.Name{
			First:  "first-input",
			Middle: "middle-input",
			Last:   "last-input",
		},
		Age:       99,
		CreatedAt: time.Now(),
		Friends: []dto.Name{
			{
				First:  "f-f-1",
				Middle: "f-m-1",
				Last:   "f-l-1",
			},
			{
				First:  "f-f-2",
				Middle: "f-m-2",
				Last:   "f-l-2",
			},
		},
	}

	output := dto.TargetDto{}

	log.Info().Msg("map from source to target")
	output2, err := mapper.MapFrom(ctx, input, output)
	if err != nil {
		log.Error().Msgf("Error while mapping.Error := %v", err)
		return
	}

	log.Info().Msgf("Mapped response := %v", *output2)
}
