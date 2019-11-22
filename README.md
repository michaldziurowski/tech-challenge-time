# Tech Challenge Time

![](https://github.com/michaldziurowski/tech-challenge-time/workflows/Main%20CI/badge.svg)

## Usecases and assumptions

-   As a user, I want to be able to start a time tracking session
-   As a user, I want to be able to stop a time tracking session
-   As a user, I want to be able to name my time tracking session
-   As a user, I want to be able to save my time tracking session when I am done with it. **Assumption**: session can be saved each time user starts, stops or restarts it.
-   As a user, I want an overview of my sessions for the day, week and month. **Assumption**: week = last 7 days, month = last month
-   As a user, I want to be able to close my browser and shut down my computer and still have my sessions visible to me when I power it up again. **Assumption**: this relates to client side - since persistance layer is done in memory all saved tracking sessions will be lost after server shuts down.

## Solution

The application is made of Go based server and React based client.

### Server architecture

Server code structure is inspired by [Clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and consists of following packages:

-   `domain` - responsible for domain entities.
-   `usecases` - responsible for implementation of application usecases (`services.go`) and definition of dependencies needed to fulfill usecases (`interfaces.go`).
-   `infrastructure` - responsible for all external dependencies (http, storage, etc).

### Client architecture

Client structure consist of following folders:

-   `components` - the place where React components are defined.

    Components can be divided into two groups:

    -   `Container` components - responsible for data fetching and logic and provision of those to presentational components.
    -   `Presentational` components - responsible for defining how component should be presented.

-   `api` - here is the code responsible for communication with the server (for now server address is harcoded and of course it shouldn't be but since this is just an demo app let it stay as it is).

### Unit tests

> Unit tests in this application are ment to show an approach not to provide full test coverage.

To run the tests use following command (requires Go to be installed on the machine)

```
> go test ./...
```

### CI

Continous Integration pipeline is running on GitHub Actions. Currently there is one workflow defined:

-   `Main CI` - builds and runs unit tests for server and client on every push to master branch.

### Deployment

Currently application is not deployed anywhere.
To run it locally use docker (from root folder):

```
> docker-compose up
```

This sets up client and server. When containers are running following urls can be used to access application:

-   `http://localhost:3000` - client application
-   `http://localhost:8080/api/v1` - web api
