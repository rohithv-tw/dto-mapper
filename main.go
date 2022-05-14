package main

import (
	"dto-mapper/dto"
	logUtil "dto-mapper/log"
	"dto-mapper/mapper"
	"fmt"
	"time"
)

func init() {

}

func main() {
	ctx := logUtil.NewLoggerInContext()
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "main",
		"Method": "main",
	})

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
	logger.Info("map from source to target")
	err := mapper.MapFrom(ctx, input, &output)
	if err != nil {
		logger.Errorf("Error while mapping.Error := %v", err)
		return
	}

	fmt.Printf("output := %v\n", output)
}
