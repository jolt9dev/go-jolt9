package primitives

import "context"

type InputDescriptor struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string      `json:"type,omitempty" yaml:"type,omitempty"`
	IsRequired  bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	IsSecret    bool        `json:"secret,omitempty" yaml:"secret,omitempty"`
}

type OutputDescriptor struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	IsSecret    bool   `json:"secret,omitempty" yaml:"secret,omitempty"`
	IsRequired  bool   `json:"required,omitempty" yaml:"required,omitempty"`
}

type Expression interface {
	Type() string
	Raw() string
	IsEvaluated() bool
	Compiled() interface{}
	Error() error
}

type StringExpression interface {
	Expression
	String() string
}

type BoolExpression interface {
	Expression
	Bool() bool
}

type Int32Expression interface {
	Expression
	Int32() int32
}

type Int64Expression interface {
	Expression
	Int64() int64
}

type Float32Expression interface {
	Expression
	Float32() float32
}

type Float64Expression interface {
	Expression
	Float64() float64
}

type ObjectMap struct {
	Items map[string]interface{} `json:"items" yaml:"items"`
	order []string
}

type Message interface {
	Kind() string
}

type MessageSink interface {
	Send(msg Message) error
}

type MessageBus interface {
	Subscribe(sink MessageSink) error
	Unsubscribe(sink MessageSink) error
	Send(msg Message) error
}

type LoggingMessageBus interface {
	MessageBus
	Enabled(level int) bool
	SetLogLevel(level int)

	Tracef(format string, args ...interface{})
	TraceErrorf(err error, format string, args ...interface{})
	Debugf(format string, args ...interface{})
	DebugErrorf(err error, format string, args ...interface{})
	Infof(format string, args ...interface{})
	InfoErrorf(err error, format string, args ...interface{})
	Warnf(format string, args ...interface{})
	WarnErrorf(err error, format string, args ...interface{})
	Errorf(err error, format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	FatalErrorf(err error, format string, args ...interface{})
}

type Context struct {
	Signal   context.Context
	Env      map[string]string
	Vars     *ObjectMap
	Secrets  map[string]string
	Services map[string]interface{}
	Outputs  *ObjectMap
	Bus      LoggingMessageBus
}

func (o *ObjectMap) Add(key string, value interface{}) bool {
	if o.Items == nil {
		o.Items = make(map[string]interface{})
	}

	if _, ok := o.Items[key]; ok {
		return false
	}

	o.Items[key] = value
	o.order = append(o.order, key)
	return true
}

func (o *ObjectMap) Has(key string) bool {
	if o.Items == nil {
		return false
	}
	_, ok := o.Items[key]
	return ok
}

func (o *ObjectMap) Get(key string) interface{} {
	if o.Items == nil {
		return nil
	}
	return o.Items[key]
}

func (o *ObjectMap) Set(key string, value interface{}) {
	if o.Items == nil {
		o.Items = make(map[string]interface{})
	}

	if _, ok := o.Items[key]; !ok {
		o.order = append(o.order, key)
	}

	o.Items[key] = value
}

func (o *ObjectMap) Delete(key string) {
	if o.Items == nil {
		return
	}
	delete(o.Items, key)
}

func (o *ObjectMap) Keys() []string {
	return o.order
}

func (o *ObjectMap) Values() []interface{} {
	values := make([]interface{}, 0, len(o.Items))
	for _, key := range o.order {
		values = append(values, o.Items[key])
	}
	return values
}

func (o *ObjectMap) Len() int {
	return len(o.Items)
}

func (o *ObjectMap) Clear() {
	o.Items = make(map[string]interface{})
	o.order = nil
}

func (o *ObjectMap) GetString(key string) string {
	if o.Items == nil {
		return ""
	}
	if value, ok := o.Items[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func (o *ObjectMap) GetInt64(key string) int64 {
	if o.Items == nil {
		return 0
	}
	if value, ok := o.Items[key]; ok {
		if i, ok := value.(int64); ok {
			return i
		}
	}
	return 0
}

func (o *ObjectMap) GetInt32(key string) int32 {
	if o.Items == nil {
		return 0
	}
	if value, ok := o.Items[key]; ok {
		if i, ok := value.(int32); ok {
			return i
		}
	}
	return 0
}

func (o *ObjectMap) GetFloat64(key string) float64 {
	if o.Items == nil {
		return 0
	}
	if value, ok := o.Items[key]; ok {
		if f, ok := value.(float64); ok {
			return f
		}
	}
	return 0
}

func (o *ObjectMap) GetFloat32(key string) float32 {
	if o.Items == nil {
		return 0
	}

	if value, ok := o.Items[key]; ok {
		if f, ok := value.(float32); ok {
			return f
		}
	}

	return 0
}

func (o *ObjectMap) GetBool(key string) bool {
	if o.Items == nil {
		return false
	}
	if value, ok := o.Items[key]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}
