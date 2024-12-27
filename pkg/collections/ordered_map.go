package collections

import "slices"

type OrderedMap[V any] struct {
	items map[string]V
	order []string
}

func (o *OrderedMap[V]) Add(key string, value V) bool {
	if o.items == nil {
		o.items = make(map[string]V)
	}

	if _, ok := o.items[key]; ok {
		return false
	}

	o.items[key] = value
	o.order = append(o.order, key)
	return true
}

func (o *OrderedMap[V]) Keys() []string {
	return o.order
}

func (o *OrderedMap[V]) Values() []V {
	var result []V
	for _, key := range o.order {
		result = append(result, o.items[key])
	}
	return result
}

func (o *OrderedMap[V]) Len() int {
	return len(o.order)
}

func (o *OrderedMap[V]) Clear() {
	o.items = nil
	o.order = nil
}

func (o *OrderedMap[V]) Copy() *OrderedMap[V] {
	var result OrderedMap[V]
	result.items = make(map[string]V)
	for key, value := range o.items {
		result.items[key] = value
	}
	result.order = make([]string, len(o.order))
	copy(result.order, o.order)
	return &result
}

func (o *OrderedMap[V]) At(index int) (string, V, bool) {
	if index < 0 || index >= len(o.order) {
		var key string
		var value V
		return key, value, false
	}
	key := o.order[index]
	value, ok := o.items[key]
	return key, value, ok
}

func (o *OrderedMap[V]) Has(key string) bool {
	if o.items == nil {
		return false
	}
	_, ok := o.items[key]
	return ok
}

func (o *OrderedMap[V]) Get(key string) V {
	if o.items == nil {
		var result V
		return result
	}
	return o.items[key]
}

func (o *OrderedMap[V]) Set(key string, value V) {
	if o.items == nil {
		o.items = make(map[string]V)
	}

	if _, ok := o.items[key]; !ok {
		o.order = append(o.order, key)
	}

	o.items[key] = value
}

func (o *OrderedMap[V]) Delete(key string) {
	if o.items == nil {
		return
	}

	if _, ok := o.items[key]; !ok {
		return
	}

	delete(o.items, key)
	index := slices.Index(o.order, key)
	if index != -1 {
		o.order = append(o.order[:index], o.order[index+1:]...)
	}
}

func (o *OrderedMap[V]) ToMap() map[string]V {
	return o.items
}
