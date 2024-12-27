package tasks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jolt9dev/go-jolt9/pkg/expr"
	"github.com/jolt9dev/go-jolt9/pkg/primitives"
	"github.com/jolt9dev/go-xstrings"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Id          string
	Name        string
	Uses        string
	Description string
	With        map[string]expr.Expression
	Env         map[string]expr.Expression
	Timeout     *expr.Expression
	Force       *expr.Expression
	If          *expr.Expression
	Cwd         *expr.Expression
	Needs       []string
	RunExpr     *expr.Expression
}

func (a Task) Compare(b Task) int {
	if a.Id == b.Id {
		return 0
	}

	if a.Id < b.Id {
		return -1
	}

	return 1
}

func (t *Task) Eval(ctx *TaskContext) error {

	if ctx.State == nil {
		ctx.State = &TaskState{
			Id:          t.Id,
			Name:        t.Name,
			Uses:        t.Uses,
			Description: t.Description,
			Inputs:      &primitives.ObjectMap{},
			Outputs:     &primitives.ObjectMap{},
			Force:       false,
			Timeout:     0,
			If:          true,
			Env:         make(map[string]string),
			Cwd:         "",
			RunExpr:     "",
		}
	}

	for k, v := range ctx.Env {
		ctx.State.Env[k] = v
	}

	data := make(map[string]interface{})
	data["env"] = ctx.Env
	data["secrets"] = ctx.Secrets
	data["outputs"] = mapOutputs(ctx.Outputs)
	data["vars"] = mapOutputs(ctx.Vars)

	if len(t.Env) > 0 {
		for key, value := range t.Env {
			if !value.IsEvaluated {
				err := value.Eval(ctx.Evaluator, data)
				if err != nil {
					return err
				}
			}

			ctx.Env[key] = value.Value.(string)
			ctx.State.Env[key] = value.Value.(string)
		}
	}

	if len(t.With) > 0 {
		if ctx.Descriptor == nil {
			return fmt.Errorf("inputs are not defined for task %s", t.Id)
		}

		for key, value := range t.With {
			in, ok := ctx.Descriptor.Inputs[key]
			if !ok {
				return fmt.Errorf("input %s is not defined for task %s", key, t.Id)
			}

			if !value.IsEvaluated {
				err := value.Eval(ctx.Evaluator, data)
				if err != nil {
					return err
				}
			}

			value := value.Value.(string)

			if in.IsRequired && value == "" {
				return fmt.Errorf("input %s is required for task %s", key, t.Id)
			}

			envName := xstrings.Underscore(key, xstrings.Screaming)
			envName = "INPUT_" + envName
			ctx.Env[envName] = value

			switch in.Type {
			case "int":
				fallthrough
			case "int64":
				i := int64(0)
				if value != "" {
					i2, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return fmt.Errorf("input %s must be a valid integer for task %s", key, t.Id)
					}

					i = i2
				} else if in.Default != nil {
					if v, ok := in.Default.(int64); ok {
						i = v
					}
				}
				ctx.State.Inputs.Set(key, i)
			case "int32":
				i := int32(0)
				if value != "" {
					i2, err := strconv.ParseInt(value, 10, 32)
					if err != nil {
						return fmt.Errorf("input %s must be a valid integer for task %s", key, t.Id)
					}

					i = int32(i2)
				} else if in.Default != nil {
					if v, ok := in.Default.(int32); ok {
						i = v
					}
				}
				ctx.State.Inputs.Set(key, i)
			case "uint":
				fallthrough
			case "uint64":
				i := uint64(0)
				if value != "" {
					i2, err := strconv.ParseUint(value, 10, 64)
					if err != nil {
						return fmt.Errorf("input %s must be a valid unsigned integer for task %s", key, t.Id)
					}

					i = i2
				} else if in.Default != nil {
					if v, ok := in.Default.(uint64); ok {
						i = v
					}
				}
				ctx.State.Inputs.Set(key, i)
			case "uint32":
				i := uint32(0)
				if value != "" {
					i2, err := strconv.ParseUint(value, 10, 32)
					if err != nil {
						return fmt.Errorf("input %s must be a valid unsigned integer for task %s", key, t.Id)
					}

					i = uint32(i2)
				} else if in.Default != nil {
					if v, ok := in.Default.(uint32); ok {
						i = v
					}
				}
				ctx.State.Inputs.Set(key, i)
			case "number":
				fallthrough
			case "float":
				fallthrough
			case "float64":
				f := float64(0)
				if value != "" {
					f2, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return fmt.Errorf("input %s must be a valid float for task %s", key, t.Id)
					}

					f = f2
				} else if in.Default != nil {
					if v, ok := in.Default.(float64); ok {
						f = v
					}
				}
				ctx.State.Inputs.Set(key, f)

			case "float32":
				f := float32(0)
				if value != "" {
					f2, err := strconv.ParseFloat(value, 32)
					if err != nil {
						return fmt.Errorf("input %s must be a valid float for task %s", key, t.Id)
					}

					f = float32(f2)
				} else if in.Default != nil {
					if v, ok := in.Default.(float32); ok {
						f = v
					}
				}
				ctx.State.Inputs.Set(key, f)
			case "bool":
				b := false
				if value != "" {
					b2, err := strconv.ParseBool(value)
					if err != nil {
						return fmt.Errorf("input %s must be a valid boolean for task %s", key, t.Id)
					}

					b = b2
				} else if in.Default != nil {
					if v, ok := in.Default.(bool); ok {
						b = v
					}
				}
				ctx.State.Inputs.Set(key, b)
			}

			for _, input := range ctx.Descriptor.Inputs {
				if input.IsRequired && !ctx.State.Inputs.Has(input.Name) {
					return fmt.Errorf("input %s is required for task %s", input.Name, t.Id)
				}
			}
		}
	}

	if t.Timeout != nil {
		if !t.Timeout.IsEvaluated {
			err := t.Timeout.Eval(ctx.Evaluator, data)
			if err != nil {
				return err
			}
		}

		ctx.State.Timeout = t.Timeout.Value.(uint32)
	}

	if t.Force != nil {
		if !t.Force.IsEvaluated {
			err := t.Force.Eval(ctx.Evaluator, data)
			if err != nil {
				return err
			}
		}

		ctx.State.Force = t.Force.Value.(bool)
	}

	if t.Cwd != nil {
		if !t.Cwd.IsEvaluated {
			err := t.Cwd.Eval(ctx.Evaluator, data)
			if err != nil {
				return err
			}
		}

		ctx.State.Cwd = t.Cwd.Value.(string)
	}

	if t.If != nil {
		if !t.If.IsEvaluated {
			err := t.If.Eval(ctx.Evaluator, data)
			if err != nil {
				return err
			}
		}

		ctx.State.If = t.If.Value.(bool)
	}

	if t.RunExpr != nil {
		if !t.RunExpr.IsEvaluated {
			err := t.RunExpr.Eval(ctx.Evaluator, data)
			if err != nil {
				return err
			}
		}

		ctx.State.RunExpr = t.RunExpr.Value.(string)
	}

	return nil
}

func mapOutputs(outputs *primitives.ObjectMap) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range outputs.Keys() {
		value := outputs.Get(key)
		if value == nil {
			continue
		}

		if v, ok := value.(map[string]interface{}); ok {
			result[key] = v
			continue
		}

		if v, ok := value.(primitives.ObjectMap); ok {
			result[key] = mapOutputs(&v)
			continue
		}

		result[key] = value
	}
	return result
}

func (t *Task) SetWith(inputs map[string]string) *Task {
	t.With = make(map[string]expr.Expression)
	for key, value := range inputs {
		t.With[key] = expr.Expression{
			Raw:         value,
			Value:       value,
			ValueString: value,
			IsEvaluated: true,
			Type:        "string",
		}
	}

	return t
}

func (t *Task) SetWithEntry(key, value string) *Task {
	if t.With == nil {
		t.With = make(map[string]expr.Expression)
	}

	e, ok := t.With[key]
	if !ok {
		e = expr.Expression{
			Raw:         value,
			Value:       value,
			ValueString: value,
			IsEvaluated: true,
			Type:        "string",
		}
		t.With[key] = e
	} else {
		e.Raw = value
		e.Value = value
		e.ValueString = value
		e.IsEvaluated = true
		t.With[key] = e
	}

	return t
}

func (t *Task) SetEnv(env map[string]string) *Task {
	t.Env = make(map[string]expr.Expression)
	for key, value := range env {
		t.Env[key] = expr.Expression{
			Raw:         value,
			Value:       value,
			ValueString: value,
			IsEvaluated: true,
			Type:        "string",
		}
	}

	return t
}

func (t *Task) SetEnvEntry(key, value string) *Task {
	if t.Env == nil {
		t.Env = make(map[string]expr.Expression)
	}

	e, ok := t.Env[key]
	if !ok {
		e = expr.Expression{
			Raw:         value,
			Value:       value,
			ValueString: value,
			IsEvaluated: true,
			Type:        "string",
		}
		t.Env[key] = e
	} else {
		e.Raw = value
		e.Value = value
		e.ValueString = value
		e.IsEvaluated = true
		t.Env[key] = e
	}

	return t
}

func (t *Task) SetTimeout(timeout uint32) *Task {
	str := strconv.FormatUint(uint64(timeout), 10)
	t.Timeout = &expr.Expression{
		Raw:         str,
		Value:       timeout,
		ValueString: str,
		IsEvaluated: true,
		Type:        "uint32",
	}

	return t
}

func (t *Task) SetForce(force bool) *Task {
	str := strconv.FormatBool(force)
	t.Force = &expr.Expression{
		Raw:         str,
		Value:       force,
		ValueString: str,
		IsEvaluated: true,
		Type:        "bool",
	}

	return t
}

func (t *Task) SetIf(condition bool) *Task {
	str := strconv.FormatBool(condition)
	t.If = &expr.Expression{
		Raw:         str,
		Value:       condition,
		ValueString: str,
		IsEvaluated: true,
		Type:        "bool",
	}

	return t
}

func (t *Task) SetCwd(cwd string) *Task {
	t.Cwd = &expr.Expression{
		Raw:         cwd,
		Value:       cwd,
		ValueString: cwd,
		IsEvaluated: true,
		Type:        "string",
	}

	return t
}

func (s *Task) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("task must be a mapping on line %d at column %d", node.Line, node.Column)
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		key := keyNode.Value

		switch key {
		case "id":
			s.Id = valueNode.Value
		case "name":
			s.Name = valueNode.Value
		case "uses":
			s.Uses = valueNode.Value
		case "description":
			s.Description = valueNode.Value
		case "needs":
			if valueNode.Kind != yaml.SequenceNode {
				return fmt.Errorf("needs must be a sequence on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			for _, n := range valueNode.Content {
				s.Needs = append(s.Needs, n.Value)
			}

		case "with":
			if valueNode.Kind != yaml.MappingNode {
				return fmt.Errorf("with must be a mapping on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			s.With = make(map[string]expr.Expression)
			for i := 0; i < len(valueNode.Content); i += 2 {
				kn := valueNode.Content[i]
				vn := valueNode.Content[i+1]

				target := expr.Expression{
					Raw:         vn.Value,
					IsEvaluated: false,
					Type:        "string",
				}

				if vn.Value == "" || !strings.Contains(vn.Value, "${{") {
					target.Value = vn.Value
					target.ValueString = vn.Value
					target.IsEvaluated = true
				}

				s.With[kn.Value] = target
			}

		case "env":
			if valueNode.Kind != yaml.MappingNode {
				return fmt.Errorf("env must be a mapping on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			s.Env = make(map[string]expr.Expression)
			for i := 0; i < len(valueNode.Content); i += 2 {
				kn := valueNode.Content[i]
				vn := valueNode.Content[i+1]

				target := expr.Expression{
					Raw:         vn.Value,
					IsEvaluated: false,
					Type:        "string",
				}

				if vn.Value == "" || !strings.Contains(vn.Value, "${{") {
					target.Value = vn.Value
					target.ValueString = vn.Value
					target.IsEvaluated = true
				}

				s.Env[kn.Value] = target
			}

		case "timeout":
			if valueNode.Kind != yaml.ScalarNode {
				return fmt.Errorf("timeout must be a scalar on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			s.Timeout = &expr.Expression{
				Raw:         valueNode.Value,
				IsEvaluated: false,
				Type:        "uint32",
			}

			actual := strings.TrimSpace(valueNode.Value)
			if actual == "" {
				s.Timeout.Value = 0
				s.Timeout.IsEvaluated = true
				s.Timeout.ValueString = "0"
				continue
			}

			if !strings.Contains(actual, "${{") {
				v, err := strconv.ParseUint(actual, 10, 32)
				if err != nil {
					return fmt.Errorf("timeout must be a valid unsigned integer on line %d at column %d", valueNode.Line, valueNode.Column)
				}

				s.Timeout.Value = uint32(v)
				s.Timeout.IsEvaluated = true
				s.Timeout.ValueString = actual
			}

		case "force":
			if valueNode.Kind != yaml.ScalarNode {
				return fmt.Errorf("force must be a scalar on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			s.Force = &expr.Expression{
				Raw:         valueNode.Value,
				IsEvaluated: false,
				Type:        "bool",
			}

			actual := strings.TrimSpace(valueNode.Value)
			if actual == "" {
				s.Force.Value = false
				s.Force.IsEvaluated = true
				s.Force.ValueString = "false"
				continue
			}

			if !strings.Contains(actual, "${{") {
				v, err := strconv.ParseBool(actual)
				if err != nil {
					return fmt.Errorf("force must be a valid boolean on line %d at column %d", valueNode.Line, valueNode.Column)
				}

				s.Force.Value = v
				s.Force.IsEvaluated = true
				s.Force.ValueString = actual
			}

		case "if":
			if valueNode.Kind != yaml.ScalarNode {
				return fmt.Errorf("if must be a scalar on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			s.If = &expr.Expression{
				Raw:         valueNode.Value,
				IsEvaluated: false,
				Type:        "bool",
			}

			actual := strings.TrimSpace(valueNode.Value)
			if actual == "" {
				s.If.Value = true
				s.If.IsEvaluated = true
				s.If.ValueString = "true"
				continue
			}

			if !strings.Contains(actual, "${{") {
				v, err := strconv.ParseBool(actual)
				if err != nil {
					return fmt.Errorf("if must be a valid boolean on line %d at column %d", valueNode.Line, valueNode.Column)
				}

				s.If.Value = v
				s.If.IsEvaluated = true
				s.If.ValueString = actual
			}

		case "cwd":
			if valueNode.Kind != yaml.ScalarNode {
				return fmt.Errorf("cwd must be a scalar on line %d at column %d", valueNode.Line, valueNode.Column)
			}

			s.Cwd = &expr.Expression{
				Raw:         valueNode.Value,
				IsEvaluated: false,
				Type:        "string",
			}

			actual := valueNode.Value
			if actual == "" {
				s.Cwd.Value = ""
				s.Cwd.IsEvaluated = true
				s.Cwd.ValueString = ""
				continue
			}

			if !strings.Contains(actual, "${{") {
				s.Cwd.Value = actual
				s.Cwd.IsEvaluated = true
				s.Cwd.ValueString = actual
			}

		default:
			return fmt.Errorf("unknown key %s on line %d at column %d", key, keyNode.Line, keyNode.Column)
		}
	}

	return nil
}
