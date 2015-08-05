# Contribution Guide

We welcome and encourage community contributions to `rack`.

Please familiarize yourself with the Contribution Guidelines and Project Roadmap before contributing.

There are many ways to help `rack` besides contributing code:

 - Find bugs or file issues
 - Improve the [documentation](https://github.com/rackspace/rack/tree/master/docs)

## Contributing Code

Unless you are fixing a known bug, we **strongly** recommend discussing it with the core team via a GitHub issue before getting started to ensure your work is consistent with `rack`'s roadmap and architecture.

All contributions are made via pull request. Note that **all patches from all contributors get reviewed**. After a pull request is made other contributors will offer feedback, and if the patch passes review a maintainer will accept it with a comment. When pull requests fail testing, authors are expected to update their pull requests to address the failures until the tests pass and the pull request merges successfully.

At least one review from a maintainer is required for all patches (even patches from maintainers).

Reviewers should leave a `LGTM` or `:+1:` comment once they are satisfied with the patch. If the patch was submitted by a maintainer with write access, the pull request should be merged by the submitter after review.

## Code Style

Please follow these guidelines when formatting source code:

* Go code should match the output of `gofmt -s`
* Shell scripts should adhere to the [Google Shell Style guide](https://google-styleguide.googlecode.com/svn/trunk/shell.xml)

## Pull request procedure

To make a pull request, you will need a GitHub account; if you are unclear on this process, see GitHub's documentation on [forking](https://help.github.com/articles/fork-a-repo) and [pull requests](https://help.github.com/articles/using-pull-requests). Pull requests should be targeted at the `master` branch. Before creating a pull request, go through this checklist:

1. Create a feature branch off of `master` so that changes do not get mixed up.
1. [Rebase](http://git-scm.com/book/en/Git-Branching-Rebasing) your local changes against the `master` branch.
1. Run the full project test suite with the `go test ./...` (or equivalent) command and confirm that it passes.
1. Run `gofmt ./...`.
1. Ensure that each commit has a subsystem prefix (ex: `files: `).

Pull requests will be treated as "review requests," and maintainers will give feedback on the style and substance of the patch.

Normally, all pull requests must include tests that test your change. Occasionally, a change will be very difficult to test for. In those cases, please include a note in your commit message explaining why.

## Conduct

Whether you are a regular contributor or a newcomer, we care about making this community a safe place for you and we've got your back.

* We are committed to providing a friendly, safe and welcoming environment for all, regardless of gender, sexual orientation, disability, ethnicity, religion, or similar personal characteristic.
* Please avoid using nicknames that might detract from a friendly, safe and welcoming environment for all.
* Be kind and courteous. There is no need to be mean or rude.
* We will exclude you from interaction if you insult, demean or harass anyone. In particular, we do not tolerate behavior that excludes people in socially marginalized groups.
* Private harassment is also unacceptable. No matter who you are, if you feel you have been or are being harassed or made uncomfortable by a community member, please contact a member of the `rack` core team immediately.
* Likewise any spamming, trolling, flaming, baiting or other attention-stealing behaviour is not welcome.

We welcome discussion about creating a welcoming, safe, and productive environment for the community. If you have any questions, feedback, or concerns [please let us know](mailto:sdk-support@rackspace.com).

## Attribution

This contributing guide was directly based off of the [Flynn contributing guide](https://github.com/flynn/flynn/blob/master/CONTRIBUTING.md). Thanks to the Flynn team for providing a great example!