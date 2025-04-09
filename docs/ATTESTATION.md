# Attestation

Attesting an Aapp allows you to get cryptographic proof linking back to the CPU's manufacturer that attests that the CPU is running a specific piece of code without allowing for any tampering with it. In turn, you can use this proof

## How to attest the Aapp

### Self-attestation

To review the attestation of the S3 Aapp, you can use the "View Attestation" link at the bottom of the [Apocryph Console](https://console.apocryph.io). That link will take you to `/.well-known/attest/view`, a page that will automatically verify the signature of the TEE environment, the deployment that created it, and the TLS certificate used for connecting to it, [Apocryph Aapp toolkit documentation](https://github.com/comrade-coop/aapp-toolkit/blob/main/docs/ONCHAINATTESTATION.md).

In addition, if you scroll down to the "Certificate Check" part of the attestion page, you can use the Manual Certificate Verification process described there to verify the certificate for yourself.

The same `/.well-known/attest/view` page exists for the S3 API and the Minio Console—append the path to the correct hostname to see it.

To fully attest the App, make sure to check that the values reported by the Self-attestation page match those in the linked GitHub job.
## Exploring the codebase

The S3 Aapp attestation just tells you that the application is running a particular version of the S3 Aapp codebase. This by itself is meaningless if you cannot preview said codebase.

You can access the whole code in the [S3 Aapp repository](https://github.com/comrade-coop/s3-aapp). In particular, the repository includes instructions for running a local instance of the Aapp, which should help you understand the code and its moving pieces faster.

A critical piece of the codebase is the [application manifest](https://github.com/comrade-coop/s3-aapp/blob/master/.github/workflows/app-manifest.json), which describes how the S3 Aapp is deployed, including the version of the Aapp toolkit used (which is in turn responsible for setting up and securing the disks the S3 Aapp uses to store all buckets on).
