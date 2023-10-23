# Registry Protocol
The Registry protocol connects clients and providers in a decentralized environment, promoting service discovery, competitive pricing, reputation system, and transparency.

## Requirements

* **Service Discovery**: Clients can browse the services available
* **Provider Accessibility**: Providers can offer either accessible endpoints for receiving pod execution requests or an endpoint service that makes these endpoints available.
* **Pricing Compliance** Providers must adhere to a specified pricing table format that encompasses all the available resources.
* **Automated Provider Selection**: The application is capable of autonomously selecting a provider based on client preferences, which may encompass Region, pricing, and optionally available funds and estimated execution time.

## Smart contract

* Holds a list of providers
* Each provider have a name, list of Regions, contact information, request endpoint(s), and attestation details.
* Providers can submit their custom pricing table
* Each pricing table is identified by a unique ID.
* An event is triggered For each submission of a new pricing table (this could be valuable for publishers who are interested in exploring competitive pricing or for providers looking to compete with one another).
* Each pricing table is associated with a set of available providers.
* It has a Reputation process that works as following:
    * When a provider consistently executes a pod, which means the funds has not been unlocked and unlock time is continiously prolonged, the payment contract informs the registry of such execution period.
    * The Registry evaluates whether the execution period has reached a predefined minimum time.
    * If a minimum execution period has been met, the provider's reputation is improved, typically manifested through points, stars, or equivalent rating metric.
* A provider can unsubscribe from a pricing table and move to another one, if the pricing table has no providers it is deleted

Note: Given that execution time is guranteed by the TEE, The minimum predefined time serves the purpose of preventing providers from setting unrealistically short execution periods by creating fake clients and fake channels to artificially boost their reputation, and adds a layer of complexity by incurring execution costs and time before reputation improvement can occur.

## Provider

* Get the list of available pricing tables
* Create or Retreive Attestation details
* Configure a Profile:
    * Contact Information,
    * Name
    * Pod execution request endpoints
    * Attestation Details
    * Available Region(s) for pod execution
    * Support for Edge computing(optional)
* Chooses a pricing table or create a custom pricing table.

## Publisher

* Retreives the list of available pricing tables from the contract
* Create a configuration which includes:
    * pricing table
    * Provider selection priority configuration:
        * by points
        * by Region(s)
    * Optionally, you can specify the amount of funds you are prepared to allocate to your pod or the desired duration for its execution. The client application will then automatically propose a pricing table(s) tailored to your preference.
* If automatic selection is enabled, it will select a provider filtered by the configuration
* The publisher creates a payment channel configured with the selected pricing planand and intiates the pod execution request protocol 

## Pricing table (wip)

the pricing table should include the necessary information for billing clients, and it could be as simple as the following:

| Resource          | Description| Price                            |
|-------------------|------------|-----------------------------|
| **CPU**           | N° of Cores/ N° of vCPUs|$0.0001 VCPU(min/s/ms)|
| **RAM Capacity**  | Capacity (ex: GB)| $0.00001 GB/(min/s/ms) |
| **Storage**       | Type (e.g., Block, Object)|$0.00001 GB/(min/s/ms)|             |
| **GPU (Optional)**| Model|$0.0001 per Execution (min/s/ms)|

or it could be split into categories with more detailed information:

### Compute Pricing

* **CPU**
    
    | Resource| Description | Number of Cores| vCPUs |     Model | TEE Type| Price per Unit |
    |-|-|-|-|-|-|-|
    | **CPU**   | Processing power   | Cores           | vCPUs    | Intel, AMD, ARM, ...etc      | Enclaves, CVMs, ..etc | $0.0001 VCPU(min/s/ms)    |

- **Ram**

    |Ressource|Description|Capacity| Price|
    |-|-|-|-|
    | **RAM**   | Memory capacity | Capacity (ex: 1GB) | $0.00001         GB/(min/s/ms) 


## Storage Pricing

| Resource      | Description        | Capacity |     Storage Type | Price per Unit |
|---------------|--------------------|----------|-------------- |-----------------|
| **Storage** | Storage resources | Capacity (ex: 10GB) | Block, Object, ...etc | $000001 per GB(min/s/ms) |

