# Backlog

This file seeks to document the tasks left to do in this repository, as well as design flaws and accumulated tech debt present in the current implementation.

Rationale for not shoving all of this in GitHub Issues: while Issues are a great way for users and contributors to voice concerns and problems, for huge planned milestones, it feels simpler to just have them listed in a format more prone to conveying information and not to discussing the minutiae details.

## Features yet to be implemented

- Integrate attestation
- Support private docker registries
- Registry support in web frontend

## Technical debt accumulated

### Prometheus metrics used for billing

Status: Alternative prototyped

Prometheus's [documentation](https://prometheus.io/docs/introduction/overview/#when-does-it-not-fit) explicitly states that Prometheus metrics are not suitable for billing since Prometheus is designed for availability (an AP system by the [CAP theorem](https://en.wikipedia.org/wiki/CAP_theorem)) and not for consistency / reliability of results (which is what a CP system would be).

Despite that, the current [implementation](../pkg/prometheus/) uses Prometheus and `kube-state-metrics` to fetch the list of metrics that billing is based on. A prototype was created in the [`metrics-monitor`](github.com/comrade-coop/apocryph/tree/metrics-monitor) branch to showcases an alternative way to fetch the same data from Kubernetes directly and avoid any possible inconsistencies in the result, yet it was decided that it's better to iterate quickly with Prometheus first instead and come back to this idea later.

### A single monolithic `tpodserver` service

Status: Correct as needed

Currently, the whole of the Apocryph' Provider client/node is implemented as a pair of long-running processes deployed within Kubernetes -- one for listening for incoming Pod deployments and one for monitoring them. Going forward, it could be beneficial to make more parts of that service reusable by splitting off libp2p connections, actual deployments, metrics collection, and smart contract invoicing into their own processes/services can be changed or reused on their own.

### Payment contract is one contract and not multiple

Status: Still evaluating, alternative prototyped

The [payment contract](../contracts/src/Payment.sol) currently takes care of absolutely all payments that pass through Apocryph. However, it might be worth splitting it into a factory/library contract and small "flyweight" contracts instead. Currently, that is prototyped in the [`contract-factory`](github.com/comrade-coop/apocryph/tree/contract-factory) branch, but it ended up using more way more gas for deployment, so it was temporarily scrapped.

### Using Kubo/IPFS p2p feature marked experimental

Status: Requires research

Kubo's [`p2p` API](https://docs.ipfs.tech/reference/kubo/rpc/#api-v0-p2p-forward) is marked as an [experimental feature](https://github.com/ipfs/kubo/issues/3994), and is predictably rather finicky to work with. Moreover, it may very well be removed one day, with or without alternative, as is happening with the [`pubsub` feature](https://github.com/ipfs/kubo/issues/9717).

As such, it would be prudent to move away from using the `p2p` features of Kubo (and away from requiring Kubo-based IPFS nodes), and instead roll out an alternative, likely based on `libp2p`. This will likely be easier once [the planned Amino/DHT refactor](https://blog.ipfs.tech/2023-09-amino-refactoring/) lands.

#### `ipfs-p2p-helper` is a sidecar

Status: Correct as needed

Currently, the `ipfs-p2p-helper`, a small piece of code responsible for registering `p2p` listeners in Kubo. Doing so is a bit tricky, as the Kubo daemon does not persist `p2p` connections between restarts, and hence we have to re-register them every time the IPFS container restarts.

This is currently done using a sidecar container (a container in the same pod), so the helper gets restarted together with IPFS -- and to top that off, it just watches the list of Services for ones that are labeled correctly. Ideally, if we keep using the `p2p` feature of Kubo, we would rewrite `ipfs-p2p-helper` to be a "proper" Kubernetes operator with a "proper" custom resource definition.

### Custom HTTP client implementation in web frontend

Status: Correct as needed

Currently, the [Libp2p Connect transport](../pkg/ipfs-ts/transport-libp2p-connect.ts) implemented in the repo ends up reimplementing a whole HTTP client, just for the sake of sending [ConnectRPC](https://github.com/connectrpc/connect-es) messages over a [libp2p connection](https://libp2p.github.io/js-libp2p/interfaces/_libp2p_interface.connection.Stream.html). This is not ideal, as HTTP clients are notoriously complicated to implement right, and while it's unlikely that ours is rifle with vulnerabilities, it's also unlikely that implementing one ourselves is the best way forward.

The two main options here would be to either drop ConnectRPC completely and implement framing ourselves (and thus reimplementing ConnectRPC/GRPC while still using Protobufs for the message serialization itself) or to use an existing implementation of the HTTP client, such as node's HTTP package. Alternatively, if we use the Kubo/IPFS p2p feature instead of importing libp2p into the browser, we might be able to directly use ConnectRPC with the correct port numbers, at the cost of losing encryption and (currently) authenticity of the requests, unless the user is running their own Kubo node.

### Secret encryption done with AESGCM directly

Status: Correct as needed

Currently, we encrypt secrets' data ([(see `EncryptWith`/`DecryptWith`)](../pkg/crypto/key_management.go)) with AESGCM directly, forgoing using any libraries that could do this for us and give us a more generic encrypted package. Ideally, given that the rest of the code uses `go-jose` we would use `go-jose`'s encryption facilities directly -- however, JWE objects base64-encode the whole ciphertext... making them ~33% less efficient in terms of space on-wire! Hence, we opt to directly write the bytes ourselves and save on some space.

Some ways to improve the situation would be to contribute `BSON` functionality to `go-jose` (unfortunately, such functionality would not be standards-compliant, unless someone goes the whole way to suggest `BSON` (or other binary) serialization for [RFC7516](https://www.rfc-editor.org/rfc/rfc7516.html)), to switch to using PKCS11 instead of JSON Web Keys, or implementing our own key provider for `ocicrypt` (which was the reason to start using JSON Web Keys in the first place), perhaps one based on [ERC-5630](https://eips.ethereum.org/EIPS/eip-5630). Alternatively, we could look into other standards for storing encrypted secrets, such as [IPFS/Ceramic's dag-jose](https://github.com/ceramicnetwork/js-dag-jose/) or [WNFS](https://github.com/wnfs-wg/) or any of the [other nascent IPFS encryption standards](https://discuss.ipfs.tech/t/encryption-private-data-and-private-swarms-with-ipfs/15363).

### Code duplication in cmd/trustedpods

Status: Correct as needed

A lot of the code in [`cmd/trustedpods`](../cmd/trustedpods) has to do with setting up the environment for things like Ethereum connections, IPFS connections, deployment files, provider addresses, etc., and even IPFS uploads are in a sense a dependency of sending requests to the provider. It would be nice if we could express everything as a pipeline of dependencies that each inserts its own flags into the command parser and then gets processed in turn so as to create the whole desired environment in the end.

This has been attempted in the past (outside of Git history), but the result was even less managable. Perhaps, this is something `cobra` is not particularly well suited for, and an additional (homegrown?) dependency management system would help. Either way, the code duplication is not horrible, and the repo will survive as-is for a long time before it becomes problematic.

## Missing features

### Constellation cluster recovery not handled

Status: Solutions outlined

Constellation, the confidential Kubernetes solution we have opted to use, works by bootstrapping additional nodes on top of an existing cluster through their JoinService -- whereby a new node asks the old node's Join service for the keys used for encrypting Kubernetes state, while the old node confirms the new node's attestation through aTLS. This makes it excellent for autoscaling scenarios; however, in the case a full-cluster power outage occurs, it leaves the cluster in an hung state, as there is no initial node to bootstrap off of, and requires manually re-attesting the cluster and inputting the key that was backed up when provisioning the cluster initially -- as documented in [the recovery procedure documentation](https://docs.edgeless.systems/constellation/workflows/recovery)

For Apocryph, however, we cannot trust the provider with a key that decrypts the whole state of the cluster - as that will destroy the confidentiality of the pods running within Apocryph. Hence, when recovering an existing cluster, or when initially provisioning a cluster, we would need a securely-stored key that can only be accessed from an attested TEE that is part of the cluster.

There are multiple ways to do so. A simple one would be to generate and store the key within a TPM, and making sure the TPM only reveals the key to the attested TEE; this still leaves attesting that the key is generated there as an open task. Another one would be to modify Constellation to allow for the master secret to be stored encrypted with the TEE's own key (inasmuch as one exists) - so that the same machine, when rebooted, can be bootstrapped on its own. And finally, a more involved solution would be to use [Integretee](https://www.integritee.network/) or an equivalent thereof to generate and store cluster keys in a cloud of working attested enclaves.

### Apocryph cluster attestation

Status: Correct as needed

Constellation allows [attesting a cluster](https://docs.edgeless.systems/constellation/workflows/verify-cluster).. however, upon closer inspection, the attestation features provided only allow attesting that the whole machine is running a real Constellation cluster in a real TEE enclave... and say nothing about the containers running inside that cluster. This is only fair, perhaps, given that the containers can be configured in ways that could allow them to escape the confines of their sandboxes; however, it does mean that attestation, if implemented, will not be sufficient to convince the publisher the peer they are talking to is a Apocryph node.

The main solution to this, other than switching away from Constellation (to, e.g. Confidential Containers, despite them not being fully ready yet), would be to modify the base Constellation image so that it includes an additional API, either running within or without a container, whose hash is verified in the boot process, and which allows querying, and hence, attesting the rest of the Kubernetes state. Alternatively, the image could be modified to attest the Apocryph server container as part of the boot process; however, this feels like too much hardcoding.

### Apocryph cluster hardening

Status: Known issue

In line with the two notes about Constellation's cluster recovery and attestation features, a third departure of a Apocryph cluster from what Constellation provides out of the box is the fact that Constellation issues an admin-level Kubectl access token upon installation; however, we would like to keep parts of the Apocryph cluster inaccessible even to the administrator.

For that, we would likely need to issue a Kubectl access token with lesser privileges, allowing for only partial configuration of the Apocryph cluster. The customizable features should be selected carefully to align with Provider needs, to allow for things like configuring backups and some kinds of dashboards and monitoring, while minimizing the leaking of user privacy.

### Storage reliability

See [the respective document](STORAGE.md) for an in-depth storage reliability design proposal.

### Uptime reliability

See [the respective document](UPTIME.md) for a more in-depth uptime reliability design proposal.

### Software licensing

See [the respective document](MARKETPLACE.md) for a more in-depth software licensing design proposal.

### Individual TEEs

Status: Needs more usecases

Currently, the architecture of Constellation uses a single TEE encompassing all Kubernetes pods running in the cluster. However, for extra isolation of individual tenants, it could be beneficial to have separate TEEs for each publisher / pod / application. To implement that, we will likely end up scrapping Constellation and revamping the whole attestation process. As this is quite a bit of design and implementation work while the gains at this stage are minimal, we have opted to let the idea ruminate for the moment.

### Forced scale-down

Status: Conceptualized

It would be great if we didn't just rely on KEDA's built-in scaling down after a certain time, but also allowed Pods to request their own scaling down. See also [this issue](https://github.com/kedacore/http-add-on/issues/840).
