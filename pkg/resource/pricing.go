package resource

import (
	"math/big"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
)

type PricingTableMap map[*Resource]*big.Float

func NewPricingTableMap(table *pb.PricingTable) PricingTableMap {
	res := make(PricingTableMap, len(table.Resources))
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

func ConvertPricingTables(tables []*pb.PricingTable) map[common.Address]PricingTableMap {
	result := make(map[common.Address]PricingTableMap, len(tables))
	for _, table := range tables {
		result[common.BytesToAddress(table.PaymentContractAddress)] = NewPricingTableMap(table)
	}
	return result
}
