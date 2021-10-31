package seekpo

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
}

func ParseType(s string) (Type, error) {
	if t, ok := parsePrimitiveType(s); ok {
		return t, nil
	}
	if t, ok := parseArrayType(s); ok {
		return t, nil
	}
	if t, ok := parseObjectType(s); ok {
		return t, nil
	}
	return nil, fmt.Errorf("missing type: '%s'", s)
}

func parsePrimitiveType(s string) (PrimitiveType, bool) {
	var (
		t  PrimitiveType
		ok bool
	)
	switch s {
	case TypeBool.String():
		t = TypeBool
	case TypeInt8.String():
		t = TypeInt8
	case TypeInt16.String():
		t = TypeInt16
	case TypeInt32.String():
		t = TypeInt32
	case TypeInt64.String():
		t = TypeInt64
	case TypeInt128.String():
		t = TypeInt128
	case TypeUint8.String():
		t = TypeUint8
	case TypeUint16.String():
		t = TypeUint16
	case TypeUint32.String():
		t = TypeUint32
	case TypeUint64.String():
		t = TypeUint64
	case TypeUInt128.String():
		t = TypeUInt128
	case TypeFloat32.String():
		t = TypeFloat32
	case TypeFloat64.String():
		t = TypeFloat64
	case TypeString.String():
		t = TypeString
	}
	return t, ok
}

func parseArrayType(s string) (ArrayType, bool) {
	if ok := strings.HasPrefix(s, arrayPrefix); !ok {
		return ArrayType{}, false
	}
	t, err := ParseType(strings.TrimPrefix(s, arrayPrefix))
	if err != nil {
		return ArrayType{}, false
	}
	return ArrayType{
		t: t,
	}, false
}

func parseObjectType(s string) (ObjectType, bool) {
	// TODO validate
	return ObjectType{
		s: s,
	}, true
}

type PrimitiveType string

const (
	TypeBool    PrimitiveType = "bool"
	TypeInt8    PrimitiveType = "i8"
	TypeInt16   PrimitiveType = "i16"
	TypeInt32   PrimitiveType = "i32"
	TypeInt64   PrimitiveType = "i64"
	TypeInt128  PrimitiveType = "i128"
	TypeUint8   PrimitiveType = "u8"
	TypeUint16  PrimitiveType = "u16"
	TypeUint32  PrimitiveType = "u32"
	TypeUint64  PrimitiveType = "u64"
	TypeUInt128 PrimitiveType = "u128"
	TypeFloat32 PrimitiveType = "f32"
	TypeFloat64 PrimitiveType = "f64"
	TypeString  PrimitiveType = "str"
)

func (t PrimitiveType) String() string {
	return string(t)
}

type ArrayType struct {
	t Type
}

const arrayPrefix = "[]"

func (t ArrayType) String() string {
	return fmt.Sprintf("%s%s", arrayPrefix, t.t)
}

type ObjectType struct {
	s string
}

func (t ObjectType) String() string {
	return t.s
}
