# ufo

Listening on TCP port 0 makes the system allocate a random port for you.

```
...
  ln, err := net.Listen("tcp", ":0")
  if err != nil {
          return err
  }
  log.Println("listening on", ln.Addr().String())
  return http.Serve(ln, http.FileServer(http.Dir(".")))
...

2022/09/24 23:09:17 listening on ::48982
```

Taking inspiration from that, ufo is a Golang service that allocates random
subdomains to clients with automatic SSL, HTTP/3 for free.

The programming interface is the almost the same as [net.Listen]

```
...
  ln, err := ufo.Listen("https://ufo.k0s.io")
  if err != nil {
          return err
  }
  log.Println("listening on", ln.URL())
  return http.Serve(ln, http.FileServer(http.Dir(".")))
...

2022/09/24 23:09:17 listening on https://1.ufo.k0s.io
```

deploy server:

```
$ kubectl apply -f https://raw.githubusercontent.com/btwiuse/ufo/main/deploy.yaml
```

install client:

```
$ go install github.com/btwiuse/ufo/cmd/ufo@latest
```

example apps:

```
$ ufo hello
$ ufo gos
$ ufo sse
$ ufo echo
```

TODO

- [x] support concurrent rw on session manager map
- [x] support custom root domain, for example `HOST=ufo.dev`
- [x] support user specified hostname
- [ ] authentication plugins: netrc, basic auth, environment variable
- [ ] reverseproxy WebTransport requests
