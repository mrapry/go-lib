package goson

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	missmatchErrorPattern = `parsing value "%v": can't set value type %s to target type %s`
)

var (
	// check basic data type (integer, float, string, boolean)
	uintCheck = map[reflect.Kind]bool{
		reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
	}
	intCheck = map[reflect.Kind]bool{
		reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true,
	}
	floatCheck = map[reflect.Kind]bool{
		reflect.Float32: true, reflect.Float64: true,
	}
	stringCheck = map[reflect.Kind]bool{
		reflect.String: true,
	}
	boolCheck = map[reflect.Kind]bool{
		reflect.Bool: true,
	}
)

func parseFromString(target reflect.Value, str string) (err error) {
	targetKind := target.Kind()
	switch {
	case stringCheck[targetKind]:
		target.SetString(str)
	case intCheck[targetKind]:
		var val int
		if val, err = strconv.Atoi(str); err == nil {
			target.SetInt(int64(val))
		}
	case uintCheck[targetKind]:
		var val int
		if val, err = strconv.Atoi(str); err == nil {
			target.SetUint(uint64(val))
		}
	case floatCheck[targetKind]:
		var val float64
		if val, err = strconv.ParseFloat(str, -1); err == nil {
			target.SetFloat(val)
		}
	case boolCheck[targetKind]:
		var val bool
		if val, err = strconv.ParseBool(str); err == nil {
			target.SetBool(val)
		}
	default:
		err = fmt.Errorf(missmatchErrorPattern, str, "string", targetKind)
	}
	return
}

func parseFromFloat(target reflect.Value, fl float64) (err error) {
	targetKind := target.Kind()
	switch {
	case floatCheck[targetKind]:
		target.SetFloat(fl)
	case intCheck[targetKind]:
		target.SetInt(int64(fl))
	case uintCheck[targetKind]:
		target.SetUint(uint64(fl))
	case stringCheck[targetKind]:
		target.SetString(strconv.FormatFloat(fl, 'f', -1, 64))
	case boolCheck[targetKind]:
		var v bool
		if v, err = strconv.ParseBool(strconv.FormatFloat(fl, 'f', -1, 64)); err == nil {
			target.SetBool(v)
		}
	default:
		err = fmt.Errorf(missmatchErrorPattern, fl, "float64", targetKind)
	}
	return
}

func parseFromBool(target reflect.Value, bl bool) (err error) {
	targetKind := target.Kind()
	switch {
	case boolCheck[targetKind]:
		target.SetBool(bl)
	case stringCheck[targetKind]:
		target.SetString(strconv.FormatBool(bl))
	case uintCheck[targetKind]:
		target.SetUint(uint64(toInt(bl)))
	case intCheck[targetKind]:
		target.SetInt(int64(toInt(bl)))
	case floatCheck[targetKind]:
		target.SetFloat(float64(toInt(bl)))
	default:
		err = fmt.Errorf(missmatchErrorPattern, bl, "boolean", targetKind)
	}
	return
}

// set if true then int is 1, if false then int is 0
func toInt(b bool) (i int) {
	if b {
		i = 1
	}
	return
}
