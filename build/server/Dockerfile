FROM golang:1.19 as builder

ENV PATH /go/bin:$PATH

RUN \
  apt update && \
  apt install -y protobuf-compiler && \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /go/src/github.com/InfJoker/helicopter

COPY go.mod go.sum ./

RUN go mod download

COPY internal internal/
COPY proto proto/
COPY cmd cmd/
COPY Makefile ./

RUN make build

FROM scratch
WORKDIR /run/helicopter/
COPY --from=builder /etc/ssl/certs /etc/ssl/certs/
COPY web web/
COPY configs configs/
COPY --from=builder /go/src/github.com/InfJoker/helicopter/bin/helicopter ./helicopter

CMD [ "./helicopter", "-config", "configs/config.yml"]
