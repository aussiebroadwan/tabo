FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download the dependencies
# RUN go mod download

# Build the application 
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o backend .

FROM scratch AS deploy

# Metadata
LABEL description="The Aussie BroadWAN Keno Backend"
LABEL vendor="Aussie BroadWAN"
LABEL version="0.1.0"

# Copy the certificates and user/group files from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/backend /backend

# Run the bot
ENTRYPOINT ["/backend"]