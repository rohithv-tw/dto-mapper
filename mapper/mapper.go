package mapper

import (
	"context"
	"dto-mapper/constants"
	jsonUtil "dto-mapper/gjson"
	logUtil "dto-mapper/log"
	"dto-mapper/reflection"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"reflect"
)

func MapFrom[TSource, TTarget any](ctx context.Context, source TSource, targetInstancePointer *TTarget) error {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "Mapper",
		"Method": "MapFrom",
	})

	logger.Info("Finding type of target instance pointer")
	targetType := reflection.FindIndirectType(ctx, reflect.ValueOf(targetInstancePointer))

	logger.Info("Check whether target type is struct")
	if !reflection.IsGivenTypeStruct(ctx, targetType) {
		invalidTypeErr := fmt.Errorf("target instance is not pointer of a struct instance. target type %s", targetType.Name())
		logger.Errorf("Error while validating target type. Error := %v", invalidTypeErr)
		return invalidTypeErr
	}

	logger.Info("Marshal source for extracting value from tags")
	sourceJson, marshallErr := json.Marshal(source)
	if marshallErr != nil {
		logger.Errorf("Error while marshaling source. Error := %v", marshallErr)
		return marshallErr
	}

	logger.Infof("Getting tagged fields with constants %s from target", constants.TAG_MAP_FROM)
	targetTaggedFields := reflection.GetTagFieldMap(ctx, targetType, constants.TAG_MAP_FROM)

	sourceResult := gjson.ParseBytes(sourceJson)

	for path, fieldName := range targetTaggedFields {
		targetFieldValue := reflection.GetFieldValue(ctx, targetInstancePointer, fieldName)

		sourceValue := jsonUtil.GetValueBasedOnType(sourceResult.Get(path), targetFieldValue.Kind())

		if ok := reflection.SetValueInField(ctx, targetFieldValue, sourceValue); !ok {
			logger.Errorf("Error while setting value %v into field %s", sourceValue, fieldName)
		}
	}

	return nil
}
