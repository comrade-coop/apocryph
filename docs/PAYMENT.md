# Payment

When it's initially deployed, the aApp creates a new Ethereum wallet for itself and uses it to deploy a smart contract it uses to manage all payments. Crucially, the aApp is the only one able to use that Ethereum wallet, since it never shares it with anyone outside the [TEE boundary](./ATTESTATION.md) - and by extension, it is the only one able to operate the payment smart contract.

## Authorized funds

At all times, the aApp keeps track of your current USDC balance and your current USDC spending cap (allowance) for the payment contract.  
The smaller of those two values is the amount the aApp considers itself "authorized" to use—it is how much it can theoretically spend on your behalf.

## Minimal required authorization

In order to log in, the aApp requires you to have some minimal authorization—currently set at 1 USDC. This protects the aApp from abuse, as it ensures that its users have at least some funds when they attempt to store data on S3.

## Maximum overdraft

In addition to your total authorized amount, the aApp keeps track of how much USDC you currently owe it for the S3 storage you have used. It would then slowly withdraw those funds from the USDC you have authorized it to use, whenever the amount owed exceeds some threshold.

However, at some point, you might authorized USDC for the aApp. At this point, the aApp considers you to be in "overdraft" (it can spend no USDC on your behalf, but you still have data stored on the S3 storage).

**The aApp will keep overdrafted buckets around for a while (even though they are inaccessible), but once the overdraft hits a maximum threshold—currently set at -10 USDC—the aApp will delete buckets.**

Hence, you should maintain some authorized amount at all times.
