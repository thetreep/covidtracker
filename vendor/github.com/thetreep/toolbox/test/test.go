package test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func Compare(t *testing.T, got, want interface{}, reasons ...string) {
	var reason string
	if len(reasons) > 0 {
		reason = reasons[0] + ": "
	}
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		var out bytes.Buffer
		diff := pretty.Diff(got, want)
		out.WriteString(fmt.Sprintf("%sgot -> want: %d lines different:\n", reason, len(diff)))
		out.WriteString("\t" + strings.Join(diff, "\n \t"))
		t.Fatal(out.String())
	}
}
