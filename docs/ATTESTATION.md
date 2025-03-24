# Attestation

Attesting an Aapp allows you to get cryptographic proof linking back to the CPU's manufacturer that attests that the CPU is running a specific piece of code without allowing for any tampering with it. In turn, you can use this proof

## How to attest the Aapp

To review the attestation of the S3 Aapp, you can use the "View Attestation" link at the bottom of the [Apocryph Console](https://console.apocryph.io). That page will guide you through reviewing the TLS certificate in use, which should be signed as described in the [Apocryph Aapp toolkit documentation](https://github.com/comrade-coop/aapp-toolkit/blob/main/docs/ONCHAINATTESTATION.md).

(TODO: Expand documentation with the information from the "View Attestation" page) <!-- e.g. https://github.com/comrade-coop/aapp-toolkit/blob/aa1625f2ad46df0a5692ce31cd268ad21dd76f7c/azure-attestation/web/index.html#L269 -->

## Exploring the codebase

The S3 Aapp attestation just tells you that the application is running a particular version of the S3 Aapp codebase. This by itself is meaningless if you cannot preview said codebase.

You can access the whole code in the [S3 Aapp repository](https://github.com/comrade-coop/s3-aapp). In particular, the repository includes instructions for running a local instance of the Aapp, which should help you understand the code and its moving pieces faster.

A critical piece of the codebase is the [application manifest](https://github.com/comrade-coop/s3-aapp/blob/master/.github/workflows/app-manifest.json), which describes how the S3 Aapp is deployed, including the version of the Aapp toolkit used (which is in turn responsible for setting up and securing the disks the S3 Aapp uses to store all buckets on).
