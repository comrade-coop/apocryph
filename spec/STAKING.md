# Staking over reputation mechanisms

(Document status: essay)

While discussing the reliability aspects (further explored in the [Storage](STORAGE.md) and [Uptime](UPTIME.md) documents), invariably someone (including, at times, the author of this document) would suggest that some kind of reputation mechanism could be used to showcase the reliable providers and indirectly punish the unreliable ones after all their customers flock away to more reliable providers. However, each of those discussions ended up in deciding that reputation systems, as a whole, are not worth for Trusted Pods. As it seems like that topic is some kind of attractor, to borrow a term from chaos theory, here is a list of some reasons that make reputation systems ultimately the wrong choice in most situations.

## Reputation systems solve the wrong problem

Whenever reliability is involved, the problem at hand is counterparty risk -- the risk that the other party of a transaction would not see their end of the deal through. In our case, we are usually seeking to protect the publisher from a malicious or abusive provider, so that a provider cannot just go offline or in some other manner not deliver the computational services they are offering.

However, reputation systems ultimately solve a different problem: a reputation system gives us a way to quickly get acquainted with peers that have a proven past track record. And as the investing maxim goes, past performance is no guarantee of future results; that is, a peer who has been reliable in the past is no guarantee that they will be reliable in the future.

In addition, a reputation system tends to outsource/externalize trust. That is, instead of trusting a particular provider (to not blunder and loose your data), you are now trusting a whole provider-rating system (to not blunder and give you a malicious provider). Or, to frame this from another angle, a reputation system makes statements of the sort of "*everybody* trusts \<X\> to not lose data" and not "*you can* trust \<X\> to not lose *your* data"; that is, it assumes that all publishers trust providers *the same* -- which is not a tenable assumption.

So overall, a reputation system, despite its apparent elegance, ends up solving the problem of reliability only indirectly, and in doing so ends up making extra assumptions on part of the providers -- and as such is to be avoided.

## Reputation systems are inherently gameable

Furthermore, as much as we enjoy reputation systems being discussed, in practice, all reputation systems we've observed end up having some kind of flaw that renders them gameable or otherwise not a good estimate for the target metric. (e.g. StackOverflow points do not strongly correlate with programming skill.) This is probably be best summed up with [Goodhart's law](https://en.wikipedia.org/wiki/Goodhart%27s_law) that "When a measure becomes a target, it ceases to be a good measure".

While one can, in theory, design a reputation system so advanced that it can avoid running into any kind of attack (be it a Sybil attack, a providers-publish-to-themselves attack, or a providers-pretend-to-have-more-CPUs-and-throttle-execution attack, for a few that have surfaced around such conversations) and do so without it becoming a worthless proxy for reliability because everybody is just optimizing for that one measure, the fact is that this feels like the task of designing a ""sufficiently smart compiler"", and as such is outside the immediate scope of the project.

## Reputation systems do not help when faults do happen

Reputation systems somewhat eventually-consistent in that they unable to help if something bad happens, but once the dust settles, they would adapt to the new state of affairs and start giving out better results than before.. until the next fault occurs, of course. So, if a provider happens experiences a catastrophic disk failure, the publishers are running their pods on that provider would lose data, and the reputation system would do nothing to tangibly help them -- they would just be the unlucky few who drew the short straws. And even though such a system would end up punishing the provider in the long term, that would only affect future publishers looking for a provider.

Unfortunately, catastrophic faults are exactly what we are trying to prevent by using a reputation system. If a reputation system's weak point is exactly what we are trying to cover up with it, it is likely that the reputation system is not the right thing to use.

In addition, reputation systems tend to ignore the realities of force-majeure, and would just as gladly over-punish providers that happened to suffer an external mishap outside their control as they would under-punish provider that happen to game the system in some way.

## Reputations systems are already not used by real-world cloud providers

Existing real-world cloud providers might not be ideal in a lot of ways, but they have already tried many different business models to assure their customers of the reliability of their services, presumably including some kind of reputation systems. However, out of the major 3 providers (GCP, Azure, AWS), only one (AWS) cites a proven track record as a selling point on their main product page. For something touted to be universal enough to be used for ranking all providers, that is surprisingly little buy-in from the major providers.

Instead, what virtually all real-world cloud providers offer is SLA (Service-Level Agreement) contracts that specify concrete uptime (and other) guarantees, and furthermore are (at least in theory) liable for damages caused by negligent or malicious operations on part of the cloud provider itself. And while reputation and advertisement likely plays a large part in why a person would pick a particular cloud provider in the first place, such SLA contracts are what allows a business to trust the reliability of such a provider in the end. (And that brings us full-circle to the first point; reputation systems do not solve the problem of reducing counterparty risk, they solve the problem of finding a well-renowned provider.)

## Reputation systems are (arguably) how we got here in the first place

Finally, a philosophical point could be made that reputation systems and similar over-leveraged mechanisms of trusting companies that look "good enough" without any hard guarantees that their words and past actions can be trusted in the future are what got the world to where it is, and that the lack of those systems when developing smart contracts is what makes smart contracts so powerful. And while it would be cynical to claim that reputation systems never work, there are other solutions that can complement reputation systems and completely replace them when it comes to reducing counterparty risk.

## Staking as an alternative to "internet points"

In the physical world, we are able to threaten providers with violence (or forceful alienation of their possessions) unless they carry out the contracts they have previously undertaken and offset any damage they have caused.

For better or worse, in the blockchain world, at least for the time being, we cannot forcefully transfer assets away from the account that owns them. The only way a transfer can take place is if it takes place between willing participants (at the moment of the transaction).

Therefore, to implement a blockchain equivalent of SLA contracts and thus an alternative to a reputation system that tracks reliability, we have to lean on staking (aka escrow) -- basically, a smart contract (or a third party trusted by the other two) which assumes ownership of some amount of assets that the provider agrees to part with in case of a provable fault; assets that, if there is no fault, would revert back to the provider once the contract is broken.

Such a solution has a three major advantages over any reputation system:

1. First, a staking solution reimburses publishers in the event of the downtime or data loss instead of just dismissing it as a freak accident. (In that sense, it is similar to the provider underwriting an insurance on themselves, claimable by the publisher.) Thus, even if publishers were to suffer harm, at least it is not that costly for them to return to business.
2. Second, a staking solution allows providers to recover quickly from faults that were not their own errors, as long as they have similar SLAs with upstream providers of electricity, internet, etc. That is in contrast with reputation systems that would happily flag them as unreliable, and offer no way for them to get back up on the list except time.
3. Third, a staking solution allows new providers to participate in the the network right away instead of having to build up reputation for months or years. In addition, it allows providers that are experimental or otherwise known-unreliable to flag themselves as such by simply offering smaller premiums in case of faults.

That being said, a staking solution does bring a few disadvantages too:

* Staking requires committing resources in escrow and letting them just sit there for the eventual disaster recovery. However, except in the case where all of those resources are simultaneously needed for faults all across the network (i.e. unless the whole internet stops working overnight), there is always going to be some fraction of resources unused. This can probably be smoothed out by insurance contracts that are cheaper but hold only a fraction of assets in custody; however that is likely to be about as hard just making a full reputation system.
* Staking requires a larger initial commitment by providers; small and at-home providers will have a harder time providing their full capacity as they would lack the fund to underwrite all the contracts required. Then again, reputation systems also tend to put at-home providers at a disadvantage, so maybe that is okay.
* Staking assumes that all publishers want some kind of recompense in case of downtime or data loss; however it's plausible that many home publishers would be happy with just having an extra backup location or a somewhat unreliable service as long as it is more reliable than their own machines / as long as it fails independently from their own machines. That usecase would probably be better handled through something other than staking.

All that being said, staking seems like the most reasonable idea for incentivizing the reliability of trusted pod providers for now. It is further explored in the documents on [Storage](STORAGE.md) and [Uptime](UPTIME.md).

# Addendum: Note on handling force-majeure in staking contracts

Handling force-majeure (overwhelming force; i.e. external events outside the control of either party) in blockchain contracts is a relatively unexplored area; however, as Trusted Pods grows and is eventually deployed in production, it can become a major point of friction and discussion. In particular, we are likely interested in three main categories of force-majeure events: natural disasters that could not be mitigated, interruptions in service of upstream providers (e.g. internet and electricity providers), government regulators demanding a particular pod be stopped or connections to it be censored (--whether one agrees about their rights to make such demands or not, observe that a government regulator can sanction a provider with increasingly heavier penalties until they comply).

The first two cases are accidents, and could be covered by various kinds of insurance (or SLAs). However, in order to not require Trusted Pods providers to sign up for such insurance, and to give them more opportunities to diversify their service, it would make sense to allow Storage and Uptime contracts to include human-readable and human-interpreted clauses concerning what kinds of accidents are covered by the stake and what kinds of accidents are not. (For example: a large-scale multi-day electricity outage can likely bring multiple providers and various external websites down simultaneously; while some might have generators and satellite links to provide service even during such times, most of them would likely not, and might be unwilling to stake money that would be lost the first time such an outage occurs.) However, this is more of a polishing touch than a required feature, and can be introduced once providers and publishers confirm its desirability.

The last case of government regulators bears additional scrunity, however. In such cases, a sufficiently-verified official document from the regulator could be deemed sufficient proof for force-majeure; as long as providers accurately state the countries under whose jurisdictions their servers operate. Providers that are willing to stake without allowing proofs from government regulators should be allowed to host applications; however in such cases, discontinuations of a service due to government pressure would result in losing stake. Likewise, providers should be able to list multiple jurisdictions and even list non-government regulators when specifying the parameters of their contracts. Unlike the previous case, providers could provide regulatory decisions upon stopping service, as opposed to having to submit a rationale after service has stopped - hence, this might be simpler to implement in terms of smart contracts.

For determining whether a certain force-majeure occurrence is sufficient proof for a given provider to retain their stakes, an arbitrage service similar to Kleros or Reality.eth could be used; potentially, the system used could even be customizable by providers, to allow for e.g. using actual national arbitrages if/once such become available on the blockchain.
