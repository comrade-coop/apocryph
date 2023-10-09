package resource

import (
	"fmt"
	"io"
	"math/big"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
)

type ResourceMeasurementsMap map[string]map[*Resource]*big.Float

func (r ResourceMeasurementsMap) Add(namespace string, resource *Resource, value *big.Float) {
	if r[namespace] == nil {
		r[namespace] = map[*Resource]*big.Float{}
	}

	if r[namespace][resource] != nil {
		r[namespace][resource] = r[namespace][resource].Add(r[namespace][resource], value)
	} else {
		r[namespace][resource] = value
	}
}

var CurrencyDisplayScale = big.NewFloat(1.0e-16)

func (r ResourceMeasurementsMap) Display(writer io.Writer, priceTable *pb.PricingTable) {
	var priceTableMap PricingTableMap
	var currency *Resource
	if priceTable != nil {
		priceTableMap = NewPricingTableMap(priceTable)
		currency = GetResource(priceTable.Currency, ResourceKindUsage)
	}

	fmt.Fprint(writer, "Resources:\n")

	for namespace, resources := range r {
		totalCost := &big.Float{}
		fmt.Fprintf(writer, "- namespace: %s\n", namespace)
		for resource, measurement := range resources {
			fmt.Fprintf(writer, "  - %s %s: %s-seconds\n", resource.Name, resource.Kind.Text(), resource.Format(measurement))
			if priceTable != nil && priceTableMap[resource] != nil {
				cost := (&big.Float{}).Mul(priceTableMap[resource], measurement)
				fmt.Fprintf(writer, "    @ %s/%s = %s\n", currency.Format((&big.Float{}).Mul(priceTableMap[resource], CurrencyDisplayScale)), resource.Unit, currency.Format((&big.Float{}).Mul(cost, CurrencyDisplayScale)))

				totalCost = totalCost.Add(totalCost, cost)
			}
		}
		if priceTable != nil {
			fmt.Fprintf(writer, "  total: %s\n", currency.Format((&big.Float{}).Mul(totalCost, CurrencyDisplayScale)))
		}
	}
}

func (r ResourceMeasurementsMap) Price(priceTable *pb.PricingTable) map[string]uint64 {
	costs := make(map[string]uint64, len(r))
	priceTableMap := NewPricingTableMap(priceTable)
	for namespace, resources := range r {
		totalCost := &big.Float{}
		for resource, measurement := range resources {
			if priceTableMap[resource] != nil {
				cost := (&big.Float{}).Mul(priceTableMap[resource], measurement)
				totalCost = totalCost.Add(totalCost, cost)
			}
		}
		totalCostInt, _ := totalCost.Int64()
		costs[namespace] = uint64(totalCostInt)
	}
	return costs
}
