// SPDX-License-Identifier: GPL-3.0

import { createHelia, Helia } from 'helia'
import type { Multiaddr } from '@multiformats/multiaddr'
import type { PeerId } from '@libp2p/interface/peer-id'
import type { Connection } from '@libp2p/interface/connection'
import type { AbortOptions } from '@libp2p/interface'
import { Transport } from '@connectrpc/connect'
import { createLibp2pConnectTransport } from './transport-libp2p-connect'
import { AllowConnectionGater } from './connection-gater'

export { createLibp2pConnectTransport, AllowConnectionGater }

/**
 * Connect to the specified peer over the specified protocol using the specified helia node. Thin wrapper around createLibp2pConnectTransport.
 *
 * @param node the Helia node
 * @param peerAddr the address of the peer
 * @param protocol the protocol to use
 * @returns a connectrpc transport
 */
export function connectTo(
  node: Helia,
  peerAddr: PeerId | Multiaddr | Multiaddr[],
  protocol: string
): Transport {
  return createLibp2pConnectTransport({
    dialStream: async () => await node.libp2p.dialProtocol(peerAddr, protocol),
    interceptors: [],
    readMaxBytes: 10000,
    writeMaxBytes: 10000,
    useBinaryFormat: true
  })
}

/**
 * Creates a helia node and configures it to allow connecting to any peers that are explicitly dialed.
 *
 * @param opts options
 * @param opts.testMode run in test mode, which disables any peer discovery mechnisms and default libp2p background operations
 */
export async function createClient({ testMode = false }): Promise<Helia> {
  const connectionGater = new AllowConnectionGater({
    onlyAllowed: testMode
  })
  const helia = await createHelia({
    libp2p: {
      connectionGater: {
        denyDialMultiaddr:
          connectionGater.denyDialMultiaddr.bind(connectionGater)
      },
      start: !testMode,
      peerDiscovery: testMode ? [] : undefined
    }
  })
  const superDial = helia.libp2p.dial.bind(helia.libp2p)
  helia.libp2p.dial = async (
    peer: PeerId | Multiaddr | Multiaddr[],
    options?: AbortOptions
  ): Promise<Connection> => {
    connectionGater.allow(peer)
    return await superDial(peer, options)
  }
  return helia
}
