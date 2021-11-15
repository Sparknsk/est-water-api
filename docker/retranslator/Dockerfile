ARG GITHUB_PATH=github.com/ozonmp/est-water-api

FROM golang:1.16-alpine

WORKDIR /home/${GITHUB_PATH}

COPY . /home/${GITHUB_PATH}
RUN go mod download
COPY . .
RUN go build ./cmd/est-water-api/main.go

CMD ["./main"]