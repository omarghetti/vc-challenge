## VC Challenge - Search API

# Tools
this project has been developed with:
- Redis
- Go 1.22

# How to setup the environment
To setup the environment, launch:

`make setup-env`

# How To Run locally
To run the service locally, you need to launch a local redis instance.

after satisfiyng the requirement, just launch

`make debug`

# How To Run With Docker-Compose
To run both a prebuild Redis image and Api image on docker using compose, just run

`make run=dependencies`

# How To Run Tests
To run Tests, launch:

`make test`

# How Would I Rank results returned from the API?

At the moment, there is not a proper ranking in the current API.

