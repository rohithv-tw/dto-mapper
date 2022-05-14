package mapper

import (
	"context"
	"dto-mapper/constants"
	jsonUtil "dto-mapper/gjson"
	"dto-mapper/reflection"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"reflect"
)

func MapFrom[TSource, TTarget any](ctx context.Context, source TSource, targetInstancePointer *TTarget) error {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "Mapper",
			"Method": "MapFrom",
		}).
		Logger()

	logger.Info().Msgf("Finding type of target instance pointer")
	targetType := reflection.FindIndirectType(ctx, reflect.ValueOf(targetInstancePointer))

	logger.Info().Msgf("Check whether target type is struct")
	if !reflection.IsGivenTypeStruct(ctx, targetType) {
		invalidTypeErr := fmt.Errorf("target instance is not pointer of a struct instance. target type %s", targetType.Name())
		logger.Error().Msgf("Error while validating target type. Error := %v", invalidTypeErr)
		return invalidTypeErr
	}

	logger.Info().Msgf("Marshal source for extracting value from tags")
	sourceJson, marshallErr := json.Marshal(source)
	if marshallErr != nil {
		logger.Error().Msgf("Error while marshaling source. Error := %v", marshallErr)
		return marshallErr
	}

	logger.Info().Msgf("Getting tagged fields with constants %s from target", constants.TAG_MAP_FROM)
	targetTaggedFields := reflection.GetTagFieldMap(ctx, targetType, constants.TAG_MAP_FROM)

	sourceResult := gjson.ParseBytes(sourceJson)

	var setterErr error = nil

	for path, fieldName := range targetTaggedFields {
		targetFieldValue := reflection.GetFieldValue(ctx, targetInstancePointer, fieldName)

		sourceValue := jsonUtil.GetValueBasedOnType(sourceResult.Get(path), targetFieldValue.Kind())

		setterErr = reflection.SetValueInField(ctx, targetFieldValue, sourceValue)
		if setterErr != nil {
			logger.Error().Msgf("Error while setting value %v into field %s. Error := %v", sourceValue, fieldName, setterErr)
		}
	}

	return nil
}
