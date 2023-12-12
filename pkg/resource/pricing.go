// SPDX-License-Identifier: GPL-3.0

package resource

import (
	"fmt"
	"math/big"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
)

type PricingTableMap map[*Resource]*big.Float

const (
	cpu         = 0
	ram         = 1
	storage     = 2
	bandwidthE  = 3
	bandwidthIn = 4
)

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

const numRessources = 5

func GetTablesPrices(pricingTables map[common.Address]PricingTableMap) ([][numRessources]*big.Int, error) {
	var allPrices [][numRessources]*big.Int
	for _, table := range pricingTables {
		var prices [numRessources]*big.Int
		for i := 0; i < numRessources; i++ {
			prices[i] = new(big.Int).SetInt64(0)
		}
		for resource, price := range table {
			p := new(big.Int)
			price.Int(p)
			switch resource.Name {
			case "cpu":
				prices[cpu] = p
				break
			case "ram":
				prices[ram] = p
				break
			case "storage":
				prices[storage] = p
				break
			case "bandwidth_egress":
				prices[bandwidthE] = p
				break
			case "bandwidth_ingress":
				prices[bandwidthIn] = p
				break
			default:
				return nil, fmt.Errorf("unrecognized ressource name")
			}
		}
		allPrices = append(allPrices, prices)
	}
	return allPrices, nil
}
