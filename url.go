package mapper

import (
	"fmt"
	"net/url"
)

func mapURL(s string) string {
	return "https://www.google.com/maps/dir//" + url.PathEscape(s)
}

func (p *POI) URL() string {
	return mapURL(fmt.Sprintf("%.06f,%.06f", p.Loc.Lat, p.Loc.Lon))
}
