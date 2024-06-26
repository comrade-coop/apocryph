# Application Registry
The Registry connects publishers and providers in a decentralized environment, promoting service discovery, competitive and transparent pricing.

## Smart contract

* Holds a list of providers
* Each provider have a name, Region(s), contact information, request endpoint(s), and attestation details.
* Providers can submit their custom [pricing table](./PRICING.md)
* Each pricing table is identified by a unique ID.
* An event is triggered For each submission of a new pricing table (this could be valuable for publishers who are interested in exploring competitive pricing or for providers looking to compete with one another).
* Each pricing table is associated with a set of available providers.
* A provider can unsubscribe from a [pricing table](./PRICING.md) and move to another one, if the pricing table has no providers it is deleted


## Provider

* Get the list of available [pricing tables](./PRICING.md) from the registry
* Create or Retrieve Attestation details
* Registering in the contract by providing the following information:
    * Contact Information,
    * Name
    * Pod execution request endpoints
    * Attestation Details
    * Available Region(s) for pod execution
    * Support for Edge computing(optional)
* Chooses a [pricing table](./PRICING.md) or create a custom pricing table.

## Publisher

* Retrieves the list of available pricing tables from the contract
* Create a configuration which includes:
    * pricing table
    * Choose the Region(s) in which the pod will be hosted.
    * Optionally, you can specify the amount of funds you are prepared to allocate to your pod or the desired duration for its execution. The client application will then automatically propose a pricing table(s) tailored to your preference.
* If automatic selection is enabled, it will select a provider filtered by the configuration
* the publisher pings the provider and checks its availability
  * In the event that the provider is offline, the publisher iterates through the provider list associated with the pricing table until it identifies an available provider.
* The publisher creates a payment channel configured with the selected pricing planand and initiates the pod execution request protocol 
