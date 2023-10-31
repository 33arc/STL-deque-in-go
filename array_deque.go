package deque

import (
	"errors"
	"fmt"
)

const (
	GrowthPolicyRelative = iota
	GrowthPolicyAbsolute

	DefaultGrowthPolicy = GrowthPolicyRelative
	DefaultGrowthFactor = 1.0
	DefaultGrowthChunk  = 10
	DefaultCapacity     = 10
)

type ArrayDeque[T any] struct {
	data         []T
	size         int
	first        int
	last         int
	growthPolicy int
	growthFactor float64
	growthChunk  int
	compare      func(T, T) bool
}

func NewArrayDeque[T any](compare func(T, T) bool) *ArrayDeque[T] {
	return NewArrayDequeWithCapacity[T](DefaultCapacity, compare)
}

func NewArrayDequeWithCapacity[T any](capacity int, compare func(T, T) bool) *ArrayDeque[T] {
	return NewArrayDequeWithCapacityAndGrowth[T](capacity, DefaultGrowthFactor, compare)
}

func NewArrayDequeWithCapacityAndGrowth[T any](capacity int, growthFactor float64, compare func(T, T) bool) *ArrayDeque[T] {
	return NewArrayDequeWithParams[T](capacity, GrowthPolicyRelative, growthFactor, DefaultGrowthChunk, compare)
}

func NewArrayDequeWithCapacityAndChunk[T any](capacity, growthChunk int, compare func(T, T) bool) *ArrayDeque[T] {
	return NewArrayDequeWithParams[T](capacity, GrowthPolicyAbsolute, DefaultGrowthFactor, growthChunk, compare)
}

func NewArrayDequeWithParams[T any](capacity, growthPolicy int, growthFactor float64, growthChunk int, compare func(T, T) bool) *ArrayDeque[T] {
	if capacity < 0 {
		panic(fmt.Sprintf("Negative capacity: %d", capacity))
	}
	if growthFactor < 0.0 {
		panic(fmt.Sprintf("Negative growth factor: %f", growthFactor))
	}
	if growthChunk < 0 {
		panic(fmt.Sprintf("Negative growth chunk: %d", growthChunk))
	}
	return &ArrayDeque[T]{
		data:         make([]T, capacity),
		size:         0,
		first:        0,
		last:         0,
		growthPolicy: growthPolicy,
		growthFactor: growthFactor,
		growthChunk:  growthChunk,
		compare:      compare,
	}
}

func (d *ArrayDeque[T]) computeCapacity(capacity int) int {
	newCapacity := 0
	if d.growthPolicy == GrowthPolicyRelative {
		newCapacity = int(float64(len(d.data)) * (1.0 + d.growthFactor))
	} else {
		newCapacity = len(d.data) + d.growthChunk
	}
	if newCapacity < capacity {
		newCapacity = capacity
	}
	return newCapacity
}

func (d *ArrayDeque[T]) EnsureCapacity(capacity int) int {
	if capacity > len(d.data) {
		newCapacity := d.computeCapacity(capacity)
		newData := make([]T, newCapacity)
		d.toArray(newData)
		d.first = 0
		d.last = d.size
		d.data = newData
		return newCapacity
	}
	return len(d.data)
}

func (d *ArrayDeque[T]) Capacity() int {
	return len(d.data)
}

func (d *ArrayDeque[T]) Add(index int, v T) {
	if index == 0 {
		d.AddFirst(v)
	} else if index == d.size {
		d.AddLast(v)
	} else {
		if index < 0 || index > d.size {
			panic(fmt.Sprintf("Index out of bounds: %d", index))
		}
		d.EnsureCapacity(d.size + 1)
		if d.first < d.last || d.last == 0 {
			iidx := index + d.first
			end := d.last
			if d.last == 0 {
				end = len(d.data)
			}
			block := end - iidx
			if d.last == 0 {
				d.data[0] = d.data[len(d.data)-1]
				copy(d.data[iidx+1:], d.data[iidx:iidx+block-1])
				d.last = 1
			} else {
				copy(d.data[iidx+1:], d.data[iidx:iidx+block])
				d.last++
				if d.last == len(d.data) {
					d.last = 0
				}
			}
			d.data[iidx] = v
		} else {
			iidx := (d.first + index) % len(d.data)
			if iidx <= d.last {
				block := d.last - iidx
				copy(d.data[iidx+1:], d.data[iidx:iidx+block])
				d.last++
				d.data[iidx] = v
			} else {
				block := iidx - d.first
				copy(d.data[d.first-1:], d.data[d.first:d.first+block])
				d.first--
				d.data[iidx-1] = v
			}
		}
		d.size++
	}
}

func (d *ArrayDeque[T]) Get(index int) (T, error) {
	if index < 0 || index >= d.size {
		var zero T
		return zero, fmt.Errorf("Index out of bounds: %d", index)
	}
	return d.data[(d.first+index)%len(d.data)], nil
}

func (d *ArrayDeque[T]) Set(index int, v T) (T, error) {
	if index < 0 || index >= d.size {
		var zero T
		return zero, fmt.Errorf("Index out of bounds: %d", index)
	}
	idx := (d.first + index) % len(d.data)
	result := d.data[idx]
	d.data[idx] = v
	return result, nil
}

func (d *ArrayDeque[T]) RemoveElementAt(index int) (T, error) {
	if index == 0 {
		return d.RemoveFirst()
	} else if index == d.size-1 {
		return d.RemoveLast()
	} else {
		if index < 0 || index >= d.size {
			var zero T
			return zero, fmt.Errorf("Index out of bounds: %d", index)
		}
		ridx := (d.first + index) % len(d.data)
		result := d.data[ridx]
		if d.first < d.last || d.last == 0 {
			block1 := ridx - d.first
			block2 := d.size - block1 - 1
			if block1 < block2 {
				copy(d.data[d.first+1:], d.data[d.first:d.first+block1])
				d.first++
			} else {
				copy(d.data[ridx:], d.data[ridx+1:ridx+1+block2])
				d.last--
				if d.last < 0 {
					d.last = len(d.data) - 1
				}
			}
		} else {
			if ridx < d.last {
				block := d.last - ridx - 1
				copy(d.data[ridx:], d.data[ridx+1:ridx+1+block])
				d.last--
			} else {
				block := ridx - d.first
				copy(d.data[d.first+1:], d.data[d.first:d.first+block])
				d.first++
				if d.first == len(d.data) {
					d.first = 0
				}
			}
		}
		d.size--
		return result, nil
	}
}

func (d *ArrayDeque[T]) TrimToSize() {
	if len(d.data) > d.size {
		newData := d.ToArray()
		d.first = 0
		d.last = 0
		d.data = newData
	}
}

func (d *ArrayDeque[T]) GetFirst() (T, error) {
	if d.size == 0 {
		var zero T
		return zero, errors.New("Deque is empty")
	}
	return d.data[d.first], nil
}

func (d *ArrayDeque[T]) GetLast() (T, error) {
	if d.size == 0 {
		var zero T
		return zero, errors.New("Deque is empty")
	}
	if d.last == 0 {
		return d.data[len(d.data)-1], nil
	}
	return d.data[d.last-1], nil
}

func (d *ArrayDeque[T]) AddFirst(v T) {
	d.EnsureCapacity(d.size + 1)
	d.first--
	if d.first < 0 {
		d.first = len(d.data) - 1
	}
	d.data[d.first] = v
	d.size++
}

func (d *ArrayDeque[T]) AddLast(v T) {
	d.EnsureCapacity(d.size + 1)
	d.data[d.last] = v
	d.last++
	if d.last == len(d.data) {
		d.last = 0
	}
	d.size++
}

func (d *ArrayDeque[T]) RemoveFirst() (T, error) {
	if d.size == 0 {
		var zero T
		return zero, errors.New("Deque is empty")
	}
	result := d.data[d.first]
	d.first++
	if d.first == len(d.data) {
		d.first = 0
	}
	d.size--
	return result, nil
}

func (d *ArrayDeque[T]) RemoveLast() (T, error) {
	if d.size == 0 {
		var zero T
		return zero, errors.New("Deque is empty")
	}
	d.last--
	if d.last < 0 {
		d.last = len(d.data) - 1
	}
	d.size--
	return d.data[d.last], nil
}

func (d *ArrayDeque[T]) Size() int {
	return d.size
}

func (d *ArrayDeque[T]) IsEmpty() bool {
	return d.size == 0
}

func (d *ArrayDeque[T]) Clear() {
	d.size = 0
	d.first = 0
	d.last = 0
}

func (d *ArrayDeque[T]) Contains(v T) bool {
	for i, idx := 0, d.first; i < d.size; i++ {
		if d.compare(d.data[idx], v) {
			return true
		}
		idx++
		if idx == len(d.data) {
			idx = 0
		}
	}
	return false
}

func (d *ArrayDeque[T]) IndexOf(v T) int {
	for i, idx := 0, d.first; i < d.size; i++ {
		if d.compare(d.data[idx], v) {
			return i
		}
		idx++
		if idx == len(d.data) {
			idx = 0
		}
	}
	return -1
}

func (d *ArrayDeque[T]) LastIndexOf(v T) int {
	for i, idx := d.size-1, d.last-1; i >= 0; i-- {
		if idx < 0 {
			idx = len(d.data) - 1
		}
		if d.compare(d.data[idx], v) {
			return i
		}
		idx--
	}
	return -1
}

func (d *ArrayDeque[T]) Remove(v T) bool {
	index := d.IndexOf(v)
	if index != -1 {
		_, err := d.RemoveElementAt(index)
		return err == nil
	}
	return false
}

func (d *ArrayDeque[T]) ToArray() []T {
	a := make([]T, d.size)
	d.toArray(a)
	return a
}

func (d *ArrayDeque[T]) toArray(a []T) {
	if d.last <= d.first {
		if d.last == 0 {
			copy(a, d.data[d.first:d.first+d.size])
		} else {
			block1 := len(d.data) - d.first
			block2 := d.size - block1
			copy(a, d.data[d.first:])
			copy(a[block1:], d.data[:block2])
		}
	} else {
		copy(a, d.data[d.first:d.first+d.size])
	}
}
