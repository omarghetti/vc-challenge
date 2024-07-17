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
after doing that, just make a copy
of `app.env.example` and name it `app.env`.

after satisfiyng the requirement, just launch

`make debug`

# How To Run With Docker-Compose
To run both a prebuild Redis image and Api image on docker using compose, just make a copy
of `app.env.example` and name it `app.env`. after doing that, just run

`make run=dependencies`

# How To Run Tests
To run Tests, launch:

`make test`

# How Would I Rank results returned from the API?

At the moment, there is not a proper ranking in the current API.

One possibility would be to rank the results based on the frequency of the searched terms inside the indexed document


