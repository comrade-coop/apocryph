# Road to MVP

## Base Protocol
A Publisher connects to a Provider to run a pod. No understanding of a "network". 

- **Pricing**: Publisher requests pricing details and attestation from Provider. ✅ 
- **Payment**: Publisher transfers payment to Provider, typically through on-chain escrow. ✅ 
- **Provisioning**: Publisher communicates image details to the Provider, including exact resource requirements, external port/routing information, and any volumes/secrets. Also allows modifying an already running pod, likely with a restart. ✅ (also includes encryption of images/secrets with an extra layer of keys for easier upload-once-deploy-anywhere)
- **Monitoring**: Publisher requests monitoring data from the Provider. ✅ 
- **Stacking**: Allows publishers to get repaid in case the provider goes dark. ⏳
- **Paired-Down Constellation**: For providers with monitoring and admin interface. Addons? ⏳ (partially done)
- **Tooling for Publishers**: To create and deploy pod configuration. ✅ 

## Autoscaler Protocol
Once deployed, allows an app to automatically redeploy itself on various providers.

- **Provider Selection Configuration**: Allows configuring criteria (needs to specify) for selecting providers to deploy to and for deciding how many instances to run. ⏳
- **Standalone Mode**: Works as a sidecar container; accepts pod configuration(s) for base protocol and redeploys that using the provider selection config and a self-contained wallet. ⏳
- **Network Mode**: Works as a separate network/autonomous application; accepts same pod config and provider selection config, but likely through an extra protocol specifically for that. ⏳
  - **Payment**: Allows paying the autoscaler; the autoscaler then keeps track of per-application balances. ⏳
- **Routing**: Providers service discovery and hosts DNS for the applications deployed through it. ⏳
  - **DNS of Autoscaler Network**: Linked to a real domain where it's discoverable as a real authoritative DNS nameserver. ⏳
- **Potential Load Balancing**: If one would rather have something redistributing the load in front of their application containers (required for "true" scale-to-zero; but probably not what one really wants - one would rather have things deployed on a provider and scaled to zero on that provider). ⏳
- **Key Management**: Manages keys for application secret decryption (so it can actually deploy something). ⏳

## Marketplace Autonomous Application
- **Primary Function**: Exists as an autonomous application, where providers can get themselves listed, and publishers (/autoscalers) can query using various criteria (needs to specify). ⏳
- **Reverse Market**: Applications get listed and providers can bid on them. ⏳
- **Integration**: Integrated with publisher (tooling) and provider (addon? so it's outside the trusted base) out-of-the-box. ⏳

## Storage Application
- **Autonomous Applications**: One or more autonomous applications of Minio/esque storage, rebalanced on various providers to provide resilience. ⏳
