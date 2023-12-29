// SPDX-License-Identifier: GPL-3.0

import { createPublicClient, createWalletClient, http } from 'viem'
import { foundry } from 'viem/chains'
import { privateKeyToAccount } from 'viem/accounts'
import { createClient } from 'apocryph-ipfs-ts'

/**
 * Client for a public Ethereum node, used for reading, estimating gas fees, and simulating transactions
 */
export const publicClient = createPublicClient({
  chain: foundry,
  transport: http()
})

/**
 * Client for a private Ethereum node, used for signing transactions
 */
export const walletClient = createWalletClient({
  chain: foundry,
  account: privateKeyToAccount('0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a'), // TODO= anvil.accounts[2] -- remove hardcode / use metamask!
  transport: http() // custom(window.ethereum)
})

/**
 * Client for IPFS/libp2p
 */
export const heliaNodePromise = createClient({ testMode: true })
