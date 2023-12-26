# Network Protocol

(Document status: barebones)

When the Publisher Client connects to the Provider Client, it makes use of a libp2p connection; likely through the IPFS DHT, unless the Provider has advertised a stable IP earlier. These connections use the protocols `/apocryph/attest/0.0.1`, `/apocryph/provisioning-capacity/0.0.1`, and `/apocryph/provision-pod/0.0.1`, which is based on a Protobufs protocol defined in the [`:/proto/`](../proto) folder.

The basic structure of this protocol is the following:

1. The Publisher requests an attestation from the Provider using the `/apocryph/attest/0.0.1` libp2p protocol.
2. The Provider replies with an attestation (and optionally, the resource capacity available), proving that the whole Provider stack (including the endpoint of the current stream) is running inside a TEE which is trusted by the Publisher.

3. Optionally, the Publisher inquires about the resources that will be requested, using  `/apocryph/provisioning-capacity/0.0.1`. Resource requirements can include amounts of CPU cores, RAM memory, GPU presence, specific CPU models, and even certain numbers of external IPs available.
4. Optionally, the Provider replies with the resources that would be offered / that are available; along with the prices (and payment address) at which the Provider is willing to offer those.

5. The Publisher sends a message that includes the on-wire Manifest and payment channel information using `/apocryph/provision-pod/0.0.1`.
6. The Provider provisions the requested services, and replies with a status message

## Manifest wire format

When the pod manifest is transferred between the Publisher Client and Provider Client, it uses a modified version of the [usual manifest format](MANIFEST.md), based on protocol buffers, in order to reduce the ambiguity of using YAML.

Main changes:
* There are no longer multiple ways to define volumes; all volumes must be in an array at the end.
* There are no longer multiple ways to define ports; all ports use the `(port, targetPort, protocol, hostIP)` fields.
* Likewise for the `command` and `args` fields.
<!--* The wire format includes fields for pricing copied from the [registry data](REGISTRY.md).-->

Finally, the `image` field is possibly converted to an IPFS hash and key, containing respectively the output of [`imgcrypt`](https://github.com/containerd/imgcrypt) uploaded with [`ipdr`](https://github.com/ipdr/ipdr) and the private key that can be used to decrypt it.
