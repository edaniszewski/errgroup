# errgroup

[![Build Status](https://travis-ci.org/edaniszewski/errgroup.svg?branch=master)](https://travis-ci.org/edaniszewski/errgroup)
[![GoDoc](https://godoc.org/github.com/edaniszewski/errgroup?status.svg)](https://godoc.org/github.com/edaniszewski/errgroup)
[![Go Report Card](https://goreportcard.com/badge/github.com/edaniszewski/errgroup)](https://goreportcard.com/report/github.com/edaniszewski/errgroup) 

A mechanism to manage long and short lived application goroutines with error preservation.

## Getting

```
go get github.com/edaniszewski/errgroup
```

## Why?

The [related projects](#related-projects) section below links to some projects which are similar
to this one and ultimately inspired this one. Each of them provide similar useful capabilities,
however they are simple and do not cover all use cases. This project does not aim to be an alternative
to any of them, but rather fill a small behavioral niche which they do not cover.

As with the related projects, functional groups may be defined which specify the normal logic to
execute, and the terminal logic which may be used to terminate the goroutine and do any clean up.
The departure from others is that the termination logic is only called when an error is returned.

In short, this means that a goroutine returning `nil` will not cause the entire group to terminate.

See [`example_test.go`](example_test.go) for some usage examples.

## Related Projects

Below are some related projects which accomplish similar things, albeit with different behaviors.
If this project is too niche for you, check them out to see if they are more suited to your use
case(s). A big thanks goes out to their authors, as these projects served as the inspiration and
starting point for this project.

- https://github.com/oklog/run
- https://gopkg.in/tomb.v2
- https://godoc.org/golang.org/x/sync/errgroup
