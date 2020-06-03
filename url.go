package mapper

import "fmt"

func mapURL(s string) string {
	return "https://www.google.com/maps/dir//" + s
}

func (p *POI) URL() string {
	return mapURL(fmt.Sprintf("%.06f,%.06f", p.Loc.Lat, p.Loc.Lon))
}
