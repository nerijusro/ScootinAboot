# Scootin with Nerijus Aboot

## Table of Contents

- [Prerequisites](#prerequisites)
- [Running the project locally](#running-the-project-locally)
- [Running the project using Docker](#running-the-project-using-docker)
- [Running the tests](#running-the-tests)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Method: `GET`, URL: `/client/auth`](#get-client-auth)
  - [Method: `GET`, URL: `/admin/auth`](#get-admin-auth)
  - [Method: `POST`, URL: `/client/users`](#post-client-users)
  - [Method: `POST`, URL: `/admin/scooters`](#post-admin-scooters)
  - [Method: `GET`, URL: `/admin/scooters`](#get-admin-scooters)
  - [Method: `GET`, URL: `/client/scooters`](#get-client-scooters)
  - [Method: `GET`, URL: `/client/scooters/:id`](#get-client-scooters)
  - [Method: `POST`, URL: `/client/trips`](#post-client-trips)
  - [Method: `PUT`, URL: `/client/trips/:id`](#put-client-trips)

## Prerequisites
There was a seperate tool used to run database migrations, so in order to run the project, the following is needed:
https://github.com/golang-migrate/migrate/tree/v4.17.0/cmd/migrate

## Running the project locally
Before following through, make sure that you have a MySql database running on your machine, or feel free to swap any other storage under `/db`.

Once storage is running, create a database (solution expects `scootin_aboot`- config file must be edited in case other name is prefered) and run the migration:
```
make migrate-up
```
That will create neccessary tables. Once that is done, run the project with:
```
make run
```

## Running the project using Docker
To run the server on docker, build the `docker-compose.yml` file using Terminal in the project's directory with:
```
docker-compose build
```

And run it:
```
docker-compose -d up
```

MySql database is launched together with the server, so no need to launch the migrations seperately.

## Running the tests
Test can be launched using command:
```
make test
```

## Authentication
Since assignment was kind enough to only require a static api key for authentication, it must be attached to a header as `x-api-key` for every request (except `auth`). As one's eye might catch, endpoints are grouped into `admin` and `user`. These groups have different api keys, though `admin` one can be used to call `client` endpoints as well.

## Endpoints
The project consists of the endpoints listed below:
### Method: `GET`, URL: `/client/auth`
Returns a static api key for `client` route group.

Example response:
```
{
    "StaticApiKey": "my_static_user_api_key"
}
```

### Method: `GET`, URL: `/admin/auth`
Returns a static api key for `admin` route group.

Example response:
```
{
    "StaticApiKey": "my_static_user_api_key"
}
```

### Method: `POST`, URL: `/client/users`
Creates a new user. For the sake of simplicity, only user's name is required at this project's stage. 
IMPORTANT: The `id` returned must be attached to a header as `client-id` for trip related endpoints.

Example request:
```
{
    "full_name": "Post Malone"
}
```
Example response:
```
{
    "id": "76341b35-ffb0-4ed6-b017-395b2156de99",
    "full_name": "Post Malone",
    "is_eligible_to_travel": true
}
```

### Method: `POST`, URL: `/admin/scooters`
Creates a new scooter.

Example request:
```
{
    "is_available": true,
    "location": {
        "latitude": 54.1234,
        "longitude": 25.5436
    }
}
```
Example response:
```
{
    "id": "6651ecbd-0d85-47c0-a30b-7c8598148ac8",
    "location": {
        "latitude": 54.1234,
        "longitude": 25.5436
    },
    "is_available": true
}
```

### Method: `GET`, URL: `/admin/scooters`
Returns all existing scooters. Since there is not a single use case where a mobile user could need it, it is only accessible to `admin`.

Example response:
```
{
    "scooters": [
        {
            "id": "097975dd-41cf-4c94-ae48-66ddb5f58fc6",
            "location": {
                "latitude": 54,
                "longitude": 24
            },
            "is_available": true
        },
        {
            "id": "6651ecbd-0d85-47c0-a30b-7c8598148ac8",
            "location": {
                "latitude": 54.1234,
                "longitude": 25.5436
            },
            "is_available": true
        }
    ]
}
```

### Method: `GET`, URL: `/client/scooters`
Returns scooters according to the following search criteria:
- `availability`: Used for filtering scooters that are currently free or used in an active trip. Only valid options are: `all`, `available` and `unavailable`.
- `x1` and `x2`: Scooters are searched in a rectangular area, so `x1` and `x2` indicates longitude range or x-axis projection.
- `y1` and `y2`: Y-axis points, which creates a latitude interval.
IMPORTANT: Assume that rectangular is being drawn from left to right and bottom to top. Accordingly, `x2` and `y2` values must to be greater than `x1` and `y1`. Failing to do so results in a validation error.

Example query:
```
localhost:8080/client/scooters?availability=available&x1=10.0&x2=30.0&y1=54.0&y2=55.0
```
Example response:
```
{
    "scooters": [
        {
            "id": "097975dd-41cf-4c94-ae48-66ddb5f58fc6",
            "location": {
                "latitude": 54,
                "longitude": 24
            },
            "is_available": true
        },
        {
            "id": "6651ecbd-0d85-47c0-a30b-7c8598148ac8",
            "location": {
                "latitude": 54.1234,
                "longitude": 25.5436
            },
            "is_available": true
        }
    ]
}
```
 
### Method: `GET`, URL: `/client/scooters/:id`
Returns a scooter by it's id.

Example query:
```
localhost:8080/client/scooters/6651ecbd-0d85-47c0-a30b-7c8598148ac8
```
Example response:
```
{
    "id": "6651ecbd-0d85-47c0-a30b-7c8598148ac8",
    "location": {
        "latitude": 54.1234,
        "longitude": 25.5436
    },
    "is_available": true
}
```

### Method: `POST`, URL: `/client/trips`
Creates a new trip. Id of a scooter is provided in a request body and `clientId` needs to be attached to a header as `client-id`.

Example request:
```
{
    "scooter_id": "6651ecbd-0d85-47c0-a30b-7c8598148ac8",
    "created_at": "2024-04-26T17:07:40.284Z"
}
```
Example response:
```
{
    "trip_id": "b6e7b1b3-685e-4982-82ed-437e651d111b",
    "event_type": "start_trip_event",
    "location": {
        "latitude": 54.1234,
        "longitude": 25.5436
    },
    "created_at": "2024-04-26T17:03:32.99Z",
    "sequence": 1
}
```

### Method: `PUT`, URL: `/client/trips/:id`
Used for making updates on trip's state and scooter's geographical coordinates. IMPORTANT: Once trip's status is updated to `"isFinished": true`- trip is considered over and no further updates are accepted.

Example query:
```
localhost:8080/client/trips/b6e7b1b3-685e-4982-82ed-437e651d111b
```
Example request:
```
{
    "location": {
        "latitude": 55.43,
        "longitude": 25.234
    },
    "created_at": "{{timestamp}}",
    "is_finishing": true,
    "sequence": 3
}
```
Example response:
```
{
    "trip_id": "b6e7b1b3-685e-4982-82ed-437e651d111b",
    "event_type": "end_trip_event",
    "location": {
        "latitude": 55.43,
        "longitude": 25.234
    },
    "created_at": "2024-04-26T17:07:40.284Z",
    "sequence": 3
}
```
