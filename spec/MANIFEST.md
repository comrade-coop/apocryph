# Application Package

(Document status: work-in-progress)

## Manifest file

The manifest file, tentatively `trustedpods.yml`, is read by the Publisher Client and used to assemble the on-wire manifest that is then sent to the Provider Client. What follows is an example manifest file as it was initially conceptualized; the current implementation simply reads off a yaml-encoded protobuf object using the on-wire manifest format and does not follow the format below. (Use `trustedpods init` to get a sample file generated.)

```yaml
# WARNING: Pseudo-code
type: "trustedpods"
version: "1.0"
containers:
- image: localregistryname:tag
  command: override command # or command: ["override", "command"] # $(VAR_NAME) as in K8s -- ENTRYPOINT
  args: override args # or args: ["override", "args"] # $(VAR_NAME) as in K8s -- CMD
  workingDir: /override/pwd/ # as in K8s
  port: 80 # HTTP port (must have only one per pod)
  host: example.com # HTTP hostname used to route requests to the container (must have only one per pod)
  ports:
  - 123:321 # port mapping, as in docker-compose
  - 123 # port mapping, as in docker-compose
  - port: 123 # as in K8s services
    targetPort: 321 # as in K8s services
    protocol: TCP # or UDP, as in K8s
    hostIP: false # request that the port be exposed to the external world; otherwise it will be accessible only using k8s DNS
  env:
  - name: XX
    value: VAL
  volumes:
  - mountPath: /vol # as in K8s
    name: vol # alternatively - without name, copy the same fields from the volume definition here.
    readOnly: false
  resources:
    cpu: 1000m # in milliCPU; equivalent to K8s Requests
    memory: 1Gi
    nvidia.com/gpu: 1
replicas:
  min: 0
  max: 1
volumes:
- name: vol
  type: volume # or emptyDir or secret
  resources: # for type: volume
    storage: 8Gi
  source: ./publisher/local/file.json # for type: secret
```

Note that the Publisher Client might eventually gain functions for reading other kinds of manifests, such as for directly consuming docker-compose files.

References:
* [`HTTPScaledObject`](https://github.com/kedacore/http-add-on/blob/main/docs/ref/v0.3.0/http_scaled_object.md)
* [`DeploymentSpec`](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/#DeploymentSpec)
* [`PodTemplateSpec`](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-template-v1/#PodTemplateSpec)
* [`PodSpec`](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#PodSpec)
* [Pod `resources:`](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)
* [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
* [`docker-compose`](https://docs.docker.com/compose/compose-file/compose-file-v3/)
