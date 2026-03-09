package main

import "fmt"

type stubTranslator struct {
	prefix string
	err    error
	calls  int
	last   []string
}

func (s *stubTranslator) Translate(_ language, texts []string, _ bool) ([]string, error) {
	s.calls++
	s.last = append([]string(nil), texts...)
	if s.err != nil {
		return nil, s.err
	}
	out := make([]string, len(texts))
	for i, text := range texts {
		out[i] = fmt.Sprintf("%s%s", s.prefix, text)
	}
	return out, nil
}
