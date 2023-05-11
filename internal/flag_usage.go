package internal

import (
	"fmt"
	"github.com/spf13/pflag"
)

// FlagUsage extracts logic from pflag.FlagSet#FlagUsages() for consistent line-based usage text.
// Some extension has been made to account for unexported functionality.
//
// Note that this differs from the FlagUsages() original in that it doesn't account for widths of
// other sibling flags.
//
// pflag is BSD-3 licensed:
//
//	Copyright (c) 2012 Alex Ogier. All rights reserved.
//	Copyright (c) 2012 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// See https://github.com/spf13/pflag/blob/master/LICENSE
func FlagUsage(flag *pflag.Flag) string {
	if flag.Hidden {
		return ""
	}

	line := ""
	if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
		line = fmt.Sprintf("  -%s, --%s", flag.Shorthand, flag.Name)
	} else {
		line = fmt.Sprintf("      --%s", flag.Name)
	}

	varname, usage := pflag.UnquoteUsage(flag)
	if varname != "" {
		line += " " + varname
	}
	if flag.NoOptDefVal != "" {
		switch flag.Value.Type() {
		case "string":
			line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
		case "bool":
			if flag.NoOptDefVal != "true" {
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		case "count":
			if flag.NoOptDefVal != "+1" {
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		default:
			line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
		}
	}

	line += fmt.Sprintf("\t%s", usage)
	defaultValue := defaultValueForFlag(flag)

	if defaultValue != "" {
		if flag.Value.Type() == "string" {
			line += fmt.Sprintf(" (default %q)", flag.DefValue)
		} else {
			line += fmt.Sprintf(" (default %s)", flag.DefValue)
		}
	}
	if len(flag.Deprecated) != 0 {
		line += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
	}

	return line
}

func defaultValueForFlag(flag *pflag.Flag) string {
	var defaultValue string
	switch flag.Value.Type() {
	case "bool":
		if flag.DefValue != "false" {
			defaultValue = flag.DefValue
		}
	case "duration":
		if flag.DefValue != "0" && flag.DefValue != "0s" {
			defaultValue = flag.DefValue
		}
	case "int", "int8", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "count", "float32", "float64":
		if flag.DefValue != "0" {
			defaultValue = flag.DefValue
		}
	case "string":
		if flag.DefValue != "" {
			defaultValue = flag.DefValue
		}
	case "ip", "ipMask", "ipNet":
		if flag.DefValue != "<nil>" {
			defaultValue = flag.DefValue
		}
	case "intSlice", "stringSlice", "stringArrayValue":
		if flag.DefValue != "[]" {
			defaultValue = flag.DefValue
		}
	default:
		switch flag.Value.String() {
		case "false", "<nil>", "", "0":
			break
		default:
			defaultValue = flag.DefValue
		}
	}
	return defaultValue
}
