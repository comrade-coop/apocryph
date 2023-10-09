package resource

import (
	"math/big"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
)

type PricingTableMap map[*Resource]*big.Float

func NewPricingTableMap(table *pb.PricingTable) PricingTableMap {
	res := make(PricingTableMap)
	for _, pr := range table.Resources {
		if pr.PriceForReservation != 0 {
			res[GetResource(pr.Resource, ResourceKindReservation)] = big.NewFloat(float64(pr.PriceForReservation))
		}
		if pr.PriceForUsage != 0 {
			res[GetResource(pr.Resource, ResourceKindUsage)] = big.NewFloat(float64(pr.PriceForUsage))
		}
	}
	return res
}
