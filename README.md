# Apocryph / Trusted Pods

**NOTE**: This repository used to host the code recently moved to https://github.com/comrade-coop/apocryph-chain as part of a repository overhaul and reorganization. With due apologies to any past Github stargazers, we hope they would find this new and more active project well worth their star.

Meanwhile, the main body of code was moved here from https://github.com/comrade-coop/trusted-pods (now gone), and parts of the documentation/code might still refer to the old location. Use `git remote set-url origin git://github.com/comrade-coop/apocryph` (or `git remote set-url origin git@github.com:comrade-coop/apocryph.git`) to update your local clone/s to point to the right repository.

See [this issue](https://github.com/comrade-coop/apocryph/issues/14) for more details.

---

Trusted Pods is a decentralized compute marketplace where developers can run container pods securely and confidentially through small and medium cloud providers.

Once complete, this project would allow a regular user to deploy their own personal instance of "cloud" software (say, a wiki, website, gallery, storage backup, AI assistant, email/chat server, etc.) to another person's specialized machine, where it would run inside a secure computing enclave that no one else can access (using TEE technology) for a modest fee (however much the machine provider charges; it's a marketplace) and with regular uptime and data storage SLAs.

[![Discord](https://img.shields.io/badge/DISCORD-COMMUNITY-informational?style=for-the-badge&logo=discord)](https://discord.gg/C4e37Xhvt4)

## Spinning up a local testing environment

To start a local environment for e.g. integration-testing or evaluating the project, you can use the end-to-end tests in the `test/e2e` folder.

Typical development involves running the minikube end-to-end test, which can be done using the following command:

```bash
test/e2e/minikube/run-test.sh
```

The command will report any missing dependencies; for a full list of the required packages, you can just read the first lines of the script.

The command, once all dependencies are met, will proceed to start a local docker registry and test ethereum node, build and upload the project to them, then spin up a minikube cluster and deploy all necessary prerequisites into it, and finally deploying a pod from a [manifest file](spec/MANIFEST.md) into the cluster and then querying it over HTTP. It should display the curl command used to query the pod, and you should be able to use it yourself after the script is finished.

In addition, once you have started the minikube end-to-end test, you can also run the web UI test, which presents a sample interface that publishers can use to deploy a predefined pod template onto the minikube cluster / provider directly from their web browser.

```bash
test/e2e/webui/run-test.sh
```

Once you are done playing around with the tests, simply run the following command to delete and stop the minikube cluster:

```bash
test/e2e/minikube/run-test.sh teardown
```

(or alternatively, pass `teardown full` to also stop any local docker containers used by the test)

## Development

After editing files in `proto/` or `contracts/`, make sure to run the following commands to sync the generated files:

```bash
forge build --root contracts
go generate
npm run -ws generate
```

<!-- Note that while committing generated files is foreign to Nodejs/NPM, it's the usual way of life in the Go ecosystem, as packages are directly cloned from git rather than downloaded from the package manager. Here we are committing both in order to not require forge/protoc for JavaScript development when it's optional for Go development. -->

## Contributing

As it is, this project is still in its infancy, and most non-trivial contributions should be done only after discussing them with the team -- or else risk missing the point. So, if you fancy contributing to the project, please feel free to hop on [our Discord server](https://discord.gg/C4e37Xhvt4) or just open/reply to an issue discussing your concrete ideas for contribution.

Also, see the [`ARCHITECTURE.md`](spec/ARCHITECTURE.md) documentation for more details on the overall structure of the project.

## License

[SPDX-License-Identifier: GPL-3.0](./LICENSE.md)
