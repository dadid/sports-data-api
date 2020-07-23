FROM golang:latest
# Create working directory for app
WORKDIR /dadid/sports-data-api
# Copy the local package files to the container's workspace (in the above WORKDIR)
ADD . .
# Switch WORKDIR to directory where server main.go lives
WORKDIR /dadid/sports-data-api/cmd
# Build the go executable inside the container at the most recent WORKDIR
RUN go build -o server
# Run server executable on container start, runs command at most recent WORKDIR
ENTRYPOINT ./server
EXPOSE 8600
EXPOSE 5432

# docker image build -t ddadi/sports-data-api:latest .