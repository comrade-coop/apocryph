import { erc20Abi, mockTokenAbi } from 's3-aapp-abi'
import { PublicClient, WalletClient, Address, getContract } from 'viem'
import { config } from './wallet'

export const tokenAddress: Address = (import.meta.env.VITE_TOKEN || '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9').trim()
export const paymentAddress: Address = (import.meta.env.VITE_STORAGE_SYSTEM || '0xef11D1c2aA48826D4c41e54ab82D1Ff5Ad8A64Ca').trim()
export const aappAddress: Address = (import.meta.env.VITE_AAPP_ADDRESS || '0x14dC79964da2C08b23698B3D3cc7Ca32193d9955').trim()

type Unsubscribe = () => void

export function watchAvailableFunds(publicClient: PublicClient, accountAddress: Address, callback: (allowance?: bigint, balance?: bigint) => void): Unsubscribe {
  async function refresh() {
    callback(undefined, undefined) // We are refreshing, blur out
    const allowance = await publicClient.readContract({
      abi: erc20Abi,
      address: tokenAddress,
      functionName: 'allowance',
      args: [accountAddress, paymentAddress],
    })
    const balance = await publicClient.readContract({
      abi: erc20Abi,
      address: tokenAddress,
      functionName: 'balanceOf',
      args: [accountAddress],
    })
    callback(allowance, balance)
  }

  const unsubscribe = publicClient.watchContractEvent({
    abi: erc20Abi,
    address: tokenAddress,
    args: { from: accountAddress, owner: accountAddress, to: paymentAddress },
    onLogs() {
      refresh()
    }
  })

  refresh()

  return unsubscribe
}

export async function depositFunds(publicClient: PublicClient, walletClient: WalletClient, depositAmount: bigint) {
  const client = {public: publicClient, wallet: walletClient}
  const wallet = walletClient.account!
  const writeOptions = {account: wallet, chain: config.chains[0]}

  const token = getContract({
    address: tokenAddress,
    abi: erc20Abi,
    client,
  })

  if (depositAmount > 0n) {
    const allowance = (await token.read.allowance([wallet.address, paymentAddress]))
    if (allowance != depositAmount) {
      await token.write.approve([paymentAddress, depositAmount], writeOptions)
    }
  } else {
    const allowance = (await token.read.allowance([wallet.address, paymentAddress]))
    if (allowance != depositAmount) {
      await token.write.approve([paymentAddress, depositAmount], writeOptions)
    }
  }
}

export async function debugMintFunds(publicClient: PublicClient, walletClient: WalletClient, depositAmount: bigint) {
  const client = {public: publicClient, wallet: walletClient}
  const wallet = walletClient.account!
  const writeOptions = {account: wallet, chain: config.chains[0]}

  const token = getContract({
    address: tokenAddress,
    abi: mockTokenAbi,
    client,
  })
  const balance = (await token.read.balanceOf([wallet.address]))
  if (depositAmount > balance) {
    await token.write.mint([depositAmount - balance], writeOptions)
  }
}
