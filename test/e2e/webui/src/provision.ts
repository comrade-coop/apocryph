import './style.css'
import { bytesToHex, hexToBytes, numberToBytes, stringToHex, type Signature, type TransactionSerializable, type Hex, keccak256, concatBytes } from 'viem'
import { heliaNodePromise, walletClient } from './connections'
import { connectTo } from 'trusted-pods-ipfs-ts'
import { multiaddr } from '@multiformats/multiaddr'
import { createPromiseClient } from '@connectrpc/connect'
import base32 from 'base32'
import { type Pod, type ProvisionPodResponse, ProvisionPodService, provisionPodProtocolName } from 'trusted-pods-proto-ts'
import { type PartialMessage } from '@bufbuild/protobuf'

// Copied and adapted from pkg/publisher/connect.go
export async function provisionPod (config: {
  payment: {
    paymentContractAddress: Uint8Array
  }
  provider: {
    ethereumAddress: Uint8Array
    libp2pAddress: string
  }
  pod: PartialMessage<Pod>
}): Promise<ProvisionPodResponse> {
  const heliaNode = await heliaNodePromise
  const connection = connectTo(heliaNode, multiaddr(config.provider.libp2pAddress), provisionPodProtocolName)
  const client = createPromiseClient(ProvisionPodService, connection)

  const token = JSON.stringify({
    PodId: bytesToHex(config.payment.podID),
    Operation: '/' + ProvisionPodService.typeName + '/' + ProvisionPodService.methods.provisionPod.name,
    ExpirationTime: new Date(Date.now() + 1000 * 10),
    Publisher: walletClient.account.address
  })

  const signature = walletClient.account.signTransaction({}, {
    serializer (tx: TransactionSerializable, signature?: Signature) {
      if (signature != null) {
        return signature
      }
      return stringToHex(token)
    }
  }) as Hex // FIXME: HACK: We should just use EIP typed data signatures and be done with it...

  const namespacePartsHash = keccak256(concatBytes([hexToBytes(walletClient.account.address), config.payment.podID]), 'bytes')
  const ns = 'tpod-' + (base32.encode(namespacePartsHash) as string).toLowerCase().replace(/=+$/, '') // Why, oh why..

  const result = await client.provisionPod({
    pod: config.pod,
    payment: {
      chainID: numberToBytes(walletClient.chain.id),
      providerAddress: config.provider.ethereumAddress,
      contractAddress: config.payment.paymentContractAddress,
      publisherAddress: hexToBytes(walletClient.account.address),
      podID: config.payment.podID
    }
  }, {
    headers: {
      token,
      authorization: signature,
      namespace: ns
    }
  })
  if (result.error !== '') {
    throw new Error(`Error from provider: ${result.error}`)
  }
  return result
}
