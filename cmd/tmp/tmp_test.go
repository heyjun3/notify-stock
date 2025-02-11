package main

import (
	"testing"
)

type User struct {
	id   int
	name string
}

func TestStructMap(t *testing.T) {
	m := map[User]int{
		{id: 10, name: "mike"}: 10,
	}

	if num, exp := m[User{id: 10, name: "mike"}], 10; num != exp {
		t.Errorf("not asserted numbers, actual: %v, expected: %v", num, exp)
	}
	if num, exp := m[User{name: "mike", id: 10}], 10; num != exp {
		t.Errorf("not asserted numbers, actual: %v, expected: %v", num, exp)
	}
	if num, exp := m[User{name: "mike", id: 11}], 10; num == exp {
		t.Errorf("not asserted numbers, actual: %v, expected: %v", num, exp)
	}
	if num, exp := m[User{}], 10; num == exp {
		t.Errorf("not asserted numbers, actual: %v, expected: %v", num, exp)
	}
}
