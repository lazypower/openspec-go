package parser

import (
	"strings"

	"github.com/chuck/openspec-go/internal/model"
)

// ParseChange parses a change proposal markdown file.
func ParseChange(content string) model.Change {
	change := model.Change{}
	lines := strings.Split(content, "\n")

	var currentSection string
	var contentBuf []string

	flushSection := func() {
		text := strings.TrimSpace(strings.Join(contentBuf, "\n"))
		switch strings.ToLower(currentSection) {
		case "why":
			change.Why = text
		case "what changes":
			change.WhatChanges = text
		case "impact":
			change.Impact = text
		}
		contentBuf = nil
	}

	for _, line := range lines {
		m := headerRe.FindStringSubmatch(line)
		if m == nil {
			contentBuf = append(contentBuf, line)
			continue
		}

		level := len(m[1])
		text := strings.TrimSpace(m[2])

		switch level {
		case 1:
			flushSection()
			// Extract title from "Change: Title" format
			if name := extractAfterColon(text, "Change"); name != "" {
				change.Title = name
			} else {
				change.Title = text
			}
			currentSection = ""
			contentBuf = nil

		case 2:
			flushSection()
			currentSection = text
			contentBuf = nil

		default:
			contentBuf = append(contentBuf, line)
		}
	}

	flushSection()
	return change
}
