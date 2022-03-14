#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/bin/app
COPY . .
RUN go mod tidy
RUN go build -o /go/src/app 
#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app/configs/ /configs
COPY --from=builder /go/bin/app/internal/api/docs/ /internal/api/docs
COPY --from=builder /go/src/app /app
ENTRYPOINT /app