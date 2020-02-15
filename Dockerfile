# Start from a Debian image with the latest version of Go installed and a workspace (GOPATH) configured at /go.
FROM golang:latest
# Create WORKDIR (working directory) for app
WORKDIR /go/src/github.com/dadid/sportsbetting-data-api-chi_router
# Copy the local package files to the container's workspace (in the above WORKDIR)
ADD . .
# Switch WORKDIR to directory where server main.go lives
WORKDIR /go/src/github.com/dadid/sportsbetting-data-api-chi_router/cmd
# Build the go-API-template userServer command inside the container at the most recent WORKDIR
RUN go build -o userServer
# Run the userServer command by default when the container starts.
# runs command at most recent WORKDIR
ENTRYPOINT ./userServer
# Document that the container uses port 8
EXPOSE 8600
# Document that the container uses port 5432
EXPOSE 5432

# Build image with gilcrest as repository name, go-api-template as build name and latest as build tag from the current directory
# $ docker image build -t ddadi/sports-betting-api:latest .