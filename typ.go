package mapper

import (
	"fmt"
)

// IngrType indicates a Location's type in Ingress.
type IngrType uint8

const (
	// IngrTypeUnknown indicates that a Location is unknown in Ingress.
	IngrTypeUnknown IngrType = iota

	// IngrTypeNone indicates that a Location is not in Ingress.
	IngrTypeNone

	// IngrTypePortal indicates that a Location is a Portal in Ingress.
	IngrTypePortal
)

func (i IngrType) String() string {
	switch i {
	case IngrTypeUnknown:
		return "Unknown"
	case IngrTypeNone:
		return "None"
	case IngrTypePortal:
		return "Portal"
	}

	return fmt.Sprintf("invalid IngrType %d", i)
}

// PkmnType indicates a Location's type in Pokémon GO.
type PkmnType uint8

const (
	// PkmnTypeUnknown indicates that a Location is unknown in Pokémon GO.
	PkmnTypeUnknown PkmnType = iota

	// PkmnTypeNone indicates that a Location is not in Pokémon GO.
	PkmnTypeNone

	// PkmnTypeStop indicates that a Location is a PokéStop in Pokémon GO.
	PkmnTypeStop

	// PkmnTypeGym indicates that a Location is a Gym in Pokémon GO.
	PkmnTypeGym

	// PkmnTypeEXGym indicates that a Location is an EX-eligible Gym in Pokémon GO.
	PkmnTypeEXGym

	// PkmnTypeNest indicates that a Location is a Nest in Pokémon GO.
	PkmnTypeNest
)

func (p PkmnType) String() string {
	switch p {
	case PkmnTypeUnknown:
		return "Unknown"
	case PkmnTypeNone:
		return "None"
	case PkmnTypeStop:
		return "PokeStop"
	case PkmnTypeGym:
		return "Gym"
	case PkmnTypeEXGym:
		return "EX Gym"
	case PkmnTypeNest:
		return "Nest"
	}

	return fmt.Sprintf("invalid PkmnType %d", p)
}

// WzrdType indicates a Location's type in Harry Potter: Wizards Unite.
type WzrdType uint8

const (
	// WzrdTypeUnknown indicates that a Location is unknown in Harry Potter: Wizards Unite.
	WzrdTypeUnknown WzrdType = iota

	// WzrdTypeNone indicates that a Location is not in Harry Potter: Wizards Unite.
	WzrdTypeNone

	// WzrdTypeInn indicates that a Location is an Inn in Harry Potter: Wizards Unite.
	WzrdTypeInn

	// WzrdTypeGreenhouse indicates that a Location is a Greenhouse in Harry Potter: Wizards Unite.
	WzrdTypeGreenhouse

	// WzrdTypeFortress indicates that a Location is a Fortress in Harry Potter: Wizards Unite.
	WzrdTypeFortress
)

func (w WzrdType) String() string {
	switch w {
	case WzrdTypeUnknown:
		return "Unknown"
	case WzrdTypeNone:
		return "None"
	case WzrdTypeInn:
		return "Inn"
	case WzrdTypeGreenhouse:
		return "Greenhouse"
	case WzrdTypeFortress:
		return "Fortress"
	}

	return fmt.Sprintf("invalid WzrdType %d", w)
}
