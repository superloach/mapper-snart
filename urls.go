package main

import "fmt"

func (p *POI) MapURL() string {
	return fmt.Sprintf(
		"https://www.google.com/maps/place/%.06f,%.06f",
		p.Lat, p.Lng,
	)
}

func (p *POI) DirectionsURL() string {
	return fmt.Sprintf(
		"https://www.google.com/maps/dir//%.06f,%.06f",
		p.Lat, p.Lng,
	)
}
