package main

import (
	"fmt"

	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var cfgFileFields = make(map[string]interface{})
var cfgFile string

var rootCmd = &cobra.Command{
	Use:          "ipfs-p2p-helper",
	Short:        "ipfs-p2p-helper is a helper for binding k8s services to ipfs/kubo p2p.",
	Long:         fmt.Sprintf("To use, set a service's meta.labels[%s] to 'true' and meta.annotations[%s] to the p2p protocol exposed by the service", tpk8s.LabelIpfsP2P, tpk8s.LabelIpfsP2P),
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(runCmd)
}
