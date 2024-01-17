# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /lb

# setting the destination for COPY
COPY go.mod ./
RUN go mod download

# Copying all the .go files
COPY . /lb

# build
RUN go build -o /traffic-wizard-lb

# the application is going to listen in the port 5000
EXPOSE 5000

# run
CMD ["/traffic-wizard-lb"]