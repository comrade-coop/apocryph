package kubernetes

import (
	"context"

	kedahttpscheme "github.com/kedacore/http-add-on/operator/generated/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetScheme() (*runtime.Scheme, error) {
	sch := runtime.NewScheme()
	if err := scheme.AddToScheme(sch); err != nil {
		return nil, err
	}
	if err := kedahttpscheme.AddToScheme(sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func GetNamespacedClient(ctx context.Context, config *rest.Config, dryRun bool) (client.Client, error) {
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

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{GenerateName: "tpods"},
	}
	if err := cl.Create(ctx, namespace); err != nil {
		return nil, err
	}

	if !dryRun {
		cl = client.NewNamespacedClient(cl, namespace.ObjectMeta.Name)
	} else {
		cl = client.NewNamespacedClient(cl, "default")
	}

	return cl, nil
}
