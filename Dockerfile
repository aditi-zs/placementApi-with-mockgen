# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Download the MySQL driver for Go
RUN go get -u github.com/go-sql-driver/mysql

# Build the Go application
RUN go build -o main .

# Expose the port your Go application will listen on
EXPOSE 8080

# Command to run your Go application
CMD ["./main"]
