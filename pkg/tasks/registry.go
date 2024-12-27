package tasks

import "github.com/jolt9dev/go-jolt9/pkg/primitives"

type TaskDescriptor struct {
	Id          string                                 `json:"id" yaml:"id"`
	Version     string                                 `json:"version,omitempty" yaml:"version,omitempty"`
	Description string                                 `json:"description,omitempty" yaml:"description,omitempty"`
	Inputs      map[string]primitives.InputDescriptor  `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Outputs     map[string]primitives.OutputDescriptor `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Redirect    bool                                   `json:"redirect,omitempty" yaml:"redirect,omitempty"`
	RunFile     string                                 `json:"run,omitempty" yaml:"run,omitempty"`
	Uses        string                                 `json:"uses,omitempty" yaml:"uses,omitempty"`
}

type TaskRegistry struct {
	tasks map[string]*TaskDescriptor
}

func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		tasks: make(map[string]*TaskDescriptor),
	}
}

func (r *TaskRegistry) Register(task *TaskDescriptor) {
	r.tasks[task.Id] = task
}

func (r *TaskRegistry) Get(id string) (*TaskDescriptor, bool) {
	task, ok := r.tasks[id]
	return task, ok
}
