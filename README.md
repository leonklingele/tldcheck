# Fast TLD checker

![build](https://github.com/leonklingele/tldcheck/actions/workflows/build.yml/badge.svg)

## Installation

```sh
go install github.com/leonklingele/tldcheck/cmd/tldcheck@latest
tldcheck -help
```

## Check domain name availability

```sh
tldcheck -name google -workers 32 | grep -v 'not available'
```
