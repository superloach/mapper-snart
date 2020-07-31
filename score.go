package mapper

import (
	"fmt"
	"sort"
	"strings"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

type locationScore struct {
	*Location
	Score int
}

func (l *locationScore) String() string {
	return fmt.Sprintf("%q:%d", l.Location.Name, l.Score)
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

func scoreLocation(q string, p *Location) *locationScore {
	names := append([]string{p.Name}, p.Aliases...)

	ps := &locationScore{
		Location: p,
		Score:    0,
	}

	for _, name := range names {
		s := scorer(q, clean(name))

		if s > ps.Score {
			ps.Score = s
		}
	}

	return ps
}

func search(q string, ps []*Location, min, num int) []*locationScore {
	pss := make([]*locationScore, len(ps))
	for i, p := range ps {
		pss[i] = scoreLocation(q, p)
	}

	sort.Slice(pss, func(i, j int) bool {
		return pss[i].Score > pss[j].Score
	})

	for i, ps := range pss {
		if ps.Score < min || i >= num {
			return pss[:i]
		}
	}

	return pss
}

func clean(s string) string {
	s = fuzzy.Cleanse(s, true)

	if len(s) > 50 {
		s = s[:50]
	}

	return s
}
