package tasks

import (
	"time"

	"github.com/jolt9dev/go-jolt9/pkg/expr"
	"github.com/jolt9dev/go-jolt9/pkg/primitives"
)

type TaskState struct {
	Id          string
	Name        string
	Uses        string
	Description string
	Inputs      *primitives.ObjectMap
	Outputs     *primitives.ObjectMap
	Force       bool
	Timeout     uint32
	If          bool
	Env         map[string]string
	Cwd         string
	Needs       []string
	RunExpr     string
}

type TaskContext struct {
	primitives.Context
	State      *TaskState
	Descriptor *TaskDescriptor
	Evaluator  expr.Evaluator
}

type TaskResult struct {
	Id         string
	Outputs    *primitives.ObjectMap
	Status     int
	Error      error
	StartedAt  time.Time
	FinishedAt time.Time
}

func (t *TaskResult) SetError(err error) *TaskResult {
	t.Error = err
	t.Status = 10
	return t
}

func (t *TaskResult) SetOutputs(outputs *primitives.ObjectMap) *TaskResult {
	t.Outputs = outputs
	return t
}

func (t *TaskResult) SetStatus(status int) *TaskResult {
	t.Status = status
	return t
}

func (t *TaskResult) Start() *TaskResult {
	t.StartedAt = time.Now()
	return t
}

func (t *TaskResult) Finish() *TaskResult {
	t.FinishedAt = time.Now()
	t.Status = 1
	return t
}

func (t *TaskResult) Cancel() *TaskResult {
	t.FinishedAt = time.Now()
	t.Status = 5
	return t
}

func (t *TaskResult) Fail(err error) *TaskResult {
	t.FinishedAt = time.Now()
	t.Status = 10
	t.Error = err
	return t
}

func (t *TaskResult) Skip() *TaskResult {
	t.Status = 2
	return t
}

type DelegateTask interface {
	Run(ctx TaskContext) (primitives.ObjectMap, error)
}
