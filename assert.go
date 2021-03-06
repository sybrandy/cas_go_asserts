package cas_go_asserts

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/andreyvit/diff"
)

type Assert struct {
	t     *testing.T
	level int
}

func (a Assert) logError(expected, actual interface{}) {
	var msg string
	expectedType := "nil"
	actualType := "nil"
	if expected != nil {
		expectedType = reflect.TypeOf(expected).String()
	}
	if actual != nil {
		actualType = reflect.TypeOf(actual).String()
	}
	if expectedType == "string" && actualType == "string" {
		a, _ := expected.(string)
		b, _ := actual.(string)
		msg = fmt.Sprintf("Expected: %+v, Actual: %+v, Diff: %s",
			expected, actual, diff.CharacterDiff(a, b))
	} else {
		msg = fmt.Sprintf("Expected: %+v (%s), Actual: %+v (%s)",
			expected, expectedType, actual, actualType)
	}
	if a.level == 0 {
		a.t.Error(msg)
	} else if a.level == 1 {
		a.t.Log(msg)
	}
}

/*
reflect.Kind Constants
const (
    Invalid Kind = iota
    Bool
    Int
    Int8
    Int16
    Int32
    Int64
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
    Uintptr
    Float32
    Float64
    Complex64
    Complex128
    Array
    Chan
    Func
    Interface
    Map
    Ptr
    Slice
    String
    Struct
    UnsafePointer
)
*/
func (a Assert) isSupported(varKind reflect.Kind) bool {
	return varKind != reflect.Invalid && (varKind <= reflect.Complex128 ||
		varKind == reflect.String || a.isArray(varKind))
}

func (a Assert) isArray(varKind reflect.Kind) bool {
	return varKind == reflect.Array || varKind == reflect.Slice
}

func (a Assert) checkArray(expected, actual interface{}) bool {
	exp := reflect.ValueOf(expected)
	act := reflect.ValueOf(actual)

	if exp.Len() != act.Len() {
		msg := fmt.Sprintf("Expected and acutal arrays are of a different length: %d vs. %d", exp.Len(), act.Len())
		if a.level == 0 {
			a.t.Error(msg)
		} else if a.level == 1 {
			a.t.Log(msg)
		}
		return false
	}

	for i := 0; i < exp.Len(); i++ {
		if !a.Equals(exp.Index(i).Interface(), act.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func (a Assert) Equals(expected, actual interface{}) bool {
	if (expected == nil || actual == nil) && expected != actual {
		a.logError(expected, actual)
		return false
	} else if expected == nil && actual == nil {
		return true
	}

	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	if !a.isSupported(expectedType.Kind()) || !a.isSupported(actualType.Kind()) {
		msg := fmt.Sprintf("Unsupported type in comparison: %s, %s", expectedType.Kind(), actualType.Kind())
		if a.level == 0 {
			a.t.Error(msg)
		} else if a.level == 1 {
			a.t.Log(msg)
		}
		return false
	}

	if expectedType != actualType {
		a.logError(expected, actual)
		return false
	}

	if a.isArray(expectedType.Kind()) && a.isArray(actualType.Kind()) {
		return a.checkArray(expected, actual)
	} else if expected != actual {
		a.logError(expected, actual)
		return false
	}
	return true
}

func (a Assert) IsNil(actual interface{}) bool {
	return a.Equals(nil, actual)
}

func (a Assert) HasError(expected string, actual error) bool {
	if actual == nil {
		a.logError(expected, actual)
		return false
	}
	return a.Equals(expected, actual.Error())
}
