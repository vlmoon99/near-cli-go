package main

import "testing"

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"message", "Message"},
		{"newMessage", "NewMessage"},
		{"a", "A"},
		{"", ""},
		{"AlreadyCap", "AlreadyCap"},
		{"user_id", "User_id"},
	}

	for _, tt := range tests {
		result := capitalizeFirst(tt.input)
		if result != tt.expected {
			t.Errorf("capitalizeFirst(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"Submit", "submit"},
		{"PDFLoad", "p_d_f_load"},
		{"myFunction", "my_function"},
	}

	for _, tt := range tests {
		result := toSnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("toSnakeCase(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
