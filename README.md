# S3 aApp


## Docker

### Building

To build the S3 aApp image, use the following command:

```bash
docker build . -t comrade-coop/s3-aapp:latest \
  --build-arg VITE_TOKEN=0xe52a82edf1f2a0b0cd69ffb4b98a29e3637cf665 \
  --build-arg VITE_GLOBAL_HOST=s3-aapp.kubocloud.io \
  --build-arg VITE_GLOBAL_HOST_CONSOLE=console-s3-aapp.kubocloud.io \
  --build-arg VITE_GLOBAL_HOST_APP=console-aapp.kubocloud.io
```

### Running

Afterwards, you can run the resulting image as follows:

```bash
docker run \
  --volume data-vol:/data \
  --volume secrets-vol:/shared_secrets \
  --env BACKEND_ETH_WITHDRAW=0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f \
  docker.io/comrade-coop/s3-aapp:latest
```

...where `BACKEND_ETH_WITHDRAW` is the public address to withdraw proceeds to.

### Initial configuration

When started, the aApp will attempt to deploy a smart contract using its own private key; this will typically fail due to insufficient funds.

To kickstart the aApp, you should transfer some initial gas funds to the aApp's own address (not the payment contract address!), listed at the bottom of the page and in the logs.

### Building a separate static frontend image

To build a separate container for the static aapp frontend, you can do the following:

1. Build a frontend-less image by passing `--build-arg FRONTEND_MODE=none` to the build command above.
2. Run the frontend-less container:
  ```bash
  docker run ... docker.io/comrade-coop/s3-aapp:latest s3-aapp-container
  ```
3. Grep the logs for the storage system address 
  ```bash
  docker logs s3-aapp-container | grep "VITE_STORAGE_SYSTEM"
  ```
4. Build a frontend-only image with the right storage system address, passing the same `VITE_*` build arguments as in the usual build:
  ```bash
  docker build . --target frontend-serve -t comrade-coop/s3-aapp-serve:latest \
    --build-arg VITE_TOKEN=0xe52a82edf1f2a0b0cd69ffb4b98a29e3637cf665 \
    --build-arg VITE_GLOBAL_HOST=s3-aapp.kubocloud.io \
    --build-arg VITE_GLOBAL_HOST_CONSOLE=console-s3-aapp.kubocloud.io \
    --build-arg VITE_GLOBAL_HOST_APP=console-aapp.kubocloud.io \
    --build-arg VITE_STORAGE_SYSTEM=$VITE_STORAGE_SYSTEM
  ```
5. Run the image serving the frontend:
  ```bash
  docker run comrade-coop/s3-aapp-serve:latest
  ```
  (Alternatively, copy the static frontend out of the image's `/usr/share/nginx/html/` path and serve through IPFS/an existing server/etc.)


## Tilt

To build and run with [`tilt`](https://tilt.dev/), you can just enter the following command.

```bash
tilt up
```

The web interface that opens should offer you the option of opening a test Firefox instance, which should in automatically open to the locally-deployed application.

Make sure to configure the MetaMask extension inside that test Firefox instance so you can use all features of the frontend. As it launches with a local `anvil` instance for Ethereum support, you should import the following private key for a pre-funded address, as well as configure the following chain:
```
Private key:      0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

Chain:            localhost
Default RPC URL:  http://anvil.local:8545
Chain ID:         31337
Currency symbol:  Local
Block explorer:   <none>
```

### Testing replication

To test replication between two instances of the backend, you can use the following steps:

1. Launch only the first ("s3-zero") instance:
  ```bash
  tilt up -- sc
  ```
2. Log into the Console using the frontend and add some files to a bucket, etc.
3. Switch tilt to run both the first ("s3-zero") and the second ("s3-one") instance
  ```bash
  tilt args -- mc
  ```
4. Wait a bit for everything to start and replicate.
5. Log into the Console again and observe that files are still present.
6. Switch tilt to run only the second ("s3-one") instance
  ```bash
  tilt args -- sc2
  ```
7. Refresh the Console and observe that files are still present, despite the original instance now being disabled.
