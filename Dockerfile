FROM btwiuse/arch:golang

COPY . /ufo
WORKDIR /ufo

RUN git config --global user.name navigaid
RUN git config --global user.email navigaid@gmail.com
RUN go mod tidy
RUN GOBIN=/usr/local/bin go install ./cmd/ufo

CMD ufo gos https://ufo.k0s.io/gos
