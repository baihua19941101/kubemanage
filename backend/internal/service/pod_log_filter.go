package service

import "strings"

func applyPodLogFilter(logs string, query PodLogQuery) string {
	if strings.TrimSpace(query.Keyword) == "" {
		return logs
	}

	lines := strings.Split(logs, "\n")
	needle := query.Keyword
	if !query.CaseSensitive {
		needle = strings.ToLower(needle)
	}

	matched := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		haystack := line
		if !query.CaseSensitive {
			haystack = strings.ToLower(line)
		}
		if strings.Contains(haystack, needle) {
			matched = append(matched, line)
		} else if !query.MatchOnly {
			matched = append(matched, line)
		}
	}
	if len(matched) == 0 {
		return ""
	}
	return strings.Join(matched, "\n") + "\n"
}
