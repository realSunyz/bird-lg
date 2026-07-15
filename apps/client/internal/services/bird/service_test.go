package bird

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadResponseReturnsFinalCode(t *testing.T) {
	t.Parallel()

	reader := bufio.NewReader(strings.NewReader("1000-first line\nsecond line\n0000 completed\n"))
	response, code, err := readResponse(reader)
	if err != nil {
		t.Fatalf("readResponse() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("readResponse() code = %d, want 0", code)
	}
	if response != "first line\nsecond line\ncompleted\n" {
		t.Fatalf("readResponse() response = %q", response)
	}
}

func TestReadResponseParsesRestrictedAccess(t *testing.T) {
	t.Parallel()

	reader := bufio.NewReader(strings.NewReader("0016 Access restricted\n"))
	response, code, err := readResponse(reader)
	if err != nil {
		t.Fatalf("readResponse() error = %v", err)
	}
	if code != responseCodeAccessRestricted {
		t.Fatalf("readResponse() code = %d, want %d", code, responseCodeAccessRestricted)
	}
	if response != "Access restricted\n" {
		t.Fatalf("readResponse() response = %q", response)
	}
}
