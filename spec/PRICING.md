# Pricing

(Document status: wip)

A Pricing Table defines the terms under which a given provider is willing to rent its compute resources for running pods. The Pricing Table is typically represented using a Protobuf object encoding the prices for reserving or using certain resources.

| resource | typical usage | suggested/example metrics |
| -------- | ------------- | ----------------- |
| `cpu` | reserve (while application is running) | `kube_pod_resource_request{pod, resource=cpu}` |
| `memory` | reserve (while application is running) | `kube_pod_resource_request{pod, resource=memory}` |
| `storage` | reserve (incl. while application is stopped) | `kubelet_volume_stats_available_bytes{namespace}` |
| `nvidia.com/gpu(\|shared)` | reserve | `kube_pod_resource_request{pod, resource=nvidia.com/gpu(\|.shared)}` |
| `apocryph.network/ip` | reserve | `kube_service_spec_type{namespace, type=NodePort}` |
| `kubernetes.io/(in\|e)gress-bandwidth` | usage | `nginx_ingress_controller_nginx_process_(read\|write)_bytes_total` |




The pricing table must encompass all essential billing details for clients, and it might be as straightforward as this:

| Resource          | Description| Price                            |
|-------------------|------------|-----------------------------|
| **CPU**           | N° of Cores/ N° of vCPUs|$0.0001 VCPU(min/s/ms)|
| **RAM Capacity**  | Capacity (ex: GB)| $0.00001 GB/(min/s/ms) |
| **Storage**       | Type (e.g., Block, Object)|$0.00001 GB/(min/s/ms)|             |
| **GPU (Optional)**| Model|$0.0001 per Execution (min/s/ms)|

Or it could be split into categories with more detailed information:

### Compute Pricing

* **CPU**
    
    | Resource| Description | Number of Cores| vCPUs |     Model | TEE Type| Price per Unit |
    |-|-|-|-|-|-|-|
    | **CPU**   | Processing power   | Cores           | vCPUs    | Intel, AMD, ARM, ...etc      | Enclaves, CVMs, ..etc | $0.0001 VCPU(min/s/ms)    |

- **Ram**

    |Resource|Description|Capacity| Price|
    |-|-|-|-|
    | **RAM**   | Memory capacity | Capacity (ex: 1GB) | $0.00001         GB/(min/s/ms) 


### Storage Pricing

| Resource      | Description        | Capacity |     Storage Type | Price per Unit |
|---------------|--------------------|----------|-------------- |-----------------|
| **Storage** | Storage resources | Capacity (ex: 10GB) | Block, Object, ...etc | $000001 per GB(min/s/ms) |




#
References:
* [Kube Scheduler Metrics](https://kubernetes.io/docs/concepts/cluster-administration/system-metrics/#kube-scheduler-metrics), [2](https://github.com/kubernetes/kubernetes/blob/a321897e77ae43011fee55cfd22092008121ccb6/test/instrumentation/testdata/stable-metrics-list.yaml#L404-L431) - resource requests/limits
* [Kube State Metrics](https://github.com/kubernetes/kube-state-metrics/blob/main/docs/service-metrics.md) - resource definitions
* [CAdvisor Metrics](https://github.com/google/cadvisor/blob/master/docs/storage/prometheus.md) - resource usage
* [Ingress Nginx Metrics](https://github.com/kubernetes/ingress-nginx/blob/main/docs/user-guide/monitoring.md)
