package reflection

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

func IsGivenTypeStruct(ctx context.Context, givenType reflect.Type) error {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "IsGivenTypeStruct",
			"Type":   givenType.Name(),
		}).
		Logger()

	if givenType.Kind() != reflect.Struct {
		invalidTypeErr := fmt.Errorf("given type is not struct %s", givenType.Name())
		logger.Error().Msgf("Error while getting type of given struct.Error := %v", invalidTypeErr)
		return invalidTypeErr
	}

	logger.Info().Msgf("Given type is struct")
	return nil
}

func initializeStruct(ctx context.Context, t reflect.Type, v reflect.Value) {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "initializeStruct",
		}).
		Logger()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		logger.Info().Msgf("Found fieldValue %s with type %s", fieldType.Name, fieldType.Type.String())
		switch fieldType.Type.Kind() {
		case reflect.Map:
			fieldValue.Set(reflect.MakeMapWithSize(fieldType.Type, 0))
		case reflect.Slice:
			fieldValue.Set(reflect.MakeSlice(fieldType.Type, 0, 0))
		case reflect.Chan:
			fieldValue.Set(reflect.MakeChan(fieldType.Type, 0))
		case reflect.Struct:
			initializeStruct(ctx, fieldType.Type, fieldValue)
		case reflect.Ptr:
			fieldValue := reflect.New(fieldType.Type.Elem())
			initializeStruct(ctx, fieldType.Type.Elem(), fieldValue.Elem())
			fieldValue.Set(fieldValue)
		default:
			logger.Warn().Msgf("Set Zero Value for fieldValue %s with type %s", fieldType.Name, fieldType.Type.String())
			zero := reflect.Zero(fieldType.Type)
			fieldValue.Set(zero)
		}
	}
}

func InitializeStruct[T any](ctx context.Context, structType reflect.Type) *T {
	logger := log.Ctx(ctx).With().
		Timestamp().
		Fields(map[string]interface{}{
			"Class":  "ReflectionUtil",
			"Method": "InitializeStruct",
		}).
		Logger()

	targetInstanceValue := reflect.New(structType)
	initializeStruct(ctx, structType, targetInstanceValue.Elem())
	targetInstance := targetInstanceValue.Interface().(*T)

	logger.Info().Msgf("Create instance successfully for %s", structType.Name())
	return targetInstance
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
