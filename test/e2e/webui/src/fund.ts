import './style.css'
import { paymentABI, ierc20ABI, mockTokenABI } from 'trusted-pods-abi-ts'
import { bytesToHex } from 'viem'
import { publicClient, walletClient } from './connections'

// Copied and adapted from pkg/publisher/fund.go
export async function fundPaymentChannel (
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
  {
    unlockTime = 5n * 60n,
    mintFunds = false
  }): Promise<void> {
  const paymentContractAddress = bytesToHex(config.payment.paymentContractAddress, { size: 20 })
  const providerAddress = bytesToHex(config.provider.ethereumAddress, { size: 20 })
  const podId = bytesToHex(config.payment.podID, { size: 32 })
  if (funds > 0n) {
    const tokenAddress = await publicClient.readContract({
      abi: paymentABI,
      address: paymentContractAddress,
      functionName: 'token'
    })

    if (mintFunds) {
      const { request } = await publicClient.simulateContract({
        abi: mockTokenABI,
        address: tokenAddress,
        account: walletClient.account.address,
        functionName: 'mint',
        args: [
          funds
        ]
      })
      await walletClient.writeContract(request)
    }

    const approved = await publicClient.readContract({
      abi: ierc20ABI,
      address: tokenAddress,
      functionName: 'allowance',
      args: [
        walletClient.account.address,
        paymentContractAddress
      ]
    })

    if (approved < funds) {
      const { request } = await publicClient.simulateContract({
        abi: ierc20ABI,
        address: tokenAddress,
        functionName: 'approve',
        account: walletClient.account.address,
        args: [
          paymentContractAddress,
          funds
        ]
      })
      await walletClient.writeContract(request)
    }

    const { request } = await publicClient.simulateContract({
      abi: paymentABI,
      address: paymentContractAddress,
      functionName: 'createChannel',
      account: walletClient.account.address,
      args: [
        providerAddress,
        podId,
        unlockTime,
        funds
      ]
    })
    await walletClient.writeContract(request)
  }
}
