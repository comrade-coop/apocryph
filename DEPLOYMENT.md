# S3 aApp Deployment

Here are the steps required to deploy a version of the S3 aApp:

1. If the newly-released version should replicate data from the previous version, configure `BACKEND_REPLICATE_SITES` in the [`app-manifest.json`](./.github/workflows/app-manifest.json) file, next to `BACKEND_ETH_RPC` - it should match the `VITE_GLOBAL_HOST` value of the previous version. Otherwise, remove that configuration.

2. Create a new release on GitHub. This will start off the [`release-tag.yml`](./.github/workflows/app-manifest.json) GitHub Action, which deploys the configured manifest.

  The release should be named like `vX.Y.ZZ`, where X, Y, and Z are digits. The digits in the release's name will be concatenated and used to form the aApp's versioned domain name (the `__AAPPCNAME__`), which must be unique.

3. Watch the newly-started GitHub Action for the `MANUAL DNS SETUP REQUIRED FOR DOMAIN` message, then follow the instructions to add the required A and TXT DNS records.

  It should prompt for each of the three separate domains required by the aApp.

4. Open the newly-released version's app page (domain specified by `VITE_GLOBAL_HOST_APP`). At the bottom of the page you should see two Ethereum addresses listed: the aApp's own address, and the payment contract's address. Ensure that the payment contract exists on [basescan](https://basescan.org/); if it's not yet created, transfer some amount of ETH on Base to the aApp's address.

5. Congratulations! You have released a new version. If you have configured replication, it should replicate all data from the previous version after a while (completion currently not indicated), after which the old version may be taken down.
