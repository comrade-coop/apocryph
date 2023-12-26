// SPDX-License-Identifier: GPL-3.0

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/comrade-coop/apocryph/pkg/crypto"
	tpipdr "github.com/comrade-coop/apocryph/pkg/ipdr"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	imageCopy "github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

func main() {
	err := mainErr()
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(1)
	}
}

func mainErr() error {
	ipfs, _, err := tpipfs.GetIpfsClient("")
	if err != nil {
		return err
	}
	transports.Register(tpipdr.NewIpdrTransport(ipfs))

	srcImageReference, err := alltransports.ParseImageName(os.Args[1])
	if err != nil {
		return err
	}

	dstImageReference, err := alltransports.ParseImageName(os.Args[2])
	if err != nil {
		return err
	}

	copyOptions := &imageCopy.Options{
		DestinationCtx: &types.SystemContext{},
		SourceCtx:      &types.SystemContext{},
	}

	if len(os.Args) > 3 {
		key := &pb.Key{}
		err := pb.UnmarshalFile(os.Args[3], "", key)
		if err != nil || len(key.Data) == 0 {
			key, err = crypto.NewKey(crypto.KeyTypeOcicrypt)
			if err != nil {
				return err
			}
			err = pb.MarshalFile(os.Args[3], "", key)
			if err != nil {
				return err
			}
		}

		cryptoConfig, err := crypto.GetCryptoConfigKey(key)
		if err != nil {
			return err
		}

		copyOptions.OciEncryptConfig = cryptoConfig.EncryptConfig
		copyOptions.OciDecryptConfig = cryptoConfig.DecryptConfig
		if dstImageReference.Transport().Name() == "ipdr" {
			copyOptions.OciEncryptLayers = &[]int{}
		}
	}

	policy, err := signature.DefaultPolicy(nil)
	if err != nil {
		return err
	}

	pc, _ := signature.NewPolicyContext(policy)
	defer pc.Destroy()

	_, err = imageCopy.Image(context.Background(), pc, dstImageReference, srcImageReference, copyOptions)
	if err != nil {
		return err
	}

	fmt.Println("Written image to:")
	fmt.Println(transports.ImageName(dstImageReference))

	return nil
}
