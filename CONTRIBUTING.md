# Contributing

We love pull requests from everyone. By participating in this project, you
agree to abide by the [code of conduct].

[code of conduct]: (CODE_OF_CONDUCT.md)

Fork, then clone the repo:

    git clone git@github.com:your-username/sqlkit.git

Set up your machine:

    dep ensure

Make sure the tests pass:

    make test

Make your change. Add tests for your change. Make the tests pass:

    rake test

Push to your fork and [submit a pull request][pr].

[pr]: https://github.com/ColDog/sqlkit/compare/

At this point you're waiting on us. We like to at least comment on pull requests
within three business days (and, typically, one business day). We may suggest
some changes or improvements or alternatives.

Some things that will increase the chance that your pull request is accepted:

* Write tests.
* Follow common go code [styles].
* Write a [good commit message][commit].

[styles]: https://github.com/golang/go/wiki/CodeReviewComments
[commit]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
