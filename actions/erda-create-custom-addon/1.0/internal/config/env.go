// Copyright (c) 2022 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"os"

	"github.com/pkg/errors"
)

var runes = map[rune]bool{
	'Z': true,
	'X': true,
	'C': true,
	'V': true,
	'B': true,
	'N': true,
	'M': true,
	'A': true,
	'S': true,
	'D': true,
	'F': true,
	'G': true,
	'H': true,
	'J': true,
	'K': true,
	'L': true,
	'Q': true,
	'W': true,
	'E': true,
	'R': true,
	'T': true,
	'Y': true,
	'U': true,
	'I': true,
	'O': true,
	'P': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
	'0': true,
	'_': true,
}

func FindEnvLiteral(s string) (string, int, int, error) {
	for i := 0; i <= len(s)-2; i++ {
		if s[i] == '$' && s[i+1] == '{' {
			for j := i; j < len(s); j++ {
				if s[j] == '\n' || s[j] == '\r' {
					return "", 0, 0, errors.New("invalid literal: '\\n' or '\\r' in string")
				}
				if s[j] == '}' {
					return s[i+2 : j], i, j + 1, nil
				}
			}
		}
	}
	return "", 0, 0, nil
}

func Interpolate(m map[string]string) error {
	for k := range m {
		s := m[k]
		for {
			key, indexStart, indexEnd, err := FindEnvLiteral(s)
			if err != nil {
				return err
			}
			if len(key) == 0 {
				break
			}
			value, ok := m[key]
			if !ok {
				value = os.Getenv(key)
			}
			if len(value) == 0 {
				return errors.Errorf("the value for %s not found in configs or envs", s[indexStart:indexEnd])
			}
			s = s[:indexStart] + value + s[indexEnd:]
		}
		m[k] = s
	}
	return nil
}
