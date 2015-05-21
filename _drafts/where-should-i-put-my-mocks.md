---
layout: post
title: Where do I put my mocks in Go?
tags: golang go mock test testing production binary stub package API godoc gomock
---

#Introduction

Go's toolchain provides a mechanism to help prevent test code from being
included in production binaries. It should be obvious why this is useful; if
stubs or mocks intended for testing make their way your production system, they
may cause it to fail silently in catastrophic and unexpected ways. This is
especially a concern if you don't have automated application and integration
tests in place. Imagine if you somehow end up with a stubbed out authentication
manager that always returns `true` when asked if a user has administration
privileges. Yikes.

This blog post will explain this mechanism, along with some of the ways you can
most effectively take advantage it, especially when dealing with mocks and
stubs.

#The Mechanism

It's well known that the `go test` tool looks in files that end in `_test.go`
for tests to run. What is less well known is that these `*_test.go` files are
_excluded_ from binaries built using `go build` and `go install`.

Excellent! This means we can put all of our testing stubs and mocks in the
`*_test.go` files and they will never accidentally be included in production
binaries!

#The Problem

This works well most of the time. You have an interface `Foo` declared in
package `foo`, and a mocked out implementation `MockFoo` also in package `foo`
but in a file named `mock_foo_test.go`. This works great, you can use `MockFoo`
instead of the real implementation for all of your tests in package `foo`!

But what if you want to use `MockFoo` for tests in packages `garply` and
`waldo`? Unfortunately you can't. Because `MockFoo` is declared in a file
ending in `_test.go`, it can only be used in tests contained in the same
package.

#The Solutions

There are a two obvious ways to work around this:

* Rename the file containing `MockFoo` from `mock_foo_test.go` to
  `mock_foo.go`. Now that it's not a test file, it can be accessed from other
packages.
* Copy `mock_foo_test.go` into the other packages that need it (modifying the
  package declaration to match its destination package).

There are same drawbacks to both approaches.

In the first approach, we lose the guarantee that `MockFoo` won't accidentally
be included the production binary.  As an additional annoyance, `MockFoo` is
now part of the public API, and as such will appear in listings produced by the
`godoc` tool.

The second approach violates the [DRY
principal](http://c2.com/cgi/wiki?DontRepeatYourself). If the interface
changes, you then have to go and change all of the mocks. If you're manually
rolling your own mocks, the process of updating them all could be error prone.

In my opinion, the second approach is the better of the two, especially when
considering that mocks can be generated automatically with tools such as
[gomock](https://code.google.com/p/gomock/). If you write a small script that
generates all of the mocks in your project using gomock, then you just have to
run it whenever you change an interface that is being mocked out. Since the
mocks are generated _from_ the interface, I don't think that the DRY principal
is of great importance here.

Happy testing!
