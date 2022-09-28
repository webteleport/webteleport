```
$ cat <(printf "GET / HTTP/1.1\nHost: 7.ufo.k0s.io\nAccept: text/event-stream\n\n") - | openssl s_client -connect 7.ufo.k0s.io:443
```
