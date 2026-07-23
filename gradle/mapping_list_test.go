package gradle

import (
	"reflect"
	"testing"
)

func TestDecodeMappingList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: nil},
		{name: "whitespace only", input: "  \n ", want: nil},
		{name: "pipe separated", input: "a.txt|b.txt", want: []string{"a.txt", "b.txt"}},
		{name: "preserves internal empty", input: "a.txt||c.txt", want: []string{"a.txt", "", "c.txt"}},
		{name: "preserves trailing empty", input: "a.txt|", want: []string{"a.txt", ""}},
		{name: "preserves leading empty", input: "|b.txt", want: []string{"", "b.txt"}},
		{name: "newline separated", input: "a.txt\nb.txt", want: []string{"a.txt", "b.txt"}},
		{name: "literal backslash-n separated", input: `a.txt\nb.txt`, want: []string{"a.txt", "b.txt"}},
		{name: "trims each entry", input: " a.txt | b.txt ", want: []string{"a.txt", "b.txt"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeMappingList(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeMappingList(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEncodeDecodeMappingListRoundTrip(t *testing.T) {
	cases := [][]string{
		{"a.txt", "b.txt"},
		{"a.txt", "", "c.txt"},
		{"a.txt", ""},
		{"", "b.txt"},
		{"only.txt"},
	}

	for _, paths := range cases {
		encoded := EncodeMappingList(paths)
		decoded := DecodeMappingList(encoded)
		if !reflect.DeepEqual(decoded, paths) {
			t.Errorf("round trip of %#v: encoded %q, decoded %#v", paths, encoded, decoded)
		}
	}
}
