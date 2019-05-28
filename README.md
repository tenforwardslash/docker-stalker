## Docker Stalker

This project is for making docker container management easier. Run the docker-stalker image, and be amazed at how you can suddenly see all your containers. 

### Deployment Instructions

Login to the server where you'd like to run Docker Stalker. The image can be found [here](https://hub.docker.com/r/10forward/docker-stalker). To start the dashboard with no password, run the following command: 

```sh

docker run -v //var/run/docker.sock:/var/run/docker.sock -p 8080:8080 10forward/docker-stalker

```

You should now see the docker dashboard running on port 8080!

### Build
To build an image locally titled `10forward/docker-stalker`, run `make`. 

### Usage

We use this [really magical nginx library](https://github.com/jwilder/nginx-proxy). It's cool because with 2 environment variables and a network definition added to your **already running** docker containers, a subdomain can be routed to the container of your choice. Comes in really handy when you only have one server but want to run a bunch of different services on various subdomains. 

Full circle, let's say I want to run docker stalker on stalker.OUR_SITE.org. Here is the docker-compose.yml file used to run the dashboard at stalker.OUR_SITE.org: 

```
version: '3'

# this ensures that our container is on the same network as the nginx-proxy library 
# if this is omitted, our nginx-proxy won't work! Check out the linked repo 
# for additional configuration options and details
networks:
  default:
    external:
      name: nginx-proxy

services:
  stalker:
    image: "10forward/docker-stalker"
    volumes:
      # this volume mapping is important, without it docker-stalker won't work
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      # these two environment variables starting with VIRTUAL are a part of nginx-proxy config
      VIRTUAL_PORT: 8080
      VIRTUAL_HOST: stalker.OUR_SITE.org
      # The password to access the dashboard here is `this_is_amazing`
      # if this env var is omitted, there will be no password
      PASSWORD: this_is_amazing

```      

run `docker-compose up -d` in the folder containing the above yaml file and BOOM your dashboard is now up. If you're also using [nginx-proxy](https://github.com/jwilder/nginx-proxy), don't forget to login to your domain provider and create a new A record mapping your subdomain to the IP address of your server!

### Backend Explained

Backend is a lightweight golang server that connects to the docker daemon, and exposes API's for retrieving running docker processes. 

It takes in the following environment variables:
 
* `PASSWORD`: Optional password for accessing main docker stalker dashboard
* `PORT`: Backend HTTP Server Port
* `TOKEN_EXPIRY_MILLI`: Optional token expiration time before user is prompted to enter in the password again (default is 6 hours)
* `APP_BUILD_FOLDER`: Folder where compiled react app can be found (default is development configuration, meaning it's assumed you're running inside of /back)

If a password environment variable is not set, the docker stalker dashboard will by default be available to everybody (WARNING: this is dangerous, docker stalker has restart capabilities and exposes *all* environment variables for a container)

### Frontend Explained

Frontend is the docker stalker dashboard, it takes in the following environment variables during build time. See `front/src/Constants/index.js` for usage 

* `REACT_APP_API_SERVER`: URL of Backend App Server