package venom

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

type testFunctions struct{}

func (t testFunctions) FormatHeader(input string) string {
	return "FormatHeader"
}

func (t testFunctions) FormatText(input string) string {
	return "FormatText"
}

func (t testFunctions) FormatOptions(input string) string {
	return "FormatOptions"
}

func (t testFunctions) FormatFlag(input Flag) string {
	return "FormatFlag"
}

func (t testFunctions) SeeAlsoPath(input string) string {
	return "SeeAlsoPath"
}

func (t testFunctions) FormatExample(input string) string {
	return "FormatExample"
}

func (t testFunctions) FormatAutoGenTag(input string) string {
	return "FormatAutoGenTag"
}

func (t testFunctions) IsLocalFlag(input Flag) bool {
	return true
}

func Test_newFuncMap(t *testing.T) {
	// just check that mappings are done appropriately as a user facing contract (no surprises)
	f := testFunctions{}
	funcMap := newFuncMap(f)
	tests := []struct {
		key    string
		expect string
	}{
		{key: "header", expect: "FormatHeader"},
		{key: "text", expect: "FormatText"},
		{key: "options", expect: "FormatOptions"},
		{key: "flag", expect: "FormatFlag"},
		{key: "see_also_path", expect: "SeeAlsoPath"},
		{key: "example", expect: "FormatExample"},
		{key: "autogen", expect: "FormatAutoGenTag"},
		{key: "is_local", expect: "IsLocalFlag"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("testFunctions '%s' invokes %s", tt.key, tt.expect)
		t.Run(name, func(t *testing.T) {
			got := funcMap[tt.key]
			v := reflect.ValueOf(got)
			if v.Kind() != reflect.Func {
				t.Fatalf(tt.key + " not a function")
			}
			fullFuncName := runtime.FuncForPC(v.Pointer()).Name()

			// We only care that the proper function ends up bound to the expected template key, so check the reflected name
			pattern := regexp.MustCompile(`.*\.([^-]+)-?.*`)
			matches := pattern.FindStringSubmatch(fullFuncName)
			actual := matches[1]
			if actual != tt.expect {
				t.Errorf("newFuncMap() = %v, want %v", actual, tt.expect)
			}
		})
	}
}
