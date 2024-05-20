FROM golang:1.22-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY *.env ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /oracle-password-manager


# Run
CMD ["/oracle-password-manager"]