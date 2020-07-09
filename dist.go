package mapper

import (
	"fmt"
	"math"
	"sort"

	"github.com/go-snart/snart/db"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/types"
)

type neighDist struct {
	*Neigh
	Dist float64
}

func (n neighDist) String() string {
	return fmt.Sprintf("%s [%.2fkm]", n.Name, n.Dist/1000)
}

// GetNeigh calculates the nearest Neigh and sets the POI's Neigh value accordingly.
func (p *POI) GetNeigh(d *db.DB) {
	_f := "(*POI).GetNeighs"

	if p.Neigh != nil {
		return
	}

	dists := []neighDist{}

	if p.Loc == nil {
		Log.Warnf(_f, "poi %#v has nil loc", p)
		return
	}

	d.Cache.Lock()
	cache := d.Cache.Get("mapper.neigh").(db.Cache)
	d.Cache.Unlock()

	cache.Lock()
	keys := cache.Keys()
	cache.Unlock()

	for _, key := range keys {
		cache.Lock()
		neigh := cache.Get(key).(*Neigh)
		cache.Unlock()

		if neigh.Loc == nil {
			Log.Warnf(_f, "neigh %#v has nil loc", neigh)
			continue
		}

		dists = append(dists, neighDist{
			Neigh: neigh,
			Dist:  distance(*p.Loc, *neigh.Loc),
		})
	}

	if len(dists) == 0 {
		return
	}

	sort.Slice(dists, func(i, j int) bool {
		return dists[i].Dist < dists[j].Dist
	})

	p.Neigh = &dists[0].Name
}

func distance(a, b types.Point) float64 {
	lat1, lng1, lat2, lng2 := a.Lat, a.Lon, b.Lat, b.Lon

	const R = 6371e3

	rlat1 := lat1 * math.Pi / 180
	rlat2 := lat2 * math.Pi / 180

	dlat := (lat2 - lat1) * math.Pi / 180
	dlng := (lng2 - lng1) * math.Pi / 180

	s := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(rlat1)*math.Cos(rlat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	t := 2 * math.Atan2(math.Sqrt(s), math.Sqrt(1-s))

	return R * t
}
