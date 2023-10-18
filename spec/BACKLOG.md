# Backlog

This file seeks to document the tasks left to do in this repository, as well as design flaws and accumulated tech debt present in the current implementation. 

Rationale for not shoving all of this in GitHub Issues: while Issues are a great way for users and contributors to voice concerns and problems, for huge planned milestones, it feels simpler to just have them listed in a format more prone to conveying information and not to discussing the minutiae details.

## Features yet to be implemented

* Write and integrate a Registry contract
* Run the Provider in Constellation instead of Minikube, to allow for attestation
* Implement a way for the Publisher to monitor, manage, and edit the deployed pod

## Technical debt accumulated

### Prometheus metrics used for billing

Status: Alternative prototyped

Prometheus's [documentation](https://prometheus.io/docs/introduction/overview/#when-does-it-not-fit) explicitly states that Prometheus metrics are not suitable for billing since Prometheus is designed for availability (an AP system by the [CAP theorem](https://en.wikipedia.org/wiki/CAP_theorem)) and not for consistency / reliability of results (which is what a CP system would be).

Despite that, the current [implementation](../pkg/prometheus/) uses Prometheus and `kube-state-metrics` to fetch the list of metrics that billing is based on. A prototype was created in the [`metrics-monitor`](https://github.com/comrade-coop/trusted-pods/tree/metrics-monitor) branch to showcases an alternative way to fetch the same data from Kubernetes directly and avoid any possible inconsistencies in the result, yet it was decided that it's better to iterate quickly with Prometheus first instead and come back to this idea later.

### A single monolithic `tpodserver` service

Status: Correct as needed

Currently, the whole of the Trusted Pods Provider client/node is implemented as a single long-running process deployed within Kubernetes. Despite this being easier to implement, it would be beneficial to make parts of that service more reusable by splitting off libp2p connections, pod deployment, metrics collection, and smart contract invoicing into their own parts can be changed or deployed on their own.

## Missing features

### Storage reliability

Status: Barebones design done

For Trusted Pods to be a viable platform for deploying mission-critical software to, it needs to provide some guarantees about the longevity of data stored on it. This is still up to be designed and incorporated into the overall design and architecture.

The best solution for this we have found is establishing contracts and a protocol for Providers to stake that they will keep specific stored data available. In theory, to identify the data, a version identifier similar to ZFS's root checksum/version should be sufficient, with the in-TEE encryption of the volume providing the confidentiality of the data - at which point a container running with the same volume attached can run a full filesystem check and confirm the root checksum/version. If the Provider fails to execute the filesystem check within some (generous) time limit in response to a challenge by the Publisher, the storage can not be shown to still exist, and therefore the stake is lost and can be transferred to the Publisher.

### Uptime reliability

Status: Needs more design work

A closely-related concept to storage reliability is that of uptime reliability. While scale-down-to-zero pods might never have great first response latency, even for them Publishers should be able to ask for reliable uptime, a guarantee that the system will remain up an running.

Here again, the best system seems to be staking -- and is similar to what the cloud industry is already doing with [SLAs](https://en.wikipedia.org/wiki/Service-level_agreement). It is, however, a bit trickier to prove that the application was down due to a fault of the Provider, and there are different kinds of downtime, such as network outages and power outages that we might or might not want to handle differently.

### Software licensing

Status: Needs more design work

A useful feature would be to allow publishers to upload application code that other parties can then deploy to a provider -- this would enable usecases such as providing a pay-for-usage subscription to a service where the publisher does not have to take on the risk and responsibility of running a master instance of that service.

This should already be achievable in theory by having a master service which confirms the attestation and licensing status of deployed pods and only then provides them with a decryption key for the rest of the application. However, if it is integrated with the rest of the Trusted Pods platform, we should be able to achieve faster scale-up-from-zero performance, as well as lower some of the counter-party risk of the master service going down.

### Individual TEEs

Status: Needs more usecases

Currently, the architecture of Constellation uses a single TEE encompassing all Kubernetes pods running in the cluster. However, for extra isolation of individual tenants, it could be beneficial to have separate TEEs for each publisher / pod / application. To implement that, we will likely end up scrapping Constellation and revamping the whole attestation process. As this is quite a bit of design and implementation work while the gains at this stage are minimal, we have opted to let the idea ruminate for the moment.

### Forced scale-down

Status: Conceptualized

It would be great if we didn't just rely on KEDA's built-in scaling down after a certain time, but also allowed Pods to request their own scaling down.
