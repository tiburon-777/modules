package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func getEnvVar(v reflect.Value, t reflect.Type, counter int, prefix string) error {
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer value")
	}
	f := reflect.StructField{}
	if counter != -1 {
		f = t.Field(counter)
	}
	v = reflect.Indirect(v)
	fName := strings.ToUpper(f.Name)
	if prefix != "" {
		prefix = strings.TrimLeft(prefix, "_")
	}

	env := os.Getenv(prefix + fName)
	if err := selector(env, &v); err != nil {
		return fmt.Errorf("could set value: %w", err)
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			if err := getEnvVar(v.Field(i).Addr(), v.Type(), i, prefix+fName+"_"); err != nil {
				return fmt.Errorf("could not apply env var: %w", err)
			}
		}
	}
	return nil
}

func parseToStruct(v reflect.Value, t reflect.Type, counter int, prefix string, kv map[string]string) error {
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer value")
	}
	f := reflect.StructField{}
	if counter != -1 {
		f = t.Field(counter)
	}
	v = reflect.Indirect(v)
	fName := strings.ToLower(f.Name)
	if prefix != "" {
		prefix = strings.TrimLeft(prefix, ".")
	}

	env := kv[prefix+fName]
	if err := selector(env, &v); err != nil {
		return fmt.Errorf("could set value: %w", err)
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			if err := parseToStruct(v.Field(i).Addr(), v.Type(), i, prefix+fName+".", kv); err != nil {
				return fmt.Errorf("could not apply env var: %w", err)
			}
		}
	}
	return nil
}

func selector(env string, v *reflect.Value) error {
	if env != "" {
		switch v.Kind() {
		case reflect.Int:
			envI, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("could not parse to int: %w", err)
			}
			v.SetInt(int64(envI))
		case reflect.String:
			v.SetString(env)
		case reflect.Bool:
			envB, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("could not parse bool: %w", err)
			}
			v.SetBool(envB)
		case reflect.Array, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64, reflect.Func, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface, reflect.Invalid, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		}
	}
	return nil
}
