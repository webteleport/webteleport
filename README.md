# skynet

install:

```
$ go install github.com/btwiuse/skynet/cmd/skynet@latest
```

server:

```
$ skynet server
```

example clients:

```
$ skynet client
$ skynet gos
$ skynet echo
```

TODO

- [ ] reverseproxy WebTransport requests
- [ ] support user specified hostname, requiring netrc authentication
- [ ] support custom root domain, for example `ROOT=usesthis.app`
- [ ] support concurrent rw on session manager map
