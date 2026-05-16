package processor

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// CastType represents the target type for a field cast.
type CastType string

const (
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastString CastType = "string"
	CastBool   CastType = "bool"
)

// FieldCaster converts the value of a named field to a target type.
type FieldCaster struct {
	field    string
	castType CastType
}

// NewFieldCaster creates a FieldCaster for the given field and target type.
// Returns an error if field is blank or castType is unsupported.
func NewFieldCaster(field string, castType CastType) (*FieldCaster, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldcast: field name must not be empty")
	}
	switch castType {
	case CastInt, CastFloat, CastString, CastBool:
		// valid
	default:
		return nil, fmt.Errorf("fieldcast: unknown cast type %q (want int, float, string, bool)", castType)
	}
	return &FieldCaster{field: field, castType: castType}, nil
}

// Process attempts to cast the named field's value on the given LogLine.
// If the field is missing or the cast fails, the line is returned unchanged.
func (fc *FieldCaster) Process(line parser.LogLine) parser.LogLine {
	val, ok := line.Fields[fc.field]
	if !ok {
		return line
	}
	str := fmt.Sprintf("%v", val)
	var casted interface{}
	var err error
	switch fc.castType {
	case CastInt:
		casted, err = strconv.ParseInt(str, 10, 64)
	case CastFloat:
		casted, err = strconv.ParseFloat(str, 64)
	case CastString:
		casted = str
	case CastBool:
		casted, err = strconv.ParseBool(str)
	}
	if err != nil {
		return line
	}
	if line.Fields == nil {
		line.Fields = make(map[string]interface{})
	}
	line.Fields[fc.field] = casted
	return line
}
