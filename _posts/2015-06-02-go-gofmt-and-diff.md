---
layout: post
title: Go, gofmt, and diffs
tags: go, golang, gofmt, fmt, diff, git, whitespace, ignore, hg, mercurial, svn, subversion, github
---

[Gofmt](https://golang.org/cmd/gofmt/) is a Go program formatter. Its universal
adoption by the Go community has resulted in a consistent formatting style
among most if not all Go source code.

When making changes to Go code, gofmt sometimes causes changes to lines of code
that weren't originally modified. This can done for example to keep indentation
consistent. This can make diffs hard to read.

For example, the following is a diff of a simple
[type declaration](https://golang.org/ref/spec#Type_declarations):

    $ git diff
    diff --git a/main.go b/main.go
    index 2c7af41..d9e59cd 100644
    --- a/main.go
    +++ b/main.go
    @@ -1,13 +1,12 @@
     type foo struct {
    -    jabberwocky  int
    -    quux         complex128
    -    galumphing   string
    -    vorpal       float64
    -    bandersnatch float64
    -    mimsy        struct {
    -        beaming  int
    -        tumtummy string
    -        brillig  float64
    +    jabberwocky int
    +    quux        complex64
    +    tumtummy    string
    +    galumphing  string
    +    vorpal      float64
    +    mimsy       struct {
    +        beaming int
    +        brillig float64
         }
     }

Because gofmt has rewritten the source code to have consistent indentation, at
least a whitespace change has been made to every line. It can be hard to see
what the _real_ change is.

To see the real changes, the diff algorithm should be configured to ignore
whitespace. The diff below makes it quite obvious what the actual change is:

    $ git diff -b
    diff --git a/main.go b/main.go
    index 2c7af41..d9e59cd 100644
    --- a/main.go
    +++ b/main.go
    @@ -1,12 +1,11 @@
     type foo struct {
         jabberwocky int
    -    quux         complex128
    +    quux        complex64
    +    tumtummy    string
         galumphing  string
         vorpal      float64
    -    bandersnatch float64
         mimsy       struct {
             beaming int
    -        tumtummy string
             brillig float64
         }
     }

It's now quite clear that `quux` has had its type modified, `bandersnatch` has
been removed, and `tumtummy` has been moved from the inner struct, to the outer
struct.

This type of diff can typically be created by passing the
`--ignore-space-change` flag (the short version is `-b`). This works with Git,
Mercurial, Subversion, and well as the `diff` command that comes with
`diffutils`.  GUI based diffs such as those [produced by
Github](https://github.com/blog/967-github-secrets) will typically have an
ignore-whitespace option as well.
