import { createPublicClient, createWalletClient, http } from 'viem'
import { foundry } from 'viem/chains'
import { privateKeyToAccount } from 'viem/accounts'
import { createClient } from 'trusted-pods-ipfs-ts'

export const publicClient = createPublicClient({
  chain: foundry,
  transport: http()
})
export const walletClient = createWalletClient({
  chain: foundry,
  account: privateKeyToAccount('0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a'), // TODO= anvil.accounts[2],
  transport: http() // custom(window.ethereum)
})
export const heliaNodePromise = createClient({ testMode: true })
