package librarian

import (
	"fmt"
	"testing"
)

type tt struct {
	One int
	Two string
}

func (o tt) Exports() string {
	return fmt.Sprintf("One:%d, Two:%s", o.One, o.Two)
}

func TestScalar(t *testing.T) {
	Set("a", "A")
	if actual := Has("a"); actual != true {
		t.Errorf("expected `Has` return %t got %t", true, actual)
	}
	if actual := IsResolved("a"); actual != true {
		t.Errorf("expected `IsResolved` return %t got %t", true, actual)
	}
	if actual := IsUnresolved("a"); actual != false {
		t.Errorf("expected `IsUnresolved` return %t got %t", false, actual)
	}
	Remove("a")
	if actual := Has("a"); actual != false {
		t.Errorf("expected `a` return %t got %t", false, actual)
	}
	if actual := Get("a"); actual != nil {
		t.Errorf("expected returned %s got %s", nil, actual)
	}

	Set("b", 1)
	if actual := Has("b"); actual != true {
		t.Errorf("expected `Has` return %t got %t", true, actual)
	}
	if actual := IsResolved("b"); actual != true {
		t.Errorf("expected `IsResolved` return %t got %t", true, actual)
	}
	if actual := IsUnresolved("b"); actual != false {
		t.Errorf("expected `IsUnresolved` return %t got %t", false, actual)
	}
	if actual := Get("b"); actual != 1 {
		t.Errorf("expected returned %s got %s", 1, actual)
	}
	Remove("b")
	if actual := Has("b"); actual != false {
		t.Errorf("expected `Has` return %t got %t", false, actual)
	}
	if actual := Get("b"); actual != nil {
		t.Errorf("expected returned %s got %s", nil, actual)
	}
}

func TestInstance(t *testing.T) {

	k := "tinstance"
	v := tt{
		One: 1,
		Two: "o3",
	}

	Set(k, v)

	if actual := IsResolved(k); actual != true {
		t.Errorf("expected `%s` return %t got %t", k, true, actual)
	}

	if actual := IsUnresolved(k); actual != false {
		t.Errorf("expected `%s` return %t got %t", k, false, actual)
	}

	if actual := Get(k).(tt); actual != v {
		t.Errorf("expected returned %T got %T", v, actual)
	}

	Remove(k)

	if actual := IsUnresolved(k); actual != false {
		t.Errorf("expected `%s` return %t got %t", k, false, actual)
	}

	if actual := IsResolved(k); actual != false {
		t.Errorf("expected `%s` return %t got %t", k, false, actual)
	}

	if actual := Get(k); actual != nil {
		t.Errorf("expected returned nil got %T", actual)
	}

}

func TestPointer(t *testing.T) {

	k := "tpointer"
	v := tt{
		One: 1,
		Two: "o3",
	}

	Set("tpointer", func() interface{} {
		return &v
	})

	if actual := IsUnresolved(k); actual != true {
		t.Errorf("expected `%s` return %t got %t", k, true, actual)
	}

	if actual := IsResolved(k); actual != false {
		t.Errorf("expected `%s` return %t got %t", k, false, actual)
	}

	actual := Get(k).(*tt)
	if actual.Exports() != v.Exports() {
		t.Errorf("expected returned %T got %T", v.Exports(), actual.Exports())
	}

	if actual := IsResolved(k); actual != true {
		t.Errorf("expected `%s` return %t got %t", k, true, actual)
	}

	Remove(k)

	if actual := IsUnresolved(k); actual != false {
		t.Errorf("expected `%s` return %t got %t", k, false, actual)
	}

	if actual := IsResolved(k); actual != false {
		t.Errorf("expected `%s` return %t got %t", k, false, actual)
	}

	if actual := Get(k); actual != nil {
		t.Errorf("expected returned nil got %T", actual)
	}
}
