import './style.css'
import { setupCounter } from './counter.ts'
import { connectTo, createClient } from 'trusted-pods-ipfs-ts'
import { createPromiseClient } from '@connectrpc/connect'
import { ProvisionPodService } from 'trusted-pods-proto-ts'
import { multiaddr, type MultiaddrInput } from '@multiformats/multiaddr'

var helia = await createClient({testMode: true})

export async function test (peerAddr: MultiaddrInput): Promise<string> {
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
  return 'SUCCESSS'
}

window.Xtest = test

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div>
    <h1>Vite + TypeScript</h1>
    <div class="card">
      <button id="counter" type="button"></button>
    </div>
  </div>
`

setupCounter(document.querySelector<HTMLButtonElement>('#counter')!)
