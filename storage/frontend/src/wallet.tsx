import { createConfig, http } from 'wagmi'
import { foundry } from 'wagmi/chains'
import { metaMask } from 'wagmi/connectors'

export const config = createConfig({
  chains: [foundry],
  connectors: [metaMask()],
  transports: {
    [foundry.id]: http(),
  },
})
