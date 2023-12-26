// SPDX-License-Identifier: GPL-3.0

import { connectTo, createClient } from 'tapocryph-ipfs-ts'
import { createPromiseClient } from '@connectrpc/connect'
import { ProvisionPodService } from 'apocryph-proto-ts'
import { multiaddr, type MultiaddrInput } from '@multiformats/multiaddr'

const helia = await createClient({ testMode: true })

const peerAddr: MultiaddrInput = process.argv[1]

const client = createPromiseClient(ProvisionPodService, connectTo(helia, multiaddr(peerAddr), '/x/apocryph/provision-pod/0.0.1'))

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

console.log('Successfully connected from javascript!')
