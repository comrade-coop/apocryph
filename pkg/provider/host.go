package provider

import "fmt"

type HostInfo struct {
	Id             string         `yaml:"id"`
	Regions        []Region       `yaml:"regions"`
	MultiAddresses []MultiAddress `yaml:"multiAddresses"`
}

type Region struct {
	Name string `yaml:"name"`
	Zone string `yaml:"zone"`
	Num  uint32 `yaml:"num"`
}

func (r Region) IsEqual(other Region) bool {
	return r.Name == other.Name && r.Zone == other.Zone && r.Num == other.Num
}

type MultiAddress struct {
	Kind  string `yaml:"kind"`
	Value string `yaml:"value"`
}

func (h *HostInfo) FormatAddresses() string {
	var addresses string = ""
	for _, addr := range h.MultiAddresses {
		addresses = addresses + addr.Kind + ":" + addr.Value + "\n"
	}
	return addresses
}

func (h *HostInfo) FormatRegions() string {
	var regions string = ""
	for _, region := range h.Regions {
		regions = regions + region.Name + "-" + region.Zone + "-" + fmt.Sprint(region.Num) + "\n"
	}
	return regions
}
