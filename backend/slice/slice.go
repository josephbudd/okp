package slice

import (
	"sort"
)

func UniqueStrings(ss []string) (unique []string) {
	var lss int
	if lss = len(ss); lss < 2 {
		// 0-1 items so there are no duplicates.
		unique = ss
		return
	}

	// Sort the slice so that equal strings follow one another.
	// Ex: {"A", "A", "B", "C", "C", "C", "D" }
	sort.Strings(ss)

	var keepAt int
	for i := 1; i < lss; i++ {
		if ss[keepAt] == ss[i] {
			// ss[i] is a duplicate of ss[keepAt] so don't keep ss[i].
			continue
		}
		// ss[i] is not a duplicate of ss[keepAt] so keep ss[i].
		// Inc keepAt to the next available item index.
		keepAt++
		ss[keepAt] = ss[i]
	}

	unique = ss[:keepAt+1]
	return
}
