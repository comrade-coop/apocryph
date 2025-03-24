
import { WalletClient } from 'viem'
import { createSiweMessage } from 'viem/siwe'

const authenticationDomain = (import.meta.env.VITE_GLOBAL_HOST_APP || "console-aapp.kubocloud.io").trim()
const consoleUrl = 'https://' + (import.meta.env.VITE_GLOBAL_HOST_CONSOLE || "console-s3-aapp.kubocloud.io").trim()

export async function getSiweToken(walletClient: WalletClient, tokenExpirationSeconds?: number) {
  const wallet = walletClient.account!

  const issuedAt = new Date()
  const expiry = tokenExpirationSeconds ? issuedAt.valueOf() + tokenExpirationSeconds * 1000 : undefined

  const message = createSiweMessage({
    address: wallet.address,
    chainId: walletClient.chain!.id,
    domain: authenticationDomain,
    nonce: issuedAt.valueOf().toString(),
    uri: consoleUrl,
    version: '1',
    expirationTime: expiry ? new Date(expiry) : undefined,
  })
  const signature = await walletClient.signMessage({
    message,
    account: wallet
  })
  const token = JSON.stringify({message: message, signature: signature})
  return token
}
