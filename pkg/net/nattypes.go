package net

type NATType int

func (n NATType) String() string {
	switch n {
	case None:
		return "None"
	case FullCone:
		return "FullCone"
	case RestrictedCone:
		return "RestrictedCone"
	case FullOrRestrictedCone:
		return "FullOrRestrictedCone"
	case PortRestrictedCone:
		return "PortRestrictedCone"
	case Symmetric:
		return "Symmetric"
	case UnKnown:
		return "UnKnown"
	default:
		return "UnKnown"
	}
}

const (
	UnKnown NATType = iota
	None
	FullCone
	RestrictedCone
	FullOrRestrictedCone
	PortRestrictedCone
	Symmetric
)
