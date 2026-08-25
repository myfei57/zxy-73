FROM golang:1.23-bookworm

ENV GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local
WORKDIR /src
COPY . .
CMD ["go", "test", "./..."]
