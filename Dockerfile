FROM golang:1.19-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/challenge

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Test
RUN CGO_ENABLED=0 go test -v

# Build
RUN go build -o ./out/challenge .

FROM alpine:3.9

# Copy executable from previous stage
COPY --from=build_base /tmp/challenge/out/challenge /app/challenge

ENTRYPOINT ["/app/challenge"]