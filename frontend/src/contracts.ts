import { erc20Abi, mockTokenAbi } from 's3-aapp-abi'
import { PublicClient, WalletClient, Address, getContract } from 'viem'
import { config } from './wallet'

export const tokenAddress: Address = (import.meta.env.VITE_TOKEN || '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9').trim()
export const storageSystemAddress: Address = (import.meta.env.VITE_STORAGE_SYSTEM || '0x14dC79964da2C08b23698B3D3cc7Ca32193d9955').trim()

type Unsubscribe = () => void

export function watchAvailableFunds(publicClient: PublicClient, accountAddress: Address, callback: (available?: bigint, reserved?: bigint) => void): Unsubscribe {
  async function refresh() {
    callback(undefined, undefined) // We are refreshing, blur out
    const available = await publicClient.readContract({
      abi: erc20Abi,
      address: tokenAddress,
      functionName: 'allowance',
      args: [accountAddress, storageSystemAddress],
    })
    callback(available, undefined)
  }

  const unsubscribe = publicClient.watchContractEvent({
    abi: erc20Abi,
    address: tokenAddress,
    args: { from: accountAddress, owner: accountAddress, to: storageSystemAddress },
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
    const balance = (await token.read.balanceOf([wallet.address]))
    const debugMintTokens = balance < depositAmount
    if (debugMintTokens) {
      const mockToken = getContract({
        address: tokenAddress,
        abi: mockTokenAbi,
        client,
      })
      await mockToken.write.mint([depositAmount], writeOptions)
    }

    const allowance = (await token.read.allowance([wallet.address, storageSystemAddress]))
    if (allowance != depositAmount) {
      await token.write.approve([storageSystemAddress, depositAmount], writeOptions)
    }
  } else {
    const allowance = (await token.read.allowance([wallet.address, storageSystemAddress]))
    if (allowance != depositAmount) {
      await token.write.approve([storageSystemAddress, depositAmount], writeOptions)
    }
  }
}
