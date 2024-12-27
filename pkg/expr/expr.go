package expr

import (
	"strconv"
	"strings"
)

type Expression struct {
	Raw         string
	ValueString string
	Value       interface{}
	IsEvaluated bool
	Type        string
}

type Evaluator interface {
	Eval(template string, ctx map[string]interface{}) (string, error)
}

func (e *Expression) Eval(evaluator Evaluator, ctx map[string]interface{}) error {
	if !e.IsEvaluated {
		v, err := evaluator.Eval(e.Raw, ctx)
		if err != nil {
			return err
		}

		e.ValueString = v
		e.IsEvaluated = true
		switch e.Type {
		case "string":
			e.Value = e.ValueString
		case "uint32":
			e.Value = uint32(0)
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				i64, err := strconv.ParseInt(str, 10, 32)
				if err != nil {
					return err
				}

				e.Value = uint32(i64)
			}
		case "int32":
			e.Value = int32(0)
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				i64, err := strconv.ParseInt(str, 10, 32)
				if err != nil {
					return err
				}

				e.Value = int32(i64)
			}
		case "int":
			fallthrough
		case "int64":
			e.Value = int64(0)
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				i64, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}

				e.Value = i64
			}
		case "uint":
			fallthrough
		case "uint64":
			e.Value = uint64(0)
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				i64, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}

				e.Value = uint64(i64)
			}

		case "float32":
			e.Value = float32(0)
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				f64, err := strconv.ParseFloat(str, 32)
				if err != nil {
					return err
				}

				e.Value = float32(f64)
			}
		case "float64":
			e.Value = float64(0)
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				f64, err := strconv.ParseFloat(str, 64)
				if err != nil {
					return err
				}

				e.Value = f64
			}

		case "bool":
			e.Value = false
			str := strings.TrimSpace(e.ValueString)
			if str != "" {
				b, err := strconv.ParseBool(str)
				if err != nil {
					return err
				}

				e.Value = b
			}

		default:
		}
	}

	return nil
}

func (s *Expression) String() string {
	if !s.IsEvaluated || s.Type != "string" {
		return ""
	}

	return s.ValueString
}

func (s *Expression) Uint32() uint32 {
	if !s.IsEvaluated || s.Type != "uint32" {
		return 0
	}

	return s.Value.(uint32)
}

func (s *Expression) Uint64() uint64 {
	if !s.IsEvaluated || (s.Type != "uint64" && s.Type != "uint") {
		return 0
	}

	return s.Value.(uint64)
}

func (s *Expression) Int32() int32 {
	if !s.IsEvaluated || s.Type != "int32" {
		return 0
	}

	return s.Value.(int32)
}

func (s *Expression) Int64() int64 {
	if !s.IsEvaluated || (s.Type != "int64" && s.Type != "int") {
		return 0
	}

	return s.Value.(int64)
}

func (s *Expression) Float32() float32 {
	if !s.IsEvaluated || s.Type != "float32" {
		return 0
	}

	return s.Value.(float32)
}

func (s *Expression) Float64() float64 {
	if !s.IsEvaluated || s.Type != "float64" {
		return 0
	}

	return s.Value.(float64)
}
