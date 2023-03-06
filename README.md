<h1 align="center">Sports News</h1>
<p align="center">A microservice that processes data from external feed providers.</p>


<p align="center">

<a style="text-decoration: none" href="https://github.com/KarolosLykos/sportsnews/actions?query=workflow%3ABuild+workflow%3Abuild">
<img src="https://img.shields.io/github/actions/workflow/status/KarolosLykos/sportsnews/build.yml?style=flat-square" alt="Build Status">
</a>

<a style="text-decoration: none" href="go.mod">
<img src="https://img.shields.io/badge/Go-v1.19-blue?style=flat-square" alt="Go version">
</a>

<a href="https://codecov.io/gh/KarolosLykos/sportsnews" style="text-decoration: none">
<img src="https://img.shields.io/codecov/c/github/KarolosLykos/sportsnews?color=magenta&style=flat-square"alt="Downloads">
</a>

## Application Structure

---
- `Cmd` folder contains the starting point of the application.
- `Internal/server` folder contains the initialization of the service, starts the consumer and the http router.
- `Internal/article` folder contains interfaces and implementations to interact with the `article` domain.
- `Internal/domain` folder contains the article model domain.
---

## Requirements
- Go 1.19+
---
## Dependencies
- [github.com/labstack/echo/v4](https://echo.labstack.com/) HTTP router
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) Structured logger
- [github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) Package envconfig implements decoding of environment variables.
- [go.mongodb.org/mongo-driver](https://https://github.com/mongodb/mongo-go-driver) The MongoDB supported driver for Go.
- [github.com/golang/mock](https://github.com/golang/mock) Mocking framework
- [github.com/stretchr/testify](https://github.com/stretchr/testify) Testing Library
- [github.com/jarcoal/httpmock](https://github.com/jarcoal/httpmock) Easy mocking of http responses from external resources.
### Extras
- [github.com/go-redis/redis/v8](https://github.com/redis/go-redis) Redis go client

## Run Instructions
*If you keep the default configuration the microservice should be running on port :8081*

- Local

Make sure you have mongo and redis running or run all containers
```shell
make run-all
```

```go
go run cmd/main.go
```

Clean
```shell
make clean
```

- Alternatively, you can build and run a composed Docker setup.
```shell
docker compose up --build -d
```

Clean
```shell
docker compose down
```

## Run tests
```shell
make test
```
## Endpoints

<details>

## List Articles
Example request:

```bash
curl -X GET http://localhost:8081/api/v1/articles
```

Example Response:

200 Status OK
```
{ 
  "status":"success",
  "data": [{"id":"640641f4b1bc7afc5cd2f855",...},{"id":"640641f4b1bc7afc5cd2f855",...}...]
}
```

## Get Article By ID
```bash
curl -X GET http://localhost:8081/api/v1/articles/640641f4b1bc7afc5cd2f855
```

Example Response:

200 Status OK
```
{ 
  "status":"success",
  "data": {"status":"success","data":{"id":"640641f4b1bc7afc5cd2f855",...}}
}
```
</details>

### Assumptions/Extensions
- Implemented the microservice using clean architecture design, to have a clear separation of the business logic and delivery mechanisms.
- I decided to create workers for processing each article after getting the list of articles.
- Probably I should not use hardcoded values (articlesCollection, maxWorkers,... etc).
- Basic validation.
- Didn't use transactions.
- Added redis cache.
- Unit testing using mocks for the different "adapters".