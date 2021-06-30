# Contributing guidelines

If you're willing to create a new PR on Starport, please make sure that you read and comply this document.

Start a new [Discussion](https://github.com/tendermint/starport/discussions/new) if you want to propose changes to this document.

Thank you for your contribution!

## Opening pull requests 

### Your pull request title
Avoid long names in your titles. Make sure your title does not have more than 60 characters.

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to find the best title.

Use brackets to point the package or feature you worked on. E.g.: `feat(services/chain)`, `fix(scaffolding)`.

### Review your own code
Make sure that you manually tested the changes you're introducing before creating a PR or pushing another commit.

Keep an eye on your PR and make sure that all CI checks becomes green.

### Avoid rebasing commits in your branch 
Avoid rebasing after you open your PRs to reviews instead, add more commits. It's OK to do force pushes if it's in the draft mode and never opened to reviews before.

A reviewer would like see a linear commit history while reviewing. If you tend to force push to an older commit, reviewer will lose track in your recent changes and will have to start reviewing from scratch.

Don't worry about adding too many commits. They will be squashed into a single commit while merging. Your PR's title will be used as the commit message.

### Ask someone else's help
If you started a PR but couldn't finish it because you lack off time, don't hesisate to ask for help from someone else to take over the ownership.

We appreciate every bits of your work!

## Coding style
