// import protos from 'gen/protos.json'
// import type { ProvisionPodRequest } from 'gen/apocryph/proto/v0/provisionPod/ProvisionPodRequest.ts'
// import { create } from 'kubo-rpc-client'
// import { Client as KuboClient } from 'kubo-rpc-client/src/lib/core'
// import { createPeerId } from '@libp2p/peer-id'
// import type { Connection } from '@libp2p/interface/connection'
// import { decode as multihashDecode } from 'multiformats/hashes/digest'
// import { fromString } from 'multiformats/bytes'

// var kuboOptions = {}
// export var kubo = create(kuboOptions)
// export var kuboClient = new KuboClient(kuboOptions)

import isPrivate from 'private-ip'
import { createHelia } from 'helia'
// import { createPeerId } from '@libp2p/peer-id'
// import type { Connection } from '@libp2p/interface/connection'
import { multiaddr, type Multiaddr, type MultiaddrInput } from '@multiformats/multiaddr'
import { createPromiseClient, type Transport } from '@connectrpc/connect'
import { createLibp2pConnectTransport } from './transport-libp2p-connect'
import { ProvisionPodService } from './generated/provision-pod_connect'

const allowedMultiaddrs = new Set<string>()
const onlyAllowed = true

export const helia = await createHelia({
  libp2p: {
    connectionGater: {
      denyDialMultiaddr: async (ma: Multiaddr): Promise<boolean> => {
        if (allowedMultiaddrs.has(ma.toString())) {
          return false
        } else if (onlyAllowed) {
          return true
        }
        const tuples = ma.stringTuples()

        if (tuples[0][0] === 4 || tuples[0][0] === 41) {
          return Boolean(isPrivate(`${tuples[0][1]}`))
        }

        return false
      }
    },
    peerDiscovery: [],
    start: false
  }
})

export function connectTo (peerAddr: MultiaddrInput, protocol: string = '/x/trusted-pods/provision-pod/0.0.1'): Transport {
  const peerMultiaddr = multiaddr(peerAddr)
  allowedMultiaddrs.add(peerMultiaddr.toString())

  return createLibp2pConnectTransport({
    dialStream: async () => await helia.libp2p.dialProtocol(peerMultiaddr, protocol),
    interceptors: [],
    readMaxBytes: 10000,
    writeMaxBytes: 10000,
    useBinaryFormat: true
  })
}

export async function test (peerAddr: MultiaddrInput): Promise<string> {
  const client = createPromiseClient(ProvisionPodService, connectTo(peerAddr, '/x/trusted-pods/provision-pod/0.0.1'))

  const result = await client.provisionPod({
    pod: {
      replicas: {
        max: 3
      }
    }
  })

  if (result.error !== '3') {
    throw new Error('assertion failed')
    // return "WRONG RESULT?!"
  }
  return 'SUCCESSS'
}
