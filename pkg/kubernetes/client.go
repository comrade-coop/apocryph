// SPDX-License-Identifier: GPL-3.0

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	kedahttpscheme "github.com/kedacore/http-add-on/operator/generated/clientset/versioned/scheme"
	policy "github.com/sigstore/policy-controller/pkg/apis/policy/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetScheme() (*runtime.Scheme, error) {
	sch := runtime.NewScheme()
	if err := scheme.AddToScheme(sch); err != nil {
		return nil, err
	}
	if err := policy.AddToScheme(sch); err != nil {
		return nil, err
	}
	if err := kedahttpscheme.AddToScheme(sch); err != nil {
		return nil, err
	}

	return sch, nil
}

func GetConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig == "" {
		config, err := rest.InClusterConfig()

		if err == rest.ErrNotInCluster {
			defaultKubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
			config, err2 := clientcmd.BuildConfigFromFlags("", defaultKubeConfigPath)
			if err2 != nil {
				return nil, errors.Join(err, err2)
			}
			return config, nil
		}

		return config, err
	} else {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
}

func GetClient(kubeConfig string, dryRun bool) (client.Client, error) {
	config, err := GetConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	sch, err := GetScheme()
	if err != nil {
		return nil, err
	}

	cl, err := client.New(config, client.Options{
		Scheme: sch,
		DryRun: &dryRun,
	})
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func RunInNamespaceOrRevert(ctx context.Context, cl client.Client, namespace *corev1.Namespace, dryRun bool, block func(client.Client) error) error {
	if err := cl.Create(ctx, namespace); err != nil {
		return err
	}

	if dryRun {
		return block(client.NewDryRunClient(client.NewNamespacedClient(cl, "default")))
	}
	clns := client.NewNamespacedClient(cl, namespace.ObjectMeta.Name)

	if err := block(clns); err != nil {
		go func() {
			var gracePeriod int64 = 0
			err := cl.Delete(ctx, namespace, &client.DeleteOptions{GracePeriodSeconds: &gracePeriod})
			if err != nil {
				fmt.Printf("Error cleaning up namespace: %v", err)
			}
		}()
		return err
	}

	return nil
}
