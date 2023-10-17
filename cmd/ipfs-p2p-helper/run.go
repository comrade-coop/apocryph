package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ipfsApi string
var kubeConfig string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start observing Kubernetes services and registering them in ipfs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := ipfs_utils.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}
		for {
			time.Sleep(time.Second * 1)
			_, err := ipfs.Key().Self(cmd.Context())
			if err == nil {
				break
			}
			fmt.Printf("IPFS is not started yet, %s", err)
		}
		ipfsp2p := ipfs_utils.NewP2pApi(ipfs, ipfsMultiaddr)

		scheme, err := tpk8s.GetScheme()
		if err != nil {
			return err
		}

		config, err := tpk8s.GetConfig(kubeConfig)
		if err != nil {
			return err
		}

		cl, err := client.NewWithWatch(config, client.Options{
			Scheme: scheme,
		})
		if err != nil {
			return err
		}

		services := &corev1.ServiceList{}
		sub, err := cl.Watch(cmd.Context(), services, client.HasLabels{tpk8s.LabelIpfsP2P})
		if err != nil {
			return err
		}
		defer sub.Stop()

		fmt.Print("Now watching the list of services.\n")

	Loop:
		for {
			select {
			case e := <-sub.ResultChan():
				if e.Type == watch.Error {
					fmt.Printf("Error: %v", e)
					return errors.New("Watch resulted in error!")
				}
				if service, ok := e.Object.(*corev1.Service); ok {
					err := handleEvent(ipfsp2p, e.Type, service)
					if err != nil {
						if e.Type == watch.Added {
							fmt.Fprintf(cmd.ErrOrStderr(), "Service %s: %s\n", service.Name, err)
						}
						continue Loop
					}
				}
			case <-cmd.Context().Done():
				break Loop
			}
		}
		return nil
	},
}

func handleEvent(ipfsp2p *ipfs_utils.P2pApi, eType watch.EventType, service *corev1.Service) error {
	if len(service.Spec.Ports) != 1 {
		return errors.New("Expected exactly one exposed port")
	}

	protocol := service.ObjectMeta.Annotations[tpk8s.LabelIpfsP2P]

	port := service.Spec.Ports[0].Port
	portProtocol := strings.ToLower(string(service.Spec.Ports[0].Protocol))
	dns := fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace)
	endpoint, err := multiaddr.NewMultiaddr(fmt.Sprintf("/dns4/%s/%s/%d", dns, portProtocol, port))
	if err != nil {
		return err
	}

	if eType != watch.Deleted {
		if eType == watch.Added {
			fmt.Printf("Forwarding %s to %s\n", protocol, endpoint)
		}
		_, err := ipfsp2p.ExposeEndpoint(protocol, endpoint, ipfs_utils.ReturnExistingEndpoint)
		return err
	} else {
		fmt.Printf("Dropping forward for %s to %s\n", protocol, endpoint)
		endpoint, err := ipfsp2p.ExposeEndpoint(protocol, endpoint, ipfs_utils.ReturnExistingEndpoint)
		if err != nil {
			return err
		}

		err = endpoint.Close()
		return err
	}
}

func init() {
	runCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to use in-cluster config)")
	runCmd.Flags().StringVar(&ipfsApi, "ipfs", "-", "multiaddr where the ipfs/kubo api can be accessed (- to use the daemon running in IPFS_PATH)")
}
