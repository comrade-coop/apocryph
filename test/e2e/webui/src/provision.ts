// SPDX-License-Identifier: GPL-3.0

import './style.css'
import {
  bytesToHex,
  hexToBytes,
  numberToBytes,
  stringToHex,
  Hex
} from 'viem'
import { heliaNodePromise, walletClient } from './connections'
import { connectTo } from 'apocryph-ipfs-ts'
import { multiaddr } from '@multiformats/multiaddr'
import { createPromiseClient } from '@connectrpc/connect'
import {
  Pod,
  ProvisionPodResponse,
  ProvisionPodService,
  provisionPodProtocolName
} from 'apocryph-proto-ts'
import { PartialMessage } from '@bufbuild/protobuf'

// Copied and adapted from pkg/publisher/connect.go

/**
 * Provision a pod on the configured provider. This function currently only does the initial provisioning of a pod, and will not update a pod if one already exists.
 *
 * @param config deployment configuration; includes payment channel details as well as provider details
 * @param config.pod the pod configuration to provision on the provider
 */
export async function provisionPod(config: {
  payment: {
    paymentContractAddress: Uint8Array
    podID: Uint8Array
  }
  provider: {
    ethereumAddress: Uint8Array
    libp2pAddress: string
  }
  pod: PartialMessage<Pod>
}): Promise<ProvisionPodResponse> {
  const heliaNode = await heliaNodePromise
  const connection = connectTo(
    heliaNode,
    multiaddr(config.provider.libp2pAddress),
    provisionPodProtocolName
  )
  const client = createPromiseClient(ProvisionPodService, connection)

  const tokenData = JSON.stringify({
    PodId: bytesToHex(config.payment.podID),
    Operation:
      '/' +
      ProvisionPodService.typeName +
      '/' +
      ProvisionPodService.methods.provisionPod.name,
    ExpirationTime: new Date(Date.now() + 1000 * 10),
    Publisher: walletClient.account.address
  })

  const signature = (await walletClient.account.signTransaction(
    {},
    {
      serializer() {
        return stringToHex(tokenData)
      }
    }
  )) as Hex // FIXME: HACK: We should just use EIP typed data signatures and be done with it...

  const tokenDataEncoded = btoa(tokenData)
  const signatureEncoded = btoa(String.fromCodePoint(...hexToBytes(signature)))
  const bearerToken = tokenDataEncoded + "." + signatureEncoded

  const result = await client.provisionPod(
    {
      pod: config.pod,
      payment: {
        chainID: numberToBytes(walletClient.chain.id),
        providerAddress: config.provider.ethereumAddress,
        contractAddress: config.payment.paymentContractAddress,
        publisherAddress: hexToBytes(walletClient.account.address),
        podID: config.payment.podID
      }
    },
    {
      headers: {
        authorization: "Bearer " + bearerToken,
      }
    }
  )
  if (result.error !== '') {
    throw new Error(`Error from provider: ${result.error}`)
  }
  return result
}
