// SPDX-License-Identifier: GPL-3.0

package resource

import (
	"math/big"
)

var scaleFactors = []struct {
	factor *big.Float
	prefix string
}{
	{big.NewFloat(1e-15), "f"},
	{big.NewFloat(1e-12), "p"},
	{big.NewFloat(1e-9), "n"},
	{big.NewFloat(1e-6), "u"},
	{big.NewFloat(1e-3), "m"},
	{big.NewFloat(1.0), ""},
	{big.NewFloat(1e3), "K"},
	{big.NewFloat(1e6), "M"},
	{big.NewFloat(1e9), "G"},
	{big.NewFloat(1e12), "T"},
	{big.NewFloat(1e15), "P"},
}

func FormatWithScaleFactor(amount *big.Float) string {
	if amount.Cmp(big.NewFloat(0.0)) == 0 {
		return "0 "
	}
	scale := scaleFactors[0]
	for _, s := range scaleFactors {
		if amount.Cmp(s.factor) > 0 {
			scale = s
		} else {
			break
		}
	}
	scaledAmount := (&big.Float{}).Quo(amount, scale.factor)
	return scaledAmount.Text('f', 1) + " " + scale.prefix
}
