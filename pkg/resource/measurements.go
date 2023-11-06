package resource

import (
	"fmt"
	"io"
	"math/big"
)

type ResourceMeasurements map[*Resource]*big.Float
type ResourceMeasurementsMap map[string]ResourceMeasurements

func (r ResourceMeasurementsMap) Add(namespace string, resource *Resource, value *big.Float) {
	if r[namespace] == nil {
		r[namespace] = ResourceMeasurements{}
	}

	if r[namespace][resource] != nil {
		r[namespace][resource] = r[namespace][resource].Add(r[namespace][resource], value)
	} else {
		r[namespace][resource] = value
	}
}

var CurrencyDisplayScale = big.NewFloat(1.0e-18)
var currency = GetResource("Tokens", ResourceKindUsage)

func (r ResourceMeasurementsMap) Display(writer io.Writer, priceTableMap PricingTableMap) {
	fmt.Fprint(writer, "Resources:\n")

	for namespace, resources := range r {
		totalCost := &big.Float{}
		fmt.Fprintf(writer, "- namespace: %s\n", namespace)
		for resource, measurement := range resources {
			fmt.Fprintf(writer, "  - %s %s: %s-seconds\n", resource.Name, resource.Kind.Text(), resource.Format(measurement))
			if priceTableMap != nil && priceTableMap[resource] != nil {
				cost := (&big.Float{}).Mul(priceTableMap[resource], measurement)
				fmt.Fprintf(writer, "    @ %s/%s = %s\n", currency.Format((&big.Float{}).Mul(priceTableMap[resource], CurrencyDisplayScale)), resource.Unit, currency.Format((&big.Float{}).Mul(cost, CurrencyDisplayScale)))

				totalCost = totalCost.Add(totalCost, cost)
			}
		}
		if priceTableMap != nil {
			fmt.Fprintf(writer, "  total: %s\n", currency.Format((&big.Float{}).Mul(totalCost, CurrencyDisplayScale)))
		}
	}
}

func (r ResourceMeasurements) Price(priceTableMap PricingTableMap) *big.Int {
	totalCost := &big.Float{}
	for resource, measurement := range r {
		if priceTableMap[resource] != nil {
			cost := (&big.Float{}).Mul(priceTableMap[resource], measurement)
			totalCost = totalCost.Add(totalCost, cost)
		}
	}
	totalCostInt := &big.Int{}
	_, _ = totalCost.Int(totalCostInt)
	return totalCostInt
}

func (r ResourceMeasurementsMap) Price(priceTableMap PricingTableMap) map[string]*big.Int {
	costs := make(map[string]*big.Int, len(r))
	for namespace, resources := range r {
		costs[namespace] = resources.Price(priceTableMap)
	}
	return costs
}
