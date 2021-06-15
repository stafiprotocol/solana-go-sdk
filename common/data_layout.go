package common

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

func SerializeData(data interface{}) ([]byte, error) {
	return serializeData(reflect.ValueOf(data))
}

func serializeData(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	case reflect.Uint8:
		return []byte{uint8(v.Uint())}, nil
	case reflect.Int16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v.Int()))
		return b, nil
	case reflect.Uint16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v.Uint()))
		return b, nil
	case reflect.Int32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v.Int()))
		return b, nil
	case reflect.Uint32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v.Uint()))
		return b, nil
	case reflect.Int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(v.Int()))
		return b, nil
	case reflect.Uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v.Uint())
		return b, nil
	case reflect.Slice:
		data := make([]byte, 0)
		lenData, err := serializeData(reflect.ValueOf(uint32(v.Len())))
		if err != nil {
			return nil, err
		}
		data = append(data, lenData...)

		for i := 0; i < v.Len(); i++ {
			d, err := serializeData(v.Index(i))
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
		return data, nil
	case reflect.Array:
		data := make([]byte, 0)
		for i := 0; i < v.Len(); i++ {
			d, err := serializeData(v.Index(i))
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
		return data, nil

	case reflect.String:
		return []byte(v.String()), nil
	case reflect.Struct:
		data := make([]byte, 0, 1024)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			d, err := serializeData(field)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
		return data, nil
	}
	return nil, fmt.Errorf("unsupport type: %v", v.Kind())
}
