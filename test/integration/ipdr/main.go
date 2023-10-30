package main

import (
	"context"
	"fmt"
	"os"

	tpipdr "github.com/comrade-coop/trusted-pods/pkg/ipdr"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	imageCopy "github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/transports/alltransports"
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

	policy, err := signature.DefaultPolicy(nil)
	if err != nil {
		return err
	}

	pc, _ := signature.NewPolicyContext(policy)
	defer pc.Destroy()

	_, err = imageCopy.Image(context.Background(), pc, dstImageReference, srcImageReference, &imageCopy.Options{})
	if err != nil {
		return err
	}

	fmt.Println("Written image to:")
	fmt.Println(transports.ImageName(dstImageReference))

	return nil
}
