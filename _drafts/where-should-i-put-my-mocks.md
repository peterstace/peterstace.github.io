---
layout: post
title: (Working Title) Where do I testing helpers, fakes, and mocks in Go?
---

Sometimes when testing, it can be useful to have:

* testing frameworks
* test helper functions
* test fixtures
* fakes
* mocks
* stubs
* etc.

For convenience, I will refer to all of these as _testing utilities_.

These pieces of code are useful when testing, but for reasons that should be
obvious, you don't want them anywhere near your production artifacts. This is
especially the case for mocks, fakes, and stubs, since these things implement
interfaces that are actually used in production code.

Sure, you can say "that could never happen to us". And you're right, it won't
happen... until it does.

Go and its toolchain imposes a fairly rigid method of storing test code. The
tests themselves must be contained in files that end in `_test.go`. If they're
not, then `go test` won't be able to find the tests. Files that end in
`_test.go` are only included in builds made using `go test`. These files are
_not_ included in normal builds produced using `go install` or `go build`.

Because files ending in `_test.go` aren't compiled into the production
artifacts, they _seem_ like the perfect place to put testing utilities. If
you're testing a struct called `foo`, then `foo_test.go` is a great place to
put the tests, any testing fixtures required, along with any helper functions
specific to `foo`.

This breaks down when testing utilities are needed by the tests for more than
one package. Identifiers defined in files ending in `_test.go` are only visible
from other files if and only if:

* the other file ends in `_test.go`, and
* the other file is in the same package.

In order for testing utilities to be used in tests from two different packages,
the testing utilities must _not_ be in files that end in `_test.go`, which
removes some of the protection against accidentally including test code in
production artifacts. This is the approach taken by
[some](https://code.google.com/p/gomock/) [common](https://labix.org/gocheck)
[testing](http://goconvey.co/) [frameworks](http://onsi.github.io/ginkgo/).

I consider accidentally including mocks in production artifacts as to posing a
different level of risk compared to doing the same for testing frameworks. Say
you have an interface `Baz` in package `quux`. `Baz` is used as a dependency in
packages `waldo` and `garply`. `MockBaz`, a mock of `Baz` by definition
satisfies the `Baz` interface. You _really_ don't want the mock to be included
in any production artifact. But in order to allow tests in `waldo` and `garply`
access to the mock , it must live in a file that _doesn't_ end in `_test.go`
(probably called `mock_baz.go` or similar).

An alternate approach would be to include a copy of `MockBaz` in both `waldo`
and `garply`. `MockBaz` can be placed in files named `mock_baz_test.go` (one
copy in each package). This brings back the additional layer of safty that
prevents test code ending up in production. The downside is that multiple
copies of `MockBaz` exist, violating the DRY principal.  However I don't think
it matters so much in this case since the repeated code is generated.
