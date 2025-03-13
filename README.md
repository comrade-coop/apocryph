# S3 aApp


## Docker

### Building

To build the S3 aApp image, use the following command:

```bash
docker build . -t comrade-coop/s3-aapp:latest \
  --build-arg VITE_TOKEN=0xe52a82edf1f2a0b0cd69ffb4b98a29e3637cf665 \
  --build-arg VITE_STORAGE_SYSTEM=0x14dC79964da2C08b23698B3D3cc7Ca32193d9955
  --build-arg VITE_GLOBAL_HOST=s3-aapp.kubocloud.io
  --build-arg VITE_GLOBAL_HOST_CONSOLE=console-s3-aapp.kubocloud.io
```

...where `VITE_STORAGE_SYSTEM` can will be output by the running instance of the container.
Alternatively, you can pass `--build-arg VITE_STORAGE_SYSTEM='$$$VITE_STORAGE_SYSTEM$$$'`, and later add `--env FIXUP_VITE_STORAGE_SYSTEM=true` to tell the container to replace the value at runtime.

### Running

Afterwards, you can run the resulting image as follows:

```bash
docker run \
  --volume data-vol:/data \
  --volume secrets-vol:/shared_secrets \
  --env BACKEND_ETH_WITHDRAW=0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f \
  docker.io/comrade-coop/s3-aapp:latest
```

...where `BACKEND_ETH_WITHDRAW` is the public address that is unique to the particular node.

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
