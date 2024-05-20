// SPDX-License-Identifier: GPL-3.0

import './style.css'
import { paymentAbi, mockTokenAbi } from 'apocryph-abi-ts'
import { bytesToHex } from 'viem'
import { publicClient, walletClient } from './connections'

// Copied and adapted from pkg/publisher/fund.go

/**
 * Fund a given payment channel. This function currently only does the initial funding of a channel, and will not deposit funds if a channel already exists.
 *
 * @param config deployment configuration; includes payment channel details as well as provider details
 * @param funds the amount of funds to approve+send to the channel
 * @param opts additional options
 * @param opts.unlockTime time (in seconds) the publisher needs to wait before unlocked tokens may be withdrawn. Provider may reject channels with unlock time that is too low.
 * @param opts.mintFunds use MockToken to mint funds before the transfer.
 */
export async function fundPaymentChannel(
  config: {
    payment: {
      paymentContractAddress: Uint8Array
      podID: Uint8Array
    }
    provider: {
      ethereumAddress: Uint8Array
    }
  },
  funds: bigint,
  { unlockTime = 5n * 60n, mintFunds = false }
): Promise<void> {
  const paymentContractAddress = bytesToHex(
    config.payment.paymentContractAddress,
    { size: 20 }
  )
  const providerAddress = bytesToHex(config.provider.ethereumAddress, {
    size: 20
  })
  const podId = bytesToHex(config.payment.podID, { size: 32 })
  if (funds > 0n) {
    const tokenAddress = await publicClient.readContract({
      abi: paymentAbi,
      address: paymentContractAddress,
      functionName: 'token'
    })

    if (mintFunds) {
      const { request } = await publicClient.simulateContract({
        abi: mockTokenAbi,
        address: tokenAddress,
        account: walletClient.account.address,
        functionName: 'mint',
        args: [funds]
      })
      await walletClient.writeContract(request)
    }

    const approved = await publicClient.readContract({
      abi: mockTokenAbi,
      address: tokenAddress,
      functionName: 'allowance',
      args: [walletClient.account.address, paymentContractAddress]
    })

    if (approved < funds) {
      const { request } = await publicClient.simulateContract({
        abi: mockTokenAbi,
        address: tokenAddress,
        functionName: 'approve',
        account: walletClient.account.address,
        args: [paymentContractAddress, funds]
      })
      await walletClient.writeContract(request)
    }

    const { request } = await publicClient.simulateContract({
      abi: paymentAbi,
      address: paymentContractAddress,
      functionName: 'createChannel',
      account: walletClient.account.address,
      args: [providerAddress, podId, unlockTime, funds]
    })
    await walletClient.writeContract(request)
  }
}
