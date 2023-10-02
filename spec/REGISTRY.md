# Registry

(Document status: barebones)

In order for the Publisher Client to connect to a Provider, it first needs to know how to contact the Provider. The main way through which this is accomplished is through the Registry. The Registry stores a list of Providers.

Each Provider stored in the Registry is stored as an IPFS Hash pointing to a Protobufs-encoded or DAG-CBOR-encoded blob, along with the on-chain address of the Provider (for on-chain Registries).

The data contained for each provider is as following:
* A list of multiaddrs at which the Provider can be found - required
* A resource capacity / pricing table, similar/equivalent to the one used in the [Network Protocol](PROTOCOL.md) - optional
* A possibly-empty list of attestations for the Provider; whether or not any of those are accepted depends on the Publisher - optional

The attestations and pricing may be empty/missing; but in that case, Publishers are less likely to query the Provider to fill in the blanks.
