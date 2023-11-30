import { connectTo, createClient } from 'trusted-pods-ipfs-ts'
import { createPromiseClient } from '@connectrpc/connect'
import { ProvisionPodService } from 'trusted-pods-proto-ts'
import { multiaddr, type MultiaddrInput } from '@multiformats/multiaddr'

var helia = await createClient({testMode: true})

var peerAddr: MultiaddrInput = process.argv[1]

const client = createPromiseClient(ProvisionPodService, connectTo(helia, multiaddr(peerAddr), '/x/trusted-pods/provision-pod/0.0.1'))

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

console.log('Sucessfully connected from javascript!')

