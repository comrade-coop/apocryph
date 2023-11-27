import isPrivate from 'private-ip'
import { type PeerId, isPeerId } from '@libp2p/interface/peer-id'
import type { Multiaddr } from '@multiformats/multiaddr'
import type { ConnectionGater } from '@libp2p/interface/connection-gater'

export class AllowConnectionGater implements ConnectionGater {
  public allowed: Set<string>;
  public onlyAllowed: boolean;

  constructor({onlyAllowed = false, allowed = []}) {
    this.onlyAllowed = onlyAllowed
    this.allowed = new Set<string>(allowed)
  }

  allow(peer: PeerId | Multiaddr | Multiaddr[]) {
    if (isPeerId(peer)) {
      console.log(peer.toString())
      this.allowed.add(peer.toString())
    } else {
      let addrs = Array.isArray(peer) ? peer : [peer];
      for (let ma of addrs) {
        let peerId = ma.getPeerId()
        if (peerId != null) {
          this.allowed.add(peerId)
        }
        this.allowed.add(ma.toString())
      }
    }
  }

  async denyDialMultiaddr(ma: Multiaddr): Promise<boolean> {
    if (this.allowed.has(ma.toString())) {
      return false
    }
    let peerId = ma.getPeerId()
    if (peerId != null && this.allowed.has(peerId)) {
      return false
    }
    if (this.onlyAllowed) {
      return true
    }
    const tuples = ma.stringTuples()

    if (tuples[0][0] === 4 || tuples[0][0] === 41) {
      return Boolean(isPrivate(`${tuples[0][1]}`))
    }

    return false
  }
}

