# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app


RUN go build -o news-aggregator-service

# Expose port 8081 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./news-aggregator-service"]
