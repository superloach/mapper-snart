package mapper

import (
	"fmt"
)

type IngrType uint8

const (
	IngrTypeUnknown IngrType = iota // 0
	IngrTypeNone                    // 1
	IngrTypePortal                  // 2
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

type PkmnType uint8

const (
	PkmnTypeUnknown PkmnType = iota // 0
	PkmnTypeNone                    // 1
	PkmnTypeStop                    // 2
	PkmnTypeGym                     // 3
	PkmnTypeEXGym                   // 4
	PkmnTypeNest                    // 5
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

type WzrdType uint8

const (
	WzrdTypeUnknown    WzrdType = iota // 0
	WzrdTypeNone                       // 1
	WzrdTypeInn                        // 2
	WzrdTypeGreenhouse                 // 3
	WzrdTypeFortress                   // 4
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
