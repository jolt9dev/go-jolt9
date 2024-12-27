package tasks

import (
	"fmt"
	"slices"
)

type TaskMap struct {
	tasks map[string]*Task
	order []string
}

func (o *TaskMap) Add(key string, value *Task) bool {
	if o.tasks == nil {
		o.tasks = make(map[string]*Task)
	}

	if _, ok := o.tasks[key]; ok {
		return false
	}

	o.tasks[key] = value
	o.order = append(o.order, key)
	return true
}

func (o *TaskMap) Has(key string) bool {
	if o.tasks == nil {
		return false
	}
	_, ok := o.tasks[key]
	return ok
}

func (o *TaskMap) Get(key string) *Task {
	if o.tasks == nil {
		return nil
	}
	return o.tasks[key]
}

func (o *TaskMap) Set(key string, value Task) {
	if o.tasks == nil {
		o.tasks = make(map[string]*Task)
	}

	if _, ok := o.tasks[key]; !ok {
		o.order = append(o.order, key)
	}

	o.tasks[key] = &value
}

func (o *TaskMap) Delete(key string) {
	if o.tasks == nil {
		return
	}
	delete(o.tasks, key)

	index := slices.Index(o.order, key)
	if index >= 0 {
		o.order = append(o.order[:index], o.order[index+1:]...)
	}
}

func (o *TaskMap) Keys() []string {
	return o.order
}

func (o *TaskMap) Values() []Task {
	values := make([]Task, 0, len(o.tasks))
	for _, key := range o.order {
		values = append(values, *o.tasks[key])
	}
	return values
}

func (o *TaskMap) Flatten(targets []Task) ([]Task, error) {
	if len(targets) == 0 {
		targets = o.Values()
	}

	return flatten(*o, targets)
}

type MissingDependencyError struct {
	message string
	Tasks   *MissingDepResult
}

type MissingDepResult struct {
	Task    *Task
	Missing []string
}

func (e *MissingDependencyError) Error() string {
	return e.message
}

func (o *TaskMap) MissingDependencies() *MissingDependencyError {
	missing := &MissingDepResult{}
	for _, task := range o.tasks {
		for _, dep := range task.Needs {
			if !o.Has(dep) {
				if missing.Task == nil {
					missing.Task = task
				}

				missing.Missing = append(missing.Missing, dep)
			}
		}
	}

	if missing.Task != nil {
		return &MissingDependencyError{
			message: "missing dependencies",
			Tasks:   missing,
		}
	}

	return nil
}

func (o *TaskMap) Len() int {
	return len(o.order)
}

func (o *TaskMap) Clear() {
	o.tasks = nil
	o.order = nil
}

func (o *TaskMap) Copy() *TaskMap {
	var result TaskMap
	result.tasks = make(map[string]*Task)
	for key, value := range o.tasks {
		result.tasks[key] = value
	}
	result.order = make([]string, len(o.order))
	copy(result.order, o.order)
	return &result
}

func (o *TaskMap) At(index int) (string, Task, bool) {
	if index < 0 || index >= len(o.order) {
		var key string
		var value Task
		return key, value, false
	}
	key := o.order[index]
	value, ok := o.tasks[key]
	return key, *value, ok
}

func (o *TaskMap) FindCyclicalReferences() []Task {
	cycles := []Task{}
	stack := []*Task{}

	for _, task := range o.tasks {
		if !findCycle(task, &stack, *o) {
			cycles = append(cycles, *task)
		}
	}

	if len(cycles) == 0 {
		return nil
	}

	return cycles
}

func findCycle(task *Task, stack *[]*Task, tasks TaskMap) bool {
	index := -1
	s := *stack
	for i, t := range s {
		if t.Id == task.Id {
			index = i
			break
		}
	}

	if index >= 0 {
		return false
	}

	s = append(s, task)
	*stack = s
	for _, dep := range task.Needs {
		if !tasks.Has(dep) {
			continue
		}

		child := tasks.Get(dep)
		if child == nil {
			continue
		}

		result := findCycle(child, stack, tasks)
		if !result {
			return false
		}
	}

	s = s[:len(s)-1]
	*stack = s

	return true
}

func flatten(tasks TaskMap, set []Task) ([]Task, error) {
	results := []Task{}
	for _, task := range set {
		for _, dep := range task.Needs {
			if !tasks.Has(dep) {
				continue
			}

			child := tasks.Get(dep)
			if child == nil {
				return nil, fmt.Errorf("task %s depends on missing task %s", task.Id, dep)
			}

			childResult, err := flatten(tasks, []Task{*child})
			if err != nil {
				return nil, err
			}

			results = append(results, childResult...)
			found := false
			for _, next := range results {
				if next.Id == child.Id {
					found = true
					break
				}
			}

			if !found {
				results = append(results, *child)
			}
		}

		found := false
		for _, next := range results {
			if next.Id == task.Id {
				found = true
				break
			}
		}

		if !found {
			results = append(results, task)
		}
	}

	return results, nil
}
