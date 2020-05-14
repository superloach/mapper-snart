package mapper

import "fmt"

func (p *POI) URL() string {
	return fmt.Sprintf(
		"https://www.google.com/maps/dir//%.06f,%.06f",
		p.Lat, p.Lng,
	)
}
