# quichost

install:

```
$ go install github.com/btwiuse/quichost/cmd/quichost@latest
```

server:

```
$ quichost server
```

example clients:

```
$ quichost client
$ quichost gos
$ quichost echo
```

TODO
- [ ] reverseproxy WebTransport requests
- [ ] support user specified hostname, requiring netrc authentication
- [ ] support custom root domain, for example `ROOT=usesthis.app`
- [ ] support concurrent rw on session manager map
