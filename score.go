package main

import (
	"fmt"
	"sort"
	"strings"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"

	"github.com/superloach/mapper/types"
)

type locationScore struct {
	*types.Location
	Score int
}

func (l *locationScore) String() string {
	return fmt.Sprintf("%q:%d", l.Location.Name, l.Score)
}

func (l *locationScore) URL() string {
	return MapURL(fmt.Sprintf(
		"%.6f,%.6f",
		l.Value.Lat, l.Value.Lng,
	))
}

func words(s1, s2 string) int {
	a := 0
	s1s, s2s := strings.Split(s1, " "), strings.Split(s2, " ")

	for _, w1 := range s1s {
		for _, w2 := range s2s {
			if w1 == w2 {
				a++
				break
			}
		}
	}

	return (a * 100) / len(s1s)
}

func scorer(s1, s2 string) int {
	return (3*fuzzy.PartialRatio(s1, s2) +
		2*fuzzy.Ratio(s1, s2) +
		1*words(s1, s2)) / 6
}

func scoreLocation(q string, p *types.Location) *locationScore {
	names := append([]string{p.Name}, p.Aliases...)

	ls := &locationScore{
		Location: p,
		Score:    0,
	}

	for _, name := range names {
		s := scorer(q, clean(name))

		if s > ls.Score {
			ls.Score = s
		}
	}

	return ls
}

func search(q string, ls []*types.Location, min, num int) []*locationScore {
	lss := make([]*locationScore, len(ls))
	for i, p := range ls {
		lss[i] = scoreLocation(q, p)
	}

	sort.Slice(lss, func(i, j int) bool {
		return lss[i].Score > lss[j].Score
	})

	for i, ls := range lss {
		if ls.Score < min || i >= num {
			return lss[:i]
		}
	}

	return lss
}

func clean(s string) string {
	s = fuzzy.Cleanse(s, true)

	if len(s) > 50 {
		s = s[:50]
	}

	return s
}
