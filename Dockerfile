#FROM golang:1.16-alpine AS builder
FROM golang:alpine

RUN apk update && apk add --no-cache git

# Move to working directory (/app).
WORKDIR /app

# Copy the code into the container.
COPY . .

# Download all the dependencies that are required in your source files and update go.mod file with that dependency.
# Remove all dependencies from the go.mod file which are not required in the source files.
RUN go mod tidy

# Build the application server.
RUN go build -o binary .

# Command to run when starting the container.
ENTRYPOINT ["/app/binary"]
