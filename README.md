## Docker Stalker

This project is for making docker container management easier. Run the docker-stalker image, and be amazed at how you can suddenly see all your containers. 

### Backend

Backend is a lightweight golang server that connects to the docker daemon, and exposes API's for retrieving running docker processes. 

It takes in the following environment variables:
 
* `PASSWORD`: Optional password for accessing main docker stalker dashboard
* `PORT`: Backend HTTP Server Port
* `TOKEN_EXPIRY_MILLI`: Optional token expiration time before user is prompted to enter in the password again (default is 6 hours)
* `APP_BUILD_FOLDER`: Folder where compiled react app can be found (default is development configuration, meaning it's assumed you're running inside of /back)

If a password environment variable is not set, the docker stalker dashboard will by default be available to everybody (WARNING: this is dangerous, docker stalker has restart capabilities and exposes *all* environment variables for a container)

### Frontend

Frontend is the docker stalker dashboard, it takes in the following environment variables during build time*

* `REACT_APP_API_SERVER`: URL of Backend App Server