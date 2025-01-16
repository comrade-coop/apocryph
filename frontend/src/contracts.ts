import { erc20Abi, mockTokenAbi, paymentV2Abi } from 's3-aapp-abi'
import { PublicClient, WalletClient, parseUnits, Address, stringToHex, keccak256, encodeAbiParameters, parseAbiParameters, getContract, zeroAddress } from 'viem'

// TODO: result from forge script script/Deploy.sol
export const paymentV2Address: Address = '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9'
// TODO: HACK Insecure!! address pulled from anvil[7], until we get a real deployment rolling
export const storageSystemAddress: Address = '0x14dC79964da2C08b23698B3D3cc7Ca32193d9955'


// TODO: manually copied from the contract
export const PERMISSION_ADMIN = 1;
export const PERMISSION_MANAGE_AUTHORIZATIONS = 2;
export const PERMISSION_WITHDRAW = 4;
export const PERMISSION_NO_LIMIT = 8;
export const STORAGE_CHANNEL_PERMISSONS = PERMISSION_MANAGE_AUTHORIZATIONS | PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT;
export const STORAGE_CHANNEL_RESERVATION = parseUnits('10', 18); // TODO: Compute to roughly equal expected usage over one unlockTime period
export const STORAGE_CHANNEL_DISCRIMINATOR = stringToHex('storage.apocryph.io', { size: 32 })

type Unsubscribe = () => void

export function getStorageChannelId(address: Address) {
  return keccak256(encodeAbiParameters(
    parseAbiParameters('address, bytes32'), [
      address,
      STORAGE_CHANNEL_DISCRIMINATOR,
    ]))
}

export function watchAvailableFunds(publicClient: PublicClient, accountAddress: Address, callback: (available?: bigint, reserved?: bigint) => void): Unsubscribe {
  const channelId = getStorageChannelId(accountAddress)
  async function refresh() {
    callback(undefined, undefined) // We are refreshing, blur out
    const channelInfo = await publicClient.readContract({
      abi: paymentV2Abi,
      address: paymentV2Address,
      functionName: 'channels',
      args: [channelId],
    })
    const [available, reserved] = channelInfo
    callback(available, reserved)
  }

  const unsubscribe = publicClient.watchContractEvent({
    abi: paymentV2Abi,
    address: paymentV2Address,
    args: { channelId: channelId },
    onLogs() {
      refresh()
    }
  })

  refresh()

  return unsubscribe
}

class TransientError extends Error {

}

export async function depositFunds(publicClient: PublicClient, walletClient: WalletClient, depositAmount: bigint) {
  const client = {public: publicClient, wallet: walletClient}
  const wallet = walletClient.account!
  const writeOptions = {account: wallet, chain: undefined}

  const paymentV2 = getContract({
    address: paymentV2Address,
    abi: paymentV2Abi,
    client,
  })
  const channelId = getStorageChannelId(wallet.address)
  //const channelId = await paymentV2.read.getChannelId([wallet.address, STORAGE_CHANNEL_DISCRIMINATOR])
  const tokenAddress = await paymentV2.read.token()
  const token = getContract({
    address: tokenAddress,
    abi: erc20Abi,
    client,
  })

  if (depositAmount == 0n) {
    return
  }
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

    const allowance = (await token.read.allowance([wallet.address, paymentV2.address]))
    const shouldApproveFirst = allowance < depositAmount

    if (shouldApproveFirst) {
      await token.write.approve([paymentV2Address, depositAmount], writeOptions)
    }

    const ownAuthorization = (await paymentV2.read.channelAuthorizations([channelId, wallet.address]))
    const ownPermissions = ownAuthorization[3]

    const shouldCreateChannel = (ownPermissions & PERMISSION_ADMIN) != PERMISSION_ADMIN

    if (shouldCreateChannel) {
      await paymentV2.write.createAndAuthorize([
        STORAGE_CHANNEL_DISCRIMINATOR,
        depositAmount,
        storageSystemAddress,
        STORAGE_CHANNEL_PERMISSONS,
        STORAGE_CHANNEL_RESERVATION,
        0n,
      ], writeOptions)
    } else {
      const storageAuthorization = (await paymentV2.read.channelAuthorizations([channelId, storageSystemAddress]))
      const storagePermissions = storageAuthorization[3]
      const shouldAuthorizeChannel = (storagePermissions & STORAGE_CHANNEL_PERMISSONS) != STORAGE_CHANNEL_PERMISSONS

      if (shouldAuthorizeChannel) {
        await paymentV2.write.authorize([
          channelId,
          storageSystemAddress,
          STORAGE_CHANNEL_PERMISSONS,
          STORAGE_CHANNEL_RESERVATION,
          0n
        ], writeOptions)
      }

      await paymentV2.write.deposit([channelId, depositAmount], writeOptions)
    }
  } else {
    const withdrawAmount = -depositAmount
    const channelInfo = await paymentV2.read.channels([channelId])
    const [available, reserved] = channelInfo

    const shouldUnreserve = (withdrawAmount > available - reserved)
    if (shouldUnreserve) {
      const storageAuthorization = await paymentV2.read.channelAuthorizations([channelId, storageSystemAddress])
      const storageReserved = storageAuthorization[0]
      let storageUnlockAt = storageAuthorization[2]
      const currentTime = (await publicClient.getBlock()).timestamp
      const shouldUnlock = storageUnlockAt == 0n
      if (shouldUnlock) {
        if (withdrawAmount > available - reserved + storageReserved) {
          throw new Error(`Cannot withdraw requested funds: there aren't enough funds reserved for the storage subsystem's authorization (Did you authorize additional systems for the same storage?)`)
        }
        await paymentV2.write.unlock([channelId, storageSystemAddress], writeOptions)
        const updatedStorageAuthorization = await paymentV2.read.channelAuthorizations([channelId, storageSystemAddress])
        storageUnlockAt = updatedStorageAuthorization[2]
      }
      if (currentTime < storageUnlockAt) {
        const unlockDate = new Date(Number(storageUnlockAt) * 1000)
        throw new TransientError(`Storage system authorization is currenly being unlocked; please wait until ${unlockDate} before continuing with the withdrawal.`)
      }
      // NOTE: Could also leave a leftoverReservation = withdrawAmount - (available - reserved + storageReserved)
      // But for now, remove the whole authorization altogether:
      await paymentV2.write.authorize([channelId, storageSystemAddress, 0, 0n, 0n], writeOptions)
    }

    await paymentV2.write.withdraw([channelId, zeroAddress, withdrawAmount], writeOptions)
  }
}
