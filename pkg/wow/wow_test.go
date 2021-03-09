package wow

import (
	"context"
	"testing"
)

func Test_InMemory(t *testing.T) {
	o := NewInMemory()
	for _, v := range data {
		w, err := o.GetNext(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got err: %v", err)
		}
		if v != w {
			t.Fatalf("expected result: %s, got: %s", v, w)
		}
	}
}
