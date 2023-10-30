# Backlog

This file seeks to document the tasks left to do in this repository, as well as design flaws and accumulated tech debt present in the current implementation. 

Rationale for not shoving all of this in GitHub Issues: while Issues are a great way for users and contributors to voice concerns and problems, for huge planned milestones, it feels simpler to just have them listed in a format more prone to conveying information and not to discussing the minutiae details.

## Features yet to be implemented

* Write and integrate a Registry contract
* Run the Provider in Constellation instead of Minikube, to allow for attestation
* Implement a way for the Publisher to monitor, manage, and edit the deployed pod (other than shutting down the payment channel)

## Technical debt accumulated

### Prometheus metrics used for billing

Status: Alternative prototyped

Prometheus's [documentation](https://prometheus.io/docs/introduction/overview/#when-does-it-not-fit) explicitly states that Prometheus metrics are not suitable for billing since Prometheus is designed for availability (an AP system by the [CAP theorem](https://en.wikipedia.org/wiki/CAP_theorem)) and not for consistency / reliability of results (which is what a CP system would be).

Despite that, the current [implementation](../pkg/prometheus/) uses Prometheus and `kube-state-metrics` to fetch the list of metrics that billing is based on. A prototype was created in the [`metrics-monitor`](https://github.com/comrade-coop/trusted-pods/tree/metrics-monitor) branch to showcases an alternative way to fetch the same data from Kubernetes directly and avoid any possible inconsistencies in the result, yet it was decided that it's better to iterate quickly with Prometheus first instead and come back to this idea later.

### A single monolithic `tpodserver` service

Status: Correct as needed

Currently, the whole of the Trusted Pods' Provider client/node is implemented as a pair of long-running processes deployed within Kubernetes -- one for listening for incomming Pod deployments and one for monitoring them. Going forward, it could be beneficial to make more parts of that service reusable by splitting off libp2p connections, actual deployments, metrics collection, and smart contract invoicing into their own processes/services can be changed or reused on their own.

### Payment contract is one contract and not multiple

Status: Still evaluating, alternative prototyped

The [payment contract](../contracts/src/Payment.sol) currently takes care of absolutely all payments that pass through Trusted Pods. However, it might be worth splitting it into a factory/library contract and small "flyweight" contracts instead. Currently, that is prototyped in the [`contract-factory`](https://github.com/comrade-coop/trusted-pods/tree/contract-factory) branch, but it ended up using more way more gas for deployment, so it was temporarily scrapped.

### Using Kubo/IPFS p2p feature marked experimental

Status: Requires research

Kubo's [`p2p` API](https://docs.ipfs.tech/reference/kubo/rpc/#api-v0-p2p-forward) is marked as an [experimental feature](https://github.com/ipfs/kubo/issues/3994), and is predictably rather finicky to work with. Moreover, it may very well be removed one day, with or without alternative, as is happening with the [`pubsub` feature](https://github.com/ipfs/kubo/issues/9717).

As such, it would be prudent to move away from using the `p2p` features of Kubo (and away from requiring Kubo-based IPFS nodes), and instead roll out an alternative, likely based on `libp2p`. This will likely be easier once [the planned Amino/DHT refactor](https://blog.ipfs.tech/2023-09-amino-refactoring/) lands.

#### `ipfs-p2p-helper` is a sidecar

Status: Correct as needed

Currently, the `ipfs-p2p-helper`, a small piece of code responsible for registering `p2p` listeners in Kubo. Doing so is a bit tricky, as the Kubo daemon does not persist `p2p` connections between restarts, and hence we have to re-register them every time the IPFS container restarts.

This is currently done using a sidecar container (a container in the same pod), so the helper gets restarted together with IPFS -- and to top that off, it just watches the list of Services for ones that are labeled correctly. Ideally, if we keep using the `p2p` feature of Kubo, we would rewrite `ipfs-p2p-helper` to be a "proper" Kubernetes operator with a "proper" custom resource definition.

### Reuploading IPDR images

Status: Correct as needed, upstream available but needs bugfixing

Currently, the [code (see `ReuploadImagesFromIpdr`)](../pkg/ipfs/images.go) dealing with transforming images that have been uploaded as IPDR takes those same images and uploads them to a local registry. Ideally, what would happen instead is that IPDR images would instead be treated as first-class citizens and downloaded on-demand (probably with some prefetching to reduce first-boot time).

There are a couple ways to implement that. One would be to run an IPDR registry in the cluster and fetch images from it. Unfortunatelly, as the [relevant issue in ipdr/ipdr notes](https://github.com/ipdr/ipdr/issues/18), the IPDR's code currently (flawedly) assumes CIDv1 multihashes are CIDv0 -- and as a whole, the `ipdr/ipdr` repository is outdated (checked 2023-10-27) and full of code which is not making use of the Go IPFS libraries nor of the OCI image-handling libraries -- making depeding on that library an overall increase of tech debt.

Another way to implement first-class IPDR images would be to develop a `containerd` [plugin](https://github.com/containerd/containerd/blob/main/docs/PLUGINS.md) which handles image downloads using our (surprisingly functional, considering the code size) [IPDR transport](../pkg/ipdr) -- or better yet, getting IPDR support merged into mainline `containerd`. A potential hurdle to actually doing that is that Constellation has hardcoded their [`containerd` config](https://github.com/edgelesssys/constellation/blob/main/image/base/mkosi.skeleton/usr/etc/containerd/config.toml) as part of the base layer that is later attested to.

### Secret encryption done with AESGCM directly

Status: Correct as needed

Currently, we encrypt secrets' data ([(see `EncryptWith`/`DecryptWith`)](../pkg/crypto/key_management.go)) with AESGCM directly, forgoing using any libraries that could do this for us and give us a more generic encrypted package. Ideally, given that the rest of the code uses `go-jose` we would use `go-jose`'s encryption facilities directly -- however, JWE objects base64-encode the whole ciphertext... making them ~33% less efficient in terms of space on-wire! Hence, we opt to directly write the bytes ourselves and save on some space.

Some ways to improve the situation would be to contribute `BSON` functionallity to `go-jose` (unfortunatelly, such functionallity would not be standards-compliant, unless someone goes the whole way to suggest `BSON` (or other binary) serialization for [RFC7516](https://www.rfc-editor.org/rfc/rfc7516.html)), to switch to using PKCS11 instead of JSON Web Keys, or implementing our own key provider for `ocicrypt` (which was the reason to start using JSON Web Keys in the first place), perhaps one based on [ERC-5630](https://eips.ethereum.org/EIPS/eip-5630).

## Missing features

### Storage reliability

Status: Barebones design done

For Trusted Pods to be a viable platform for deploying mission-critical software to, it needs to provide some guarantees about the longevity of data stored on it. This is still up to be designed and incorporated into the overall design and architecture.

The best solution for this we have found is establishing contracts and a protocol for Providers to stake that they will keep specific stored data available. In theory, to identify the data, a version identifier similar to ZFS's root checksum/version should be sufficient, with the in-TEE encryption of the volume providing the confidentiality of the data - at which point a container running with the same volume attached can run a full filesystem check and confirm the root checksum/version. If the Provider fails to execute the filesystem check within some (generous) time limit in response to a challenge by the Publisher, the storage can not be shown to still exist, and therefore the stake is lost and can be transferred to the Publisher.

### Uptime reliability

Status: Needs more design work

A closely-related concept to storage reliability is that of uptime reliability. While scale-down-to-zero pods might never have great first response latency, even for them Publishers should be able to ask for reliable uptime, a guarantee that the system will remain up an running.

Here again, the best system seems to be staking -- and is similar to what the cloud industry is already doing with [SLAs](https://en.wikipedia.org/wiki/Service-level_agreement). It is, however, a bit trickier to prove that the application was down due to a fault of the Provider, and there are different kinds of downtime, such as network outages and power outages that we might or might not want to handle differently.

### Software licensing

Status: Needs more design work

A useful feature would be to allow publishers to upload application code that other parties can then deploy to a provider -- this would enable usecases such as providing a pay-for-usage subscription to a service where the publisher does not have to take on the risk and responsibility of running a master instance of that service.

This should already be achievable in theory by having a master service which confirms the attestation and licensing status of deployed pods and only then provides them with a decryption key for the rest of the application. However, if it is integrated with the rest of the Trusted Pods platform, we should be able to achieve faster scale-up-from-zero performance, as well as lower some of the counter-party risk of the master service going down.

### Individual TEEs

Status: Needs more usecases

Currently, the architecture of Constellation uses a single TEE encompassing all Kubernetes pods running in the cluster. However, for extra isolation of individual tenants, it could be beneficial to have separate TEEs for each publisher / pod / application. To implement that, we will likely end up scrapping Constellation and revamping the whole attestation process. As this is quite a bit of design and implementation work while the gains at this stage are minimal, we have opted to let the idea ruminate for the moment.

### Forced scale-down

Status: Conceptualized

It would be great if we didn't just rely on KEDA's built-in scaling down after a certain time, but also allowed Pods to request their own scaling down.
