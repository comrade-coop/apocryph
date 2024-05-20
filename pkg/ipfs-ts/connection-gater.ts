// SPDX-License-Identifier: GPL-3.0

import isPrivate from 'private-ip'
import { PeerId, isPeerId } from '@libp2p/interface'
import type { Multiaddr } from '@multiformats/multiaddr'
import type { ConnectionGater } from '@libp2p/interface'

/**
 * A libp2p connection gater implementation that never denies peers that have been explicitly allow()-ed
 * Necessary to connect to peers running on localhost without an external IP address.
 */
export class AllowConnectionGater implements ConnectionGater {
  public allowed: Set<string>
  public onlyAllowed: boolean

  /**
   * @param opts Options
   * @param opts.onlyAllowed Only allow explicitly-allowed peers, overriding the default connection-gater behavior
   */
  constructor({ onlyAllowed = false }) {
    this.onlyAllowed = onlyAllowed
    this.allowed = new Set<string>()
  }

  /**
   * Allow dialing the specified peer
   */
  allow(peer: PeerId | Multiaddr | Multiaddr[]): void {
    if (isPeerId(peer)) {
      this.allowed.add(peer.toString())
    } else {
      const addrs = Array.isArray(peer) ? peer : [peer]
      for (const ma of addrs) {
        const peerId = ma.getPeerId()
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
    const peerId = ma.getPeerId()
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
