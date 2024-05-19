/*
 * gomake
 * Copyright (C) 2024 gomake contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package scanner

import (
	"strings"

	"github.com/ricochhet/gomake/token"
)

type Scanner struct {
	Text        string
	Position    int
	CurrentRune rune
}

func NewScanner(text string) *Scanner {
	s := &Scanner{Text: text, Position: 0, CurrentRune: 0}
	s.ReadNext()

	return s
}

func (s *Scanner) ReadNext() {
	if s.Position < len(s.Text) {
		s.CurrentRune = rune(s.Text[s.Position])
		s.Position++
	} else {
		s.CurrentRune = 0
	}
}

func (s *Scanner) ReadAhead(n int) {
	if s.Position+n-1 < len(s.Text) {
		s.CurrentRune = rune(s.Text[s.Position+n-1])
		s.Position += n
	} else {
		s.CurrentRune = 0
	}
}

//nolint:lll // wontfix
func (s *Scanner) SkipWhitespace() {
	for s.CurrentRune != 0 && (s.CurrentRune == token.TokenSpace || s.CurrentRune == token.TokenTab || s.CurrentRune == token.TokenNewLine || s.CurrentRune == token.TokenReturn) {
		s.ReadNext()
	}
}

func (s *Scanner) ReadWhile(predicate func(rune) bool) string {
	var result strings.Builder
	for s.CurrentRune != 0 && predicate(s.CurrentRune) {
		result.WriteRune(s.CurrentRune)
		s.ReadNext()
	}

	return result.String()
}

func (s *Scanner) IsIndentifiable(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func (s *Scanner) ScanIdentifier() string {
	return s.ReadWhile(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
	})
}

func (s *Scanner) ScanToUnescaped(target rune) string {
	return s.ReadWhile(func(r rune) bool {
		if s.Peek(-2) != token.TokenEscape && r == target {
			return false
		}

		return true
	})
}

func (s *Scanner) ScanToEndOfLine() string {
	return s.ReadWhile(func(r rune) bool {
		return r != token.TokenNewLine && r != token.TokenReturn && r != 0
	})
}

func (s *Scanner) Peek(n int) rune {
	peek := s.Position + n

	if peek < len(s.Text) {
		return rune(s.Text[peek])
	}

	return 0
}

func (s *Scanner) PeekAhead(n int) string {
	end := s.Position + n

	if end > len(s.Text) {
		end = len(s.Text)
	}

	return s.Text[s.Position:end]
}

func (s *Scanner) ScanBlockWithParams() (string, []string) {
	blockName := s.ScanIdentifier()

	if s.CurrentRune == token.TokenLeftParen {
		s.ReadNext()

		params := make([]string, 0)

		for {
			if s.CurrentRune == token.TokenRightParen {
				s.ReadNext()
				break
			}

			s.SkipWhitespace()
			s.ReadNext()

			param := s.ScanToUnescaped(token.TokenQuote)

			params = append(params, param)

			if s.Peek(0) == token.TokenDelimiter {
				s.ReadNext()
			}

			s.SkipWhitespace()

			if s.CurrentRune == token.TokenDelimiter {
				s.ReadNext()
			} else {
				s.ReadNext()
				break
			}
		}

		s.ReadNext()

		return blockName, params
	}

	return blockName, nil
}

func (s *Scanner) ScanParams() []string {
	params := make([]string, 0)

	for {
		if s.CurrentRune == token.TokenRightParen {
			s.ReadNext()
			break
		}

		s.SkipWhitespace()
		s.ReadNext()

		param := s.ScanToUnescaped(token.TokenQuote)

		params = append(params, param)

		if s.Peek(0) == token.TokenDelimiter {
			s.ReadNext()
		}

		s.SkipWhitespace()

		if s.CurrentRune == token.TokenDelimiter {
			s.ReadNext()
		} else {
			s.ReadNext()
			break
		}
	}

	s.ReadNext()

	return params
}
