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

func (s *Scanner) ScanToDelimiter() string {
	return s.ReadWhile(func(r rune) bool {
		return r != token.TokenDelimiter && r != token.TokenLeftParen && r != token.TokenRightParen
	})
}

func (s *Scanner) ScanToEndOfLine() string {
	return s.ReadWhile(func(r rune) bool {
		return r != token.TokenNewLine && r != token.TokenReturn && r != 0
	})
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

			param := s.ScanToDelimiter()

			params = append(params, param)

			s.SkipWhitespace()

			if s.CurrentRune == token.TokenDelimiter {
				s.ReadNext()
			} else {
				break
			}
		}

		s.ReadNext()

		return blockName, params
	}

	return blockName, nil
}
