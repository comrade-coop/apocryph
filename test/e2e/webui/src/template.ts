import { field, type FieldOrRaw } from './field'
import type { ProviderConfig, Pod } from 'trusted-pods-proto-ts'
import type { PartialMessage } from "@bufbuild/protobuf";

export function template(): FieldOrRaw<{provider: PartialMessage<ProviderConfig>, pod: PartialMessage<Pod>}> {
  return {
    provider: {
      "ethereumAddress": field(new Uint8Array(), {encoding: 'eth-address'}),
      "libp2pAddress": field("")
    },
    pod: {
      "containers": [
        {
          "name": field("nginx-hello"),
          "image": {
            "url": "docker.io/nginxdemos/nginx-hello:latest"
          },
          "ports": [
            {
              name: "http",
              containerPort: 8080n,
              exposedPort: {
                case: "hostHttpHost",
                value: field("pod.172b786856847582.hostname.example")
              }
            }
          ],
          "resourceRequests": [
            {
              "resource": "cpu",
              quantity: {
                case: 'amountMillis',
                value: field(100n, {min: 0n})
              }
            },
            {
              "resource": "memory",
              quantity: {
                case: 'amount',
                value: field(100000000n, {min: 0n})
              }
            }
          ]
        }
      ],
      "replicas": {
        "max": 1
      }
    }
  }
}
