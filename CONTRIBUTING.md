# Contributing

## Getting started
Before you begin:
- Apocryph repositories use [C# (.NET 5.0)](https://dotnet.microsoft.com/download) and [Solidity](https://ethereum.org/en/developers/local-environment/);
- Check the README / documentation of the repository to which you are interested in contributing;
- Setup your local environment according to the README / documentation.

## How to contribute?
When contributing to this repository, please first pick an issue marked as "open to contributors" or initiate a discussion using a new issue. You can check the suitable issues for the repository or check the [Apocryph GitHub Project](https://github.com/orgs/comrade-coop/projects/1) for overview of all open issues across all repositories.

### What are "open to contributors" issues?
All of the work in the project is organized using milestones and issues. Certain issues are marked as "open to contributors". These issues are specifically curated by the maintainers as suitable for contributors. 

The issues have the following structure:
- **description** including goals and relevant notes
- well defined **acceptance criteria**
- **story points** as part of the issue description

The story points gives notation on the complexity of the issue and can be indicative for the estimated time duration. The guidance for assigning story points is the following:
- **1** and **2** story points are reserved for minor tasks that take less than **4 hours**;
- **3** points are reserved for tasks that will likely take **a day**;
- **5** points are for tasks that will take at least a day, but no more than **two days**;
- **8** points are for tasks that will likely take **a few of days**;
- **13** points are for tasks that will likely take **a week**.

### Ready to make a change? Fork the repo and prepare a pull request
All contributions to the repo are managed using pull requests, reviewed and approved by the maintainers. The pull requests should have the following structure:
- the name of the request should be `<issue-title>`
- all of the changes should be squashed into a single commit before submitting the pull request
- the message of the single commit should follow [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) structure, where the body should contain link to the issue number and the issue story points: `close #<issue-number> for <story-points> story points`  
- the single commit should be signed using a [GPG key verifiable by GitHub](https://docs.github.com/en/github/authenticating-to-github/managing-commit-signature-verification/signing-commits).

Often you will have to implement changes to your open pull request. In this situation, we recommend that you ammend the single commit used for the pull request (using `git commit --amend`). Have in mind that after amending your commit, you will have to force-push it (using `git push --force`).

If you open a pull request without an issue or in relation to an issue that is not marked as "open to contribution" and your pull request is considered as valueable contribution by the maintainers, they will create / update the corresponding issue and they will ask you to ammend your commit message and pull request title.

## Rewarding contributions
Apocryph project is using a decentralized governence model based on [DAO](https://en.wikipedia.org/wiki/Decentralized_autonomous_organization) concept. Apocryph DAO allocates part of Apocryph token for project contributors. On a regular basis, part of that token is be distributed to contributors according to story points of their contributions over period of time. 
