/*
==========
Cariddi
==========

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.

	@Repository:  https://github.com/ocervell/cariddi

	@Author:      edoardottt, https://www.edoardoottavianelli.it

	@License: https://github.com/ocervell/cariddi/blob/main/LICENSE

*/

package url_test

import (
	"reflect"
	"testing"

	urlUtils "github.com/ocervell/cariddi/internal/url"
)

func TestGetHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "no protocol",
			input: "edoardottt.com/ciao?id=1",
			want:  "edoardottt.com",
		},
		{
			name:  "ok1",
			input: "http://edoardottt.com/ciao?id=1",
			want:  "edoardottt.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlUtils.GetHost(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProtocol(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "no protocol",
			input: "ciao.com",
			want:  "",
		},
		{
			name:  "ok1",
			input: "http://edoardottt.com/ciao?id=1",
			want:  "http",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlUtils.GetProtocol(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRootHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "no protocol",
			input: "ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "no protocol2",
			input: "ciao.ciao.com",
			want:  "ciao.com",
		},
		{
			name:  "ok1",
			input: "http://sub.edoardottt.com/ciao?id=1",
			want:  "edoardottt.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := urlUtils.GetRootHost(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRootHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasProtocol(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "no protocol",
			input: "ciao.com",
			want:  false,
		},
		{
			name:  "no protocol2",
			input: "ciao.ciao.com",
			want:  false,
		},
		{
			name:  "ok1",
			input: "http://sub.edoardottt.com/ciao?id=1",
			want:  true,
		},
		{
			name:  "ok2",
			input: "ftp://sub.edoardottt.com/://ciao",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlUtils.HasProtocol(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HasProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}
