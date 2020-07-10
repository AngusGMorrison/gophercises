package link

import (
	"os"
	"reflect"
	"testing"
)

const fixturePath = "fixtures/"

func TestParse(t *testing.T) {
	const testText = "Test text"
	tests := []struct {
		fixture       string
		wantNLinks    int
		wantFirstLink Link
	}{
		{
			"standard_a_tag.html",
			1,
			Link{"/test", testText},
		},
		{
			"nested_i_tag.html",
			1,
			Link{"https://www.test.com/test", testText},
		},
		{
			"nested_comment.html",
			1,
			Link{"/test", testText},
		},
		{
			"multiple_a_tags.html",
			3,
			Link{"#", testText},
		},
	}

	for _, test := range tests {
		t.Run(test.fixture, func(t *testing.T) {
			f, err := os.Open(fixturePath + test.fixture)
			if err != nil {
				t.Fatalf(err.Error())
			}
			links, err := Parse(f)
			if err != nil {
				t.Fatalf("parsing %s: %v", test.fixture, err)
			}
			if len(links) != test.wantNLinks {
				t.Fatalf("parsing %s: got %d links, want %d",
					test.fixture, len(links), test.wantNLinks)
			}
			if !reflect.DeepEqual(links[0], test.wantFirstLink) {
				t.Fatalf("parsing %s: first link returned was %+v, want %+v",
					test.fixture, links[0], test.wantFirstLink)
			}
		})
	}
}
