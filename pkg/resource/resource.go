package resource

import (
	"math/big"
)

type ResourceKind uint8

const (
	ResourceKindReservation ResourceKind = 0
	ResourceKindUsage       ResourceKind = 1
)

func (r ResourceKind) Text() string {
	switch r {
	case ResourceKindReservation:
		return "reserved"
	case ResourceKindUsage:
		return "used"
	}
	return "?"
}

type Resource struct {
	Name      string
	Kind      ResourceKind
	Unit      string
	UseMillis bool
}

func (r *Resource) Format(amount *big.Float) string {
	return FormatWithScaleFactor(amount) + r.Unit
}
