package reflection

import (
	"context"
	logUtil "dto-mapper/log"
	"reflect"
)

func FindIndirectType(ctx context.Context, value reflect.Value) reflect.Type {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "ReflectionUtil",
		"Method": "FindIndirectType",
	})

	if value.Kind() == reflect.Ptr {
		logger.Info("Found pointer, finding indirect type")
		return FindIndirectType(ctx, reflect.Indirect(value))
	}

	logger.Info("Found non pointer type. returning")
	return value.Type()
}

func IsGivenTypeStruct(ctx context.Context, givenType reflect.Type) bool {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "ReflectionUtil",
		"Method": "IsGivenTypeStruct",
		"Type":   givenType.Name(),
	})

	logger.Errorf("given type is %s", givenType.Name())
	return givenType.Kind() == reflect.Struct
}

func GetTagFieldMap(ctx context.Context, structType reflect.Type, tag string) (taggedFields map[string]string) {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "ReflectionUtil",
		"Method": "InitializeStruct",
		"Type":   structType.Name(),
		"Tag":    tag,
	})

	logger.Info("Extracting tags")
	taggedFields = make(map[string]string, 0)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if mapFromValue, ok := field.Tag.Lookup(tag); ok {
			taggedFields[mapFromValue] = field.Name
		}
	}

	logger.Infof("Returning mapper constants %v successfully", taggedFields)
	return taggedFields
}

func findSetter(ctx context.Context, kind reflect.Kind) func(value reflect.Value, fieldValue interface{}) {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "ReflectionUtil",
		"Method": "InitializeStruct",
		"Kind":   kind.String(),
	})

	logger.Info("Returning setter func")

	switch kind {

	case reflect.Slice:
		return func(sourceField reflect.Value, fieldValue interface{}) {
			fieldValueAsArray := fieldValue.([]interface{})
			fieldValueAsArrayLength := len(fieldValueAsArray)
			slice := reflect.MakeSlice(sourceField.Type(), fieldValueAsArrayLength, fieldValueAsArrayLength)
			for i := 0; i < fieldValueAsArrayLength; i++ {
				element := slice.Index(i)
				setter := findSetter(ctx, element.Kind())
				setter(element, fieldValueAsArray[i])
			}
			sourceField.Set(slice)
		}

	default:
		return func(sourceField reflect.Value, fieldValue interface{}) {
			reflectValue := reflect.ValueOf(fieldValue)
			sourceField.Set(reflectValue)
		}
	}

}

func GetFieldValue[T any](ctx context.Context, instance *T, field string) reflect.Value {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "ReflectionUtil",
		"Method": "GetFieldValue",
		"Field":  field,
	})

	logger.Info("Finding field with given name")
	return reflect.ValueOf(instance).Elem().FieldByName(field)
}

func SetValueInField[T any](ctx context.Context, sourceField reflect.Value, fieldValue T) bool {
	logger := logUtil.GetLogger(ctx, map[string]interface{}{
		"Class":  "ReflectionUtil",
		"Method": "InitializeStruct",
	})

	if sourceField.CanSet() {
		setter := findSetter(ctx, sourceField.Kind())
		setter(sourceField, fieldValue)
		logger.Info("Set value successfully")
		return true
	}

	logger.Error("unable to set value into field")
	return false
}
