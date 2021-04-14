# Fast TLD checker

## Installation

```sh
go get -u github.com/leonklingele/tldcheck/...
tldcheck -help
```

## Check domain name availability

```sh
tldcheck -name google -workers 32 | grep -v 'not available'
```
