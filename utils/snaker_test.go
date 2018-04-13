package utils

import "testing"

func TestToSnake(t *testing.T) {
	expected := "a_b_c"
	input := "a.b-c"
	got := ToSnake(input)

	if expected != got {
		t.Fatalf("ToSnake failed: input=%v, expected=%v, got=%v", input, expected, got)
	}

	expected = "a_new_world"
	input = "aNewWorld"
	got = ToSnake(input)

	if expected != got {
		t.Fatalf("ToSnake failed: input=%v, expected=%v, got=%v", input, expected, got)
	}
}
