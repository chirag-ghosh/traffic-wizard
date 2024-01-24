# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /server

# setting the destination for COPY
COPY go.mod ./
RUN go mod download

# Copying all the .go files
COPY *.go ./

# build
RUN go build -o /traffic-wizard-server

# the application is going to listen in the port 3002
EXPOSE 3002

# run
CMD ["/traffic-wizard-server"]