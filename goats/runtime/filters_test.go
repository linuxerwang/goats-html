package runtime

import (
	"testing"
)

type Struct0 struct {
	A string
	B []int
}

type Struct1 struct {
	X []*Struct0
	Y map[int]*Struct0
}

type DummyStruct struct {
	Name    string
	Titles  []string
	Prices  map[string]float32
	Embeded *Struct1
}

func TestBuiltinFilter_Debug(t *testing.T) {
	data := &DummyStruct{
		Name:   "Hello, world",
		Titles: []string{"A", "B", "C", "D"},
		Prices: map[string]float32{
			"A": 1.0,
			"B": 6.5,
			"C": 2.2,
			"D": 8.0,
		},
		Embeded: &Struct1{
			X: []*Struct0{
				&Struct0{
					A: "000",
					B: []int{10, 15},
				},
				&Struct0{
					A: "111",
					B: []int{100, 150},
				},
			},
			Y: map[int]*Struct0{
				1000: &Struct0{
					A: "T1000",
					B: []int{101, 102, 103},
				},
				1001: &Struct0{
					A: "T1001",
					B: []int{201, 202, 203},
				},
			},
		},
	}

	expected := `<div style="padding-left:20px">data (runtime.DummyStruct): {<br>` +
		`<div style="padding-left:20px">Name (string): Hello, world<br></div>` +
		`<div style="padding-left:20px">Titles ([]string): [A B C D]<br></div>` +
		`<div style="padding-left:20px">Prices (map[string]float32): map[A:1 B:6.5 C:2.2 D:8]<br></div>` +
		`<div style="padding-left:20px">Embeded (runtime.Struct1): {<br>` +
		`<div style="padding-left:20px">X ([]*runtime.Struct0): [0xc2000794e0 0xc200079510]<br></div>` +
		`<div style="padding-left:20px">Y (map[int]*runtime.Struct0): map[1000:0xc200079540 1001:0xc200079570]<br></div>` +
		`}<br></div>}<br></div>`

	filter := NewBuiltinFilter()
	output := filter.Debug(data)
	if output != expected {
		t.Errorf("Expected output not found: ", output)
	}
}

func TestTitle(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"this is a good example", "This Is A Good Example"},
		{"", ""},
	}
	filter := NewBuiltinFilter()
	for _, tt := range tests {
		output := filter.Title(tt.input)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
}

func TestCapfirst(t *testing.T) {
	filter := NewBuiltinFilter()
	var tests = []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello world"},
		{"hello 世界", "Hello 世界"},
		{"世界, hello", "世界, hello"},
		{"", ""},
	}
	for _, tt := range tests {
		output := filter.Capfirst(tt.input)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
	expected := filter.Capfirst(12)
	if expected != "<int Value>" {
		t.Errorf("Expected output: ", expected, " is wrong: ", "12")
	}
}

func TestCenter(t *testing.T) {
	filter := NewBuiltinFilter()
	var tests = []struct {
		input    string
		expected string
		width    int
	}{
		{"h", "  h  ", 5},
		{"h", "  h   ", 6},
		{"hello", "hello", 3},
	}
	for _, tt := range tests {
		output := filter.Center(tt.input, tt.width)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
}

func TestCut(t *testing.T) {
	filter := NewBuiltinFilter()
	var tests = []struct {
		input    string
		removed  string
		expected string
	}{
		{"hello", "l", "heo"},
		{"h", "h", ""},
	}
	for _, tt := range tests {
		output := filter.Cut(tt.input, tt.removed)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
}

func TestLjust(t *testing.T) {
	filter := NewBuiltinFilter()
	var tests = []struct {
		input    string
		width    int
		expected string
	}{
		{"hello", 4, "hello"},
		{"hello", 6, "hello "},
		{"hello 世界", 13, "hello 世界 "},
	}
	for _, tt := range tests {
		output := filter.Ljust(tt.input, tt.width)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
}

func TestRjust(t *testing.T) {
	filter := NewBuiltinFilter()
	var tests = []struct {
		input    string
		width    int
		expected string
	}{
		{"hello", 4, "hello"},
		{"hello", 6, " hello"},
		{"hello 世界", 13, " hello 世界"},
	}
	for _, tt := range tests {
		output := filter.Rjust(tt.input, tt.width)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
}

func TestFloatFormat(t *testing.T) {
	filter := NewBuiltinFilter()
	var tests = []struct {
		input    float64
		width    int
		expected string
	}{
		{34.23234, 1, "34.2"},
		{34.0000, 1, "34"},
		{34.2600, 1, "34.3"},
		{34.23234, 3, "34.232"},
		{34.0000, 3, "34.000"},
		{34.2600, 3, "34.260"},
		{34.23234, 0, "34"},
		{34.0000, 0, "34"},
		{34.2600, 0, "34"},
		{34.23234, -3, "34.232"},
		{34.0000, -3, "34"},
		{34.2600, -3, "34.260"},
	}
	for _, tt := range tests {
		output := filter.FloatFormat(tt.input, tt.width)
		if tt.expected != output {
			t.Errorf("Expected output: ", tt.expected, " is wrong: ", output)
		}
	}
}
