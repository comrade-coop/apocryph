// SPDX-License-Identifier: GPL-3.0

package resource

var units = map[string]string{
	"cpu":     "cpu",
	"memory":  "B",
	"storage": "B",
}

var resources = map[struct {
	string
	ResourceKind
}]*Resource{}

func GetResource(name string, kind ResourceKind) *Resource {
	key := struct {
		string
		ResourceKind
	}{name, kind}
	if res, ok := resources[key]; ok {
		return res
	}

	unit, ok := units[name]
	if !ok {
		unit = name
	}

	res := &Resource{
		Name: name,
		Kind: kind,
		Unit: unit,
	}

	resources[key] = res

	return res
}
