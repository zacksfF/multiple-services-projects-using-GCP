FROM golang:1.22.0 as builder
WORKDIR /app

# Initialize a new Go module.
RUN go mod init quickstart-go

# Copy local code to the container image.
COPY *.go ./

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -o /quickstart-go

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM gcr.io/distroless/base-debian11

# Change the working directory.
WORKDIR /

# Copy the binary to the production image from the builder stage.
COPY --from=builder /quickstart-go /quickstart-go

# Run the web service on container startup.
USER nonroot:nonroot
ENTRYPOINT ["/quickstart-go"]