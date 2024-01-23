# <div align="center">Traffic Wizard</div>

> Traffic Wizard is your friendly load balancer written in Go _(hopefully)_. It has a very simple life goal - to balance the load of our hard-working servers. It is also very religious and follows it's own [bible](./bible.pdf) very strictly. Oh and yes, this is the first assignment for Distributed Systems course taken by [Dr. Sandip Chakraborty](https://cse.iitkgp.ac.in/~sandipc/) for Spring-2024.

### Production

1. Run `docker-compose build`
2. Run `docker-compose up`

### Debugging

To run the server:

1. Go to `server` folder.
2. Run `docker build --tag traffic-wizard-server .` to build the docker image
3. Run `docker run -e id=1 -p 5000:5000 traffic-wizard-server:latest` to run the docker container. You can set the id as you wish.

To run the loadbalancer:

1. Go to `loadbalancer` folder.
2. Run `docker build --tag traffic-wizard-lb .` to build the docker image
3. Run `docker run -p 5000:5000 traffic-wizard-lb:latest` to run the docker container.
