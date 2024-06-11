# Use the official Golang image
FROM golang:1.22.4

# Set the Current Working Directory inside the container
WORKDIR /go/src/app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all the dependencies
RUN go mod download

# Command to run the executable
CMD ["go", "run", "cmd/app/main.go"]
