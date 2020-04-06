package main

func (p *POI) MapURL() string {
	_f := "(*POI).MapURL"

	Log.Warn(_f, "placeholder map url")
	return "https://google.com/"
}

func (p *POI) DirectionsURL() string {
	_f := "(*POI).DirectionsURL"

	Log.Warn(_f, "placeholder directions url")
	return "https://google.com/"
}
