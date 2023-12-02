import { field, FieldOrRaw } from './field'
import type {
  ProviderConfig,
  Pod,
  PaymentChannelConfig
} from 'trusted-pods-proto-ts'
import type { PartialMessage } from '@bufbuild/protobuf'
import { hexToBytes } from 'viem'

export function template(): FieldOrRaw<{
  funds: bigint
  unlockTime: bigint
  payment: PartialMessage<PaymentChannelConfig>
  provider: PartialMessage<ProviderConfig>
  pod: PartialMessage<Pod>
}> {
  return {
    funds: field(10000000000000000000000n, { min: 0n }),
    unlockTime: field(240n, { min: 0n }),
    payment: {
      paymentContractAddress: hexToBytes(
        '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
      ) // TODO= result of forge create
    },
    provider: {
      ethereumAddress: field(
        hexToBytes('0x70997970C51812dc3A010C7d01b50e0d17dc79C8'),
        { encoding: 'eth-address' }
      ), // TODO= anvil.accounts[1]
      libp2pAddress: field(import.meta.env.VITE_PROVIDER_MULTIADDR ?? '') // TODO TODO - fetch from registry!
    },
    pod: {
      containers: [
        {
          name: field('nginx-hello'),
          image: {
            url: 'docker.io/nginxdemos/nginx-hello:latest'
          },
          ports: [
            {
              name: 'http',
              containerPort: 8080n,
              exposedPort: {
                case: 'hostHttpHost',
                value: field(
                  `x${Math.random()
                    .toString(16)
                    .slice(2)}.podhostname.localhost`
                )
              }
            }
          ],
          resourceRequests: [
            {
              resource: 'cpu',
              quantity: {
                case: 'amountMillis',
                value: field(100n, { min: 0n })
              }
            },
            {
              resource: 'memory',
              quantity: {
                case: 'amount',
                value: field(100000000n, { min: 0n })
              }
            }
          ]
        }
      ],
      replicas: {
        max: 1
      }
    }
  }
}
