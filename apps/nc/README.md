```
$ cat <(printf "GET / HTTP/1.1\nHost: 6.ufo.k0s.io\n\n") - nc 6.ufo.k0s.io 80
GET / HTTP/1.1
Host: 6.ufo.k0s.io

asdf
asdf
asdf
asdf
asdf
asdf
hello
hello
```
