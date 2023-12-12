// SPDX-License-Identifier: GPL-3.0

// this package contains ipfs helper functions
package ipfs

import (
	"os"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/multiformats/go-multiaddr"
)

func GetIpfsClient(ipfsApi string) (api *rpc.HttpApi, apiMultiaddr multiaddr.Multiaddr, err error) {
	if ipfsApi == "" {
		// via rpc.NewLocalApi()
		ipfspath := os.Getenv(rpc.EnvDir)
		if ipfspath == "" {
			ipfspath = rpc.DefaultPathRoot
		}
		apiMultiaddr, err = rpc.ApiAddr(ipfspath)
		if err != nil {
			if os.IsNotExist(err) {
				err = rpc.ErrApiNotFound
			}
			return
		}
	} else {
		apiMultiaddr, err = multiaddr.NewMultiaddr(ipfsApi)
		if err != nil {
			return
		}
	}
	api, err = rpc.NewApi(apiMultiaddr)
	return
}
