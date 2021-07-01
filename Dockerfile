FROM golang:1.13 as builder

LABEL maintainer="manovisut.ktp@gmail.com"

#
# Testing & Building
#

WORKDIR /src/build
ADD src /src/build
RUN go mod download
RUN CGO_ENABLED=0 go build -o main http/main.go

#
# Final Stage
#

FROM alpine:3.14

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/build/main .
CMD ["./main"]