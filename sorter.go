package main

import (
	"sort"
)

// By ...
type By func(p1, p2 *Runner) bool

// Sort ...
func (by By) Sort(runners []Runner) {
	ps := &runnerSorter{
		runners: runners,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type runnerSorter struct {
	runners []Runner
	by      func(p1, p2 *Runner) bool // Closure used in the Less method.
}

func (s *runnerSorter) Len() int {
	return len(s.runners)
}

func (s *runnerSorter) Swap(i, j int) {
	s.runners[i], s.runners[j] = s.runners[j], s.runners[i]
}

func (s *runnerSorter) Less(i, j int) bool {
	return s.by(&s.runners[i], &s.runners[j])
}
