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

...where `VITE_STORAGE_SYSTEM` can be computed by running `cast wallet address $BACKEND_ETH_PRIVATE_KEY`.

### Running

Afterwards, you can run the resulting image as follows:

```bash
docker run \
  --env BACKEND_ETH_PRIVATE_KEY=4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356 \
  --env BACKEND_ETH_WITHDRAW=0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f \
  docker.io/comrade-coop/s3-aapp:latest
```

...where  `BACKEND_ETH_PRIVATE_KEY` is the private key for the whole groups of S3 aApp nodes (..once replication rolls around), and `BACKEND_ETH_WITHDRAW` is the public address that is unique to the particular node.
