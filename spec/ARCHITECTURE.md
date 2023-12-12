# Architecture

(Document status: mostly complete)

Much of the architecture of Trusted Pods revolves around two key pieces, the "Publisher Client" and the "Provider Client". A publisher is a buyer in the Trusted Pods network seeking to provision their pod/container on the network. A provider is a seller seeking to offer their hardware for rent.

The planned version of Trusted Pods takes care of matching the two and provisioning the pod on a specific target provider. While it offers a way to manually pick a provider, it defaults to automatically picking one for the user. However, it does not take care of operational issues that might subsequently arise; specifically, it does not take care of rescheduling pods when a provider becomes unavailable nor does it make any guarantees about uptime or availability of data—those will be handled in later versions. Such concerns are kept track of in the [backlog](BACKLOG.md).

## Bird's-eye view

```mermaid
flowchart
  classDef Go stroke:#0ff
  classDef Chain stroke:#f0f
  classDef Lib stroke:#000
  classDef User stroke:#ff0

  subgraph Publisher machine
    RegistryPub[K8s Registry]:::Lib
    PubCl[Publisher Client]:::Go
  end
  External[External Clients]:::User
  Network[libp2p/IPFS]:::Lib
  subgraph TEE
    ProCl[Provider Client]:::Go
    KubeCtl[K8s Control Plane]:::Lib
    RegistryPro[K8s Registry]:::Lib
    Http[HTTP Facade]:::Lib
    App[Application Pod]:::User
  end
  subgraph On-chain
    RegistryC[Registry Contract]:::Chain
    PaymentC[Payment Contract]:::Chain
  end

  %% KubeCtl -- Execute --> ProCl
  
  RegistryC -- List Providers --> PubCl -- Create --> PaymentC -- Monitor --> ProCl
  RegistryPub -- App Container --> PubCl
  PubCl -- Encrypted App Container --> Network -- Execution request --> ProCl
  ProCl -- App Container --> RegistryPro --> App
  ProCl -- Configuration --> KubeCtl
  KubeCtl -- Execute --> Http -- Metrics --> KubeCtl -- Execute --> App
  External -- Requests --> Http -- Requests --> App
```

As a sequence of steps:

0. The user starts the [Publisher Client](../cmd/trustedpods/) to deploy a container
1. The Publisher Client collects the [Pod Manifest](MANIFEST.md)
2. The Publisher Client gets the list of providers from the [Registry Contract](REGISTRY.md)
3. The Publisher Client selects a provider, using the configured strategy (automatic or by asking the user to manually make a choice)
4. The Publisher Client creates a [Payment Contract](../) and transfers the initial payment amount (possibly in parallel with steps 5-6)
5. The Publisher Client bundles up the Pod Manifest, any related resources, and the Payment Contract's address and sends them to the Provider Client over the [Network Protocol](PROTOCOL.md), encrypted
6. The Provider Client creates the relevant configurations for the Pod using the [Kubernetes API](https://kubernetes.io/docs/reference/kubernetes-api/), including an HTTP Scaler, an Application Pod, and Monitoring
7. The Provider Client confirms receiving the manifest and resources
8. When HTTP requests come in, the HTTP Scaler contacts the Kubernetes API in order to scale the Application Pod up
9. The Application Pod is started using the configuration from earlier
10. The Application Pod handles the incoming requests
11. After a period of no requests the Scaler uses the Kubernetes API to scale the Application Pod down
12. The Monitoring component keeps track of how many e.g. CPU-seconds the Application Pod has run for, and forwards these metrics to the Provider Client
13. The Provider Client submits the metrics to the Payment Contract and is then able to claim the payment due
14. Whenever the Payment Contract runs out of funds, the Provider Client removes the related configurations from Kubernetes

<details><summary>Sequence Diagram</summary>

```mermaid
sequenceDiagram
  box
  actor User as User
  participant PubCl as Publisher Client
  end
  box
  participant RegistryC as Registry Contract
  participant PaymentC as Payment Contract
  participant IPFS as IPFS Network
  end
  box
  participant ProCl as Provider Client
  participant RegistryD as Docker Registry
  participant K8s as Kubernetes API
  participant Monitoring as Monitoring
  participant HTTP as HTTP Scaler
  participant App as Application Pod
  end
  %% 0.
  %% 1.
  User ->>+ PubCl: Pod Manifest
  User ->> PubCl: Pod Container
  %% 2.
  PubCl -->>+ RegistryC: List Providers
  RegistryC ->>- PubCl: 
  %% 3.
  PubCl ->> PubCl: Select Provider
  %% 4.
    %%PubCl ->>+ ProCl: Initial execution request
    %%ProCl -->>- PubCl: Confirmation
  PubCl ->> PaymentC: Create & Configure
  %% 5.
  PubCl ->>+ IPFS: Upload Encrypted Container
  PubCl ->>+ ProCl: Execution request
  %% 6.
  ProCl -->> IPFS: Download Container
  IPFS ->>- ProCl: 
  ProCl ->>+ PaymentC: Monitor
  ProCl ->> K8s: Configure Application & Scaler
  K8s ->>+ HTTP: Start Intercepting
  %% 7.
  ProCl -->>- PubCl: Confirmation
  PubCl -->>- User: 

  %% 8.
  User ->>+ HTTP: HTTP Request
  HTTP ->>+ K8s: Scale up
  %% 9.
  K8s ->>- App: Start w/ Secrets
  K8s ->>+ Monitoring: App up
  App ->>+ RegistryD: Fetch Image
  RegistryD ->>+ IPFS: 
  IPFS ->>- RegistryD: OCI Image
  RegistryD ->>- App: 
  App ->> App: Decrypt Image
  %% 10.
  HTTP ->>+ App: Request
  App ->>- HTTP: Response
  HTTP ->>- User: 
  %% 11.
  note over HTTP: Time passes
  HTTP ->>+ K8s: Scale down
  K8s -x- App: Stop
  Monitoring ->>- K8s: App down

  %% 12.
  ProCl ->>+ Monitoring: Get Metrics
  Monitoring ->>- ProCl: 
  %% 13.
  ProCl ->>+ PaymentC: Submit Metrics
  PaymentC -->>- ProCl: 

  %% 14.
  User ->>+ PaymentC: Unlock Funds
  PaymentC -->> User: 
  PaymentC -->>+ ProCl: 
  deactivate PaymentC
  opt insufficient funds left
  ProCl ->>- K8s: Remove Configurations
  K8s ->> HTTP: Stop intercepting
  deactivate HTTP
  end
  deactivate PaymentC
```

</details>
