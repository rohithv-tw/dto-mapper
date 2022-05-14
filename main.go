package main

import (
	"context"
	"dto-mapper/dto"
	"dto-mapper/mapper"
	"fmt"
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
				First:  "first-1",
				Middle: "middle-1",
				Last:   "last-1",
			},
			{
				First:  "first-2",
				Middle: "middle-2",
				Last:   "last-2",
			},
		},
		Data: map[string]interface{}{
			"key1": 1,
		},
	}

	output := dto.TargetDto{}

	fmt.Printf("Input := %v\n", input)
	log.Info().Msg("map from source to target")
	err := mapper.MapFrom(ctx, input, &output)
	if err != nil {
		log.Error().Msgf("Error while mapping.Error := %v", err)
		return
	}

	fmt.Printf("output := %v\n", output)
}
