---
name: Release tracker
about: Create an issue to track release progress

---

<!-- < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < ☺ 
v                            ✰  Thanks for opening an issue! ✰    
v    Before smashing the submit button please review the template.
v    Word of caution: poorly thought-out proposals may be rejected 
v                     without deliberation 
☺ > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > >  -->

## QA

- [ ] Tutorial tests verification
- [ ] Test `serve` on suite of chains

### Backwards compatibility

<!-- List of tests that need be performed with previous
versions of cli to guarantee that no regression is introduced -->


### Other testing

## Migration 

<!-- Link to migration document -->

## Checklist

<!-- Remove any items that are not applicable. -->

- [ ] Branch off main to create release branch in the form  of `release/vx.y.z`.
- [ ] Add branch protection rules to new release branch.

## Post-release checklist

- [ ] Update [`changelog.md`](https://github.com/ignite/cli/blob/main/changelog.md)
- [ ] Update [version matrix](https://github.com/ignite/cli/blob/main/readme.md) in `readme.md`.
- [ ] Update docs site:
  - [ ] Add new release tag to [`docs/versioned_docs`](https://github.com/ignite/cli/tree/main/docs/versioned_docs).
- [ ] After changes to docs site are deployed, check [docs.ignite.com/](https://docs.ignite.com/) is updated.

____
