package deque

import (
	"testing"
)

func intCompare(a, b int) bool {
	return a == b
}

func TestNewArrayDeque(t *testing.T) {
	d := NewArrayDeque(intCompare)
	if d.Capacity() != DefaultCapacity {
		t.Errorf("Expected capacity %d, got %d", DefaultCapacity, d.Capacity())
	}
	if !d.IsEmpty() {
		t.Error("New deque should be empty")
	}
}

func TestAddFirst(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddFirst(1)
	d.AddFirst(2)
	if d.Size() != 2 {
		t.Errorf("Expected size 2, got %d", d.Size())
	}
	first, err := d.GetFirst()
	if err != nil || first != 2 {
		t.Errorf("Expected first element 2, got %d", first)
	}
}

func TestAddLast(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	if d.Size() != 2 {
		t.Errorf("Expected size 2, got %d", d.Size())
	}
	last, err := d.GetLast()
	if err != nil || last != 2 {
		t.Errorf("Expected last element 2, got %d", last)
	}
}

func TestRemoveFirst(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	removed, err := d.RemoveFirst()
	if err != nil || removed != 1 {
		t.Errorf("Expected removed element 1, got %d", removed)
	}
	if d.Size() != 1 {
		t.Errorf("Expected size 1, got %d", d.Size())
	}
}

func TestRemoveLast(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	removed, err := d.RemoveLast()
	if err != nil || removed != 2 {
		t.Errorf("Expected removed element 2, got %d", removed)
	}
	if d.Size() != 1 {
		t.Errorf("Expected size 1, got %d", d.Size())
	}
}

func TestGet(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	val, err := d.Get(1)
	if err != nil || val != 2 {
		t.Errorf("Expected element 2 at index 1, got %d", val)
	}
}

func TestSet(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	oldVal, err := d.Set(1, 4)
	if err != nil || oldVal != 2 {
		t.Errorf("Expected old value 2, got %d", oldVal)
	}
	newVal, _ := d.Get(1)
	if newVal != 4 {
		t.Errorf("Expected new value 4 at index 1, got %d", newVal)
	}
}

func TestContains(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	if !d.Contains(2) {
		t.Error("Deque should contain 2")
	}
	if d.Contains(4) {
		t.Error("Deque should not contain 4")
	}
}

func TestIndexOf(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	d.AddLast(2)
	if d.IndexOf(2) != 1 {
		t.Errorf("Expected index 1 for first occurrence of 2, got %d", d.IndexOf(2))
	}
	if d.IndexOf(4) != -1 {
		t.Errorf("Expected index -1 for non-existent element, got %d", d.IndexOf(4))
	}
}

func TestLastIndexOf(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	d.AddLast(2)
	if d.LastIndexOf(2) != 3 {
		t.Errorf("Expected index 3 for last occurrence of 2, got %d", d.LastIndexOf(2))
	}
	if d.LastIndexOf(4) != -1 {
		t.Errorf("Expected index -1 for non-existent element, got %d", d.LastIndexOf(4))
	}
}

func TestClear(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.Clear()
	if !d.IsEmpty() {
		t.Error("Deque should be empty after clear")
	}
}

func TestToArray(t *testing.T) {
	d := NewArrayDeque(intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	arr := d.ToArray()
	if len(arr) != 3 || arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
		t.Errorf("ToArray returned unexpected result: %v", arr)
	}
}

func TestEnsureCapacity(t *testing.T) {
	d := NewArrayDeque(intCompare)
	initialCap := d.Capacity()
	d.EnsureCapacity(initialCap * 2)
	if d.Capacity() <= initialCap {
		t.Errorf("Capacity should have increased, got %d", d.Capacity())
	}
}

func TestTrimToSize(t *testing.T) {
	d := NewArrayDequeWithCapacity(100, intCompare)
	for i := 0; i < 10; i++ {
		d.AddLast(i)
	}
	d.TrimToSize()
	if d.Capacity() != 10 {
		t.Errorf("Expected capacity 10 after trim, got %d", d.Capacity())
	}
}

func TestWraparound(t *testing.T) {
	d := NewArrayDequeWithCapacity(4, intCompare)
	d.AddLast(1)
	d.AddLast(2)
	d.AddLast(3)
	d.AddLast(4)
	d.RemoveFirst()
	d.RemoveFirst()
	d.AddLast(5)
	d.AddLast(6)
	expected := []int{3, 4, 5, 6}
	actual := d.ToArray()
	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("Expected %v, got %v", expected, actual)
			break
		}
	}
}
