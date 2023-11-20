package resource

import (
	"fmt"
	"io"
	"math/big"
)

type consumption *big.Float
type ResourceConsumptions map[*Resource]consumption
type NamespaceConsumptions map[string]ResourceConsumptions

func (r NamespaceConsumptions) Add(namespace string, resource *Resource, value *big.Float) {
	if r[namespace] == nil {
		r[namespace] = ResourceConsumptions{}
	}

	if r[namespace][resource] != nil {
		r[namespace][resource] = (*big.Float)(r[namespace][resource]).Add(r[namespace][resource], value)
	} else {
		r[namespace][resource] = value
	}
}

var CurrencyDisplayScale = big.NewFloat(1.0e-18)
var currency = GetResource("Tokens", ResourceKindUsage)

func (r NamespaceConsumptions) Display(writer io.Writer, priceTableMap PricingTableMap) {
	fmt.Fprint(writer, "Resources:\n")

	for namespace, resources := range r {
		totalCost := &big.Float{}
		fmt.Fprintf(writer, "- namespace: %s\n", namespace)
		for resource, consumption := range resources {
			fmt.Fprintf(writer, "  - %s %s: %s-seconds\n", resource.Name, resource.Kind.Text(), resource.Format(consumption))
			if priceTableMap != nil && priceTableMap[resource] != nil {
				cost := (&big.Float{}).Mul(priceTableMap[resource], consumption)
				fmt.Fprintf(writer, "    @ %s/%s = %s\n", currency.Format((&big.Float{}).Mul(priceTableMap[resource], CurrencyDisplayScale)), resource.Unit, currency.Format((&big.Float{}).Mul(cost, CurrencyDisplayScale)))

				totalCost = totalCost.Add(totalCost, cost)
			}
		}
		if priceTableMap != nil {
			fmt.Fprintf(writer, "  total: %s\n", currency.Format((&big.Float{}).Mul(totalCost, CurrencyDisplayScale)))
		}
	}
}

func (r ResourceConsumptions) Price(priceTableMap PricingTableMap) *big.Int {
	totalCost := &big.Float{}
	for resource, consumption := range r {
		if priceTableMap[resource] != nil {
			cost := (&big.Float{}).Mul(priceTableMap[resource], consumption)
			totalCost = totalCost.Add(totalCost, cost)
		}
	}
	totalCostInt := &big.Int{}
	_, _ = totalCost.Int(totalCostInt)
	return totalCostInt
}

func (r NamespaceConsumptions) Price(priceTableMap PricingTableMap) map[string]*big.Int {
	costs := make(map[string]*big.Int, len(r))
	for namespace, resources := range r {
		costs[namespace] = resources.Price(priceTableMap)
	}
	return costs
}
