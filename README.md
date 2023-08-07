# Fastrand
[![License][license-img]][license]
[![GoDev Reference][godev-img]][godev]
[![Go Report Card][goreportcard-img]][goreportcard]

Package fastrand provides quickly generated pseudo-random numbers with no repeatability guarantees
on the stream of values.

It uses internals of the go runtime to generate pseudo-random numbers without requiring mutexes or
syncronization. As a result, it's much faster than `math/rand` and scales very well to many cores.
It doesn't allow you to provide your own seed or generate a stable sequence of values.



[license]: https://raw.githubusercontent.com/abursavich/fastrand/main/LICENSE
[license-img]: https://img.shields.io/badge/license-mit-blue.svg?style=for-the-badge

[godev]: https://pkg.go.dev/bursavich.dev/fastrand
[godev-img]: https://img.shields.io/static/v1?logo=go&logoColor=white&color=00ADD8&label=dev&message=reference&style=for-the-badge

[goreportcard]: https://goreportcard.com/report/bursavich.dev/fastrand
[goreportcard-img]: https://goreportcard.com/badge/bursavich.dev/fastrand?style=for-the-badge
