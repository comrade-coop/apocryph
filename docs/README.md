# Introduction

> **Warning**: The Apocryph S3 aApp is beta software
> 
> While this aApp might work for a particular usecase, no features of it should be considered stable and final until a stable release rolls around.  
> In particular, there may be security vulnerabilities within the aApp's code (that you can [freely review](https://github.com/comrade-coop/s3-aapp)) that will yet be fixed in future releases.

The Apocryph S3 aApp (autonomous application) allows you to store files in an S3-compatible filesystem backed by an attested trusted execution enclave.  
This allows you to store files, application data, backups, and more, in a private and secure manner, where no one else can read them.

## Features

* **S3-compatible API**: Courtesy of [Minio](https://min.io/), we support most S3 operations that modern applications require. [Full list of available S3 features.](https://min.io/docs/minio/linux/reference/s3-api-compatibility.html)
* **Cryptocurrency-based payments**: The aApp accepts payments in [USDC](https://www.circle.com/multi-chain-usdc/base) on the [BASE network](https://www.base.org), billed in 5 minute increments.
* **Fully attested execution**: The S3 aApp executes within an attested environment powered by the [Apocryph aApp toolkit](https://github.com/comrade-coop/aapp-toolkit). You can [verify the attested environment](./ATTESTATION.md) to confirm that no one can impersonate the S3 aApp.
* **Open source code**: The [S3 aApp's repository](https://github.com/comrade-coop/s3-aapp) is fully open-source, allowing anyone to review and trust the code that's running within the attested environment.
* **Web interface for managing payments**: Accessible at the [Apocryph Console](https://console.apocryph.io). See the [usage documentation](./USAGE.md) to learn how to use it.
* **Web interface for managing stored files**: Courtesy of [Minio](https://min.io/); can be accessed through the [Apocryph Console](https://console.apocryph.io).

## Next steps

* Read more about how to [attest the S3 aApp](./ATTESTATION.md)
* Learn how to [use the S3 aApp](./USAGE.md)
