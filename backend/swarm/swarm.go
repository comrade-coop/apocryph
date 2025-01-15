package swarm

import (
	"math/rand"
	"sort"
	"strconv"

	"github.com/hashicorp/serf/client"
)

type BucketStatus string

const (
	Syncing BucketStatus = "Syncing"
	Ready   BucketStatus = "Ready"
	Missing BucketStatus = ""
)
const BucketPrefix = "b_"
const CapacityTag = "meta_capacity"

type Swarm struct {
	serf    *client.RPCClient
	OwnName string
}

func NewSwarm(serfAddress string, ownNodeName string) (*Swarm, error) {
	serf, err := client.ClientFromConfig(&client.Config{
		Addr: serfAddress,
	})
	if err != nil {
		return nil, err
	}
	return &Swarm{
		serf:    serf,
		OwnName: ownNodeName,
	}, nil
}

func (s *Swarm) Join(existingNode string) error {
	_, err := s.serf.Join([]string{existingNode}, false)
	return err
}

func (s *Swarm) FindBucketBestNodes(bucketId string) ([]string, error) {
	bucketKey := BucketPrefix + bucketId
	members, err := s.serf.MembersFiltered(map[string]string{
		bucketKey: string(Ready),
	}, "Alive", "")
	if err != nil {
		return nil, err
	}
	if len(members) == 0 { // No ready members, try looking for syncing ones
		members, err = s.serf.MembersFiltered(map[string]string{
			bucketKey: string(Syncing) + "|" + string(Ready),
		}, "Alive", "")
		if err != nil {
			return nil, err
		}
	}
	if len(members) == 0 { // No alive members, try returning a few dead ones
		members, err = s.serf.MembersFiltered(map[string]string{
			bucketKey: string(Syncing) + "|" + string(Ready),
		}, "", "")
		if err != nil {
			return nil, err
		}
	}
	resultAddresses := make([]string, len(members))
	for i := range members {
		resultAddresses[i] = members[i].Name
	}
	return resultAddresses, nil
}

func (s *Swarm) FindBucketReplicas(bucketId string) ([]string, error) {
	bucketKey := BucketPrefix + bucketId
	members, err := s.serf.MembersFiltered(map[string]string{
		bucketKey: string(Syncing) + "|" + string(Ready),
	}, "Alive", "")
	if err != nil {
		return nil, err
	}
	resultAddresses := make([]string, len(members))
	for i := range members {
		resultAddresses[i] = members[i].Name
	}
	return resultAddresses, nil
}

func (s *Swarm) UpdateBucket(bucketId string, status BucketStatus) error {
	bucketKey := BucketPrefix + bucketId
	if status == Missing {
		return s.serf.UpdateTags(map[string]string{}, []string{bucketKey})
	} else {
		return s.serf.UpdateTags(map[string]string{
			bucketKey: string(status),
		}, []string{})
	}
}

// TODO: Make sure to fliter ourselves out of the list
func (s *Swarm) FindFreeNode() (string, error) {
	members, err := s.serf.MembersFiltered(map[string]string{
		CapacityTag: "",
	}, "Alive", "")
	if err != nil {
		return "", err
	}

	parsedCapacities := make([]uint64, len(members))
	for i := range members {
		parsedCapacities[i], _ = strconv.ParseUint(members[i].Tags[CapacityTag], 10, 16)
	}
	sort.Slice(members, func(i, j int) bool {
		return parsedCapacities[i] > parsedCapacities[j]
	})
	// Arbitrarily pick a server in the top half
	// TODO: there probably are better load-balancing algorithms
	picked := rand.Intn((len(members) + 1) / 2)
	return members[picked].Name, nil
}

func (s *Swarm) UpdateCapacity(capacityLeft uint64) error {
	return s.serf.UpdateTags(map[string]string{
		CapacityTag: strconv.FormatUint(capacityLeft, 10),
	}, []string{})
}
