package paths

import (
	"fmt"
	"os"
	"testing"
)

func TestExpand(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("Cannot read $HOME")
	}

	tests := []struct {
		path        string
		expected    string
		description string
	}{
		{"aaa.txt", "aaa.txt", "Does not modify plain files"},
		{"/aaa/aaa.txt", "/aaa/aaa.txt", "Does not modify full paths"},
		{"~", home, "Expands tilde"},
		{"~/aaa.txt", fmt.Sprintf("%s/aaa.txt", home), "Expands tilde for files"},
		{"~/bbb/aaa.txt", fmt.Sprintf("%s/bbb/aaa.txt", home), "Expands tilde for files in hierarchy"},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := Path(test.path).expand()
			if actual != test.expected {
				t.Errorf("paths.Expand(%s) = %s, want %s", test.path, actual, test.expected)
			}
		})
	}
}
