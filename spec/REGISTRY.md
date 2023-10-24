# Registry Protocol
The Registry protocol connects publishers and providers in a decentralized environment, promoting service discovery, competitive and transparent pricing.


## Smart contract

* Holds a list of providers
* Each provider have a name, Region(s), contact information, request endpoint(s), and attestation details.
* Providers can submit their custom [pricing table](https://github.com/comrade-coop/trusted-pods/blob/master/spec/PRICING.md)
* Each pricing table is identified by a unique ID.
* An event is triggered For each submission of a new pricing table (this could be valuable for publishers who are interested in exploring competitive pricing or for providers looking to compete with one another).
* Each pricing table is associated with a set of available providers.
* A provider can unsubscribe from a [pricing table](https://github.com/comrade-coop/trusted-pods/blob/master/spec/PRICING.md) and move to another one, if the pricing table has no providers it is deleted


## Provider

* Get the list of available [pricing tables](https://github.com/comrade-coop/trusted-pods/blob/master/spec/PRICING.md) from the registry
* Create or Retreive Attestation details
* Registering in the contract by providing the following informations:
    * Contact Information,
    * Name
    * Pod execution request endpoints
    * Attestation Details
    * Available Region(s) for pod execution
    * Support for Edge computing(optional)
* Chooses a [pricing table](https://github.com/comrade-coop/trusted-pods/blob/master/spec/PRICING.md) or create a custom pricing table.

## Publisher

* Retreives the list of available pricing tables from the contract
* Create a configuration which includes:
    * pricing table
    * Choose the Region(s) in which the pod will be hosted.
    * Optionally, you can specify the amount of funds you are prepared to allocate to your pod or the desired duration for its execution. The client application will then automatically propose a pricing table(s) tailored to your preference.
* If automatic selection is enabled, it will select a provider filtered by the configuration
* The publisher creates a payment channel configured with the selected pricing planand and intiates the pod execution request protocol 

