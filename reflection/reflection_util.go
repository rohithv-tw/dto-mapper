package reflection

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

func FindIndirectType(ctx context.Context, value reflect.Value) reflect.Type {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "FindIndirectType",
		}).
		Logger()

	if value.Kind() == reflect.Ptr {
		logger.Info().Msgf("Found pointer, finding indirect type")
		return FindIndirectType(ctx, reflect.Indirect(value))
	}

	logger.Info().Msgf("Found non pointer type. returning")
	return value.Type()
}

func IsGivenTypeStruct(ctx context.Context, givenType reflect.Type) bool {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "IsGivenTypeStruct",
			"Type":   givenType.Name(),
		}).
		Logger()

	logger.Error().Msgf("given type is %s", givenType.Name())
	return givenType.Kind() == reflect.Struct
}

func GetTagFieldMap(ctx context.Context, structType reflect.Type, tag string) (taggedFields map[string]string) {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "InitializeStruct",
			"Type":   structType.Name(),
			"Tag":    tag,
		}).
		Logger()

	logger.Info().Msg("Extracting tags")
	taggedFields = make(map[string]string, 0)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if mapFromValue, ok := field.Tag.Lookup(tag); ok {
			taggedFields[mapFromValue] = field.Name
		}
	}

	log.Info().Msgf("Returning mapper constants %v successfully", taggedFields)
	return taggedFields
}

func findSetter(ctx context.Context, kind reflect.Kind) func(value reflect.Value, fieldValue interface{}) {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "InitializeStruct",
			"Kind":   kind.String(),
		}).
		Logger()

	logger.Info().Msg("Returning setter func")

	switch kind {

	case reflect.Slice:
		return func(value reflect.Value, fieldValue interface{}) {
			fieldArray := fieldValue.([]interface{})
			fieldArrayLength := len(fieldArray)
			slice := reflect.MakeSlice(value.Type(), fieldArrayLength, fieldArrayLength)
			for i := 0; i < fieldArrayLength; i++ {
				element := slice.Index(i)
				setter := findSetter(ctx, element.Kind())
				setter(element, fieldArray[i])
			}
			value.Set(slice)
		}

	default:
		return func(value reflect.Value, fieldValue interface{}) {
			reflectValue := reflect.ValueOf(fieldValue)
			value.Set(reflectValue)
		}
	}

}

func GetFieldValue[T any](ctx context.Context, instance *T, field string) reflect.Value {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "GetFieldValue",
			"Field":  field,
		}).
		Logger()

	logger.Info().Msg("Finding field with given name")
	return reflect.ValueOf(instance).Elem().FieldByName(field)
}

func SetValueInField(ctx context.Context, v reflect.Value, fieldValue interface{}) error {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "InitializeStruct",
		}).
		Logger()

	if v.CanSet() {
		setter := findSetter(ctx, v.Kind())
		setter(v, fieldValue)
		logger.Info().Msgf("Set value successfully")
		return nil
	}

	invalidSetterErr := fmt.Errorf("unable to set value into field")
	log.Err(invalidSetterErr)
	return invalidSetterErr
}
