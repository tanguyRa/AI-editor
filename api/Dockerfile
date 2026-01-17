# ---------------------
# --- Staging Stage ---
# --------------------
# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:latest AS staging

ARG TARGETOS
ARG TARGETARCH

# Create and change to the app directory.
WORKDIR /app



# -------------------------
# --- Development Stage ---
# -------------------------
FROM staging AS development

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

EXPOSE 8080
ENTRYPOINT ["air"]



# -------------------
# --- Build Stage ---
# -------------------
FROM staging AS builder

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o server ./cmd/server



# ------------------------
# --- Production Stage ---
# ------------------------

# Use the official Alpine image for a lean production container.
FROM alpine:3 AS production

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server
COPY --from=builder /app/web web

EXPOSE 8080
# Run the web service on container startup.
ENTRYPOINT ["/server", "serve", "--http=0.0.0.0:8080"]
