package mapper

import (
	"context"
	"dto-mapper/constants"
	jsonUtil "dto-mapper/gjson"
	"dto-mapper/reflection"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"reflect"
)

func MapFrom[TSource, TTarget any](ctx context.Context, source TSource, target TTarget) (*TTarget, error) {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "Mapper",
			"Method": "MapFrom",
		}).
		Logger()

	sourceJson, marshallErr := json.Marshal(source)
	if marshallErr != nil {
		logger.Error().Msgf("Error while marshaling source. Error := %v", marshallErr)
		return nil, marshallErr
	}

	targetType := reflect.TypeOf(target)

	if invalidTypeErr := reflection.IsGivenTypeStruct(ctx, targetType); invalidTypeErr != nil {
		logger.Error().Msgf("Error while validating target type. Error := %v", invalidTypeErr)
		return nil, invalidTypeErr
	}

	logger.Info().Msgf("Getting tagged fields with constants %s from target", constants.TAG_MAP_FROM)
	targetTaggedFields := reflection.GetTagFieldMap(ctx, targetType, constants.TAG_MAP_FROM)

	targetInstance := reflection.InitializeStruct[TTarget](ctx, targetType)

	sourceResult := gjson.ParseBytes(sourceJson)

	var setErr error = nil

	for path, fieldName := range targetTaggedFields {
		targetFieldValue := reflection.GetFieldValue(ctx, targetInstance, fieldName)

		sourceValue := jsonUtil.GetValueBasedOnType(sourceResult.Get(path), targetFieldValue.Kind())

		setErr = reflection.SetValueInField(ctx, targetFieldValue, sourceValue)
		if setErr != nil {
			logger.Error().Msgf("Error while setting value %v into field %s. Error := %v", sourceValue, fieldName, setErr)
		}
	}

	return targetInstance, nil
}
