package deque

import (
	"reflect"
	"sync"
	"testing"
)

func TestDeque_Front(t *testing.T) {
	type fields struct {
		first *node
		last  *node
		size  uint64
		mutex sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
		wantE  interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deque{
				first: tt.fields.first,
				last:  tt.fields.last,
				size:  tt.fields.size,
				mutex: tt.fields.mutex,
			}
			if gotE := d.Front(); !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("Front() = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}
