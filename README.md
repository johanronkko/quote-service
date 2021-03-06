# Quote Service

Provides external customers with an API to send shipment information to Sendify, and for us to let them know how much the shipment will cost. 

## How to Run

Here are the steps to start the application:

1. `make quote` to build quote-api image.
2. `make up` to run API, PostgreSQL and pgAdmin.
3. `make logs` to see docker-compose logs.

If it is the first time you run the application:

4. `make migrate` to create the PostgreSQL schemas. 

Optionally, seed the database with `make seed`.

## Endpoints

Every endpoint will respond with an HTTP status `code` and a `success` indicator. If an error happens, the response body will consist of an `error` field. All data is sent via a `data` field. See below for examples of response bodies for the endpoints. `cmd/quote-api/handler/handler_test.go` also serves as documention.

### Healthcheck

Do `GET http://localhost:3000/api.v1/healthcheck`.
  
### Quote by ID  

Do `GET http://localhost:3000/api.v1/quotes/1cf37266-3473-4006-984f-9325122678b7` and expect a respons with the following format.

```json
{
    "code": 200,
    "data": {
        "quote": {
            "id": "1cf37266-3473-4006-984f-9325122678b7",
            "to": {
                "name": "Sven Svensson",
                "email": "sven.svensson@example.com",
                "address": "Teststreet 42A, CityA 12345",
                "country_code": "SE"
            },
            "from": {
                "name": "John Doe",
                "email": "john.doe@example.com",
                "address": "Teststreet 42B, CityB 12345",
                "country_code": "US"
            },
            "weight": 45,
            "shipment_cost": 1250
        }
    },
    "success": true
}
```

### List Quotes  

Do `GET http://localhost:3000/api.v1/quotes/` and expect a response body with the following format.

```json
{
    "code": 200,
    "data": {
        "quotes": [...]
    },
    "success": true
}
```

### Add quote:

Do `POST http://localhost:3000/api.v1/quotes/` with a request body of the following format.

```json
{
    "to": {
        "name": "Hmm Hmmson",
        "email": "hmm.hmmson@example.com",
        "address": "Teststreet 11A, Xcity 55555",
        "country_code": "FR"
    },
    "from": {
        "name": "Wihh a",
        "email": "wihh.a@example.com",
        "address": "Galzstreet 1B, GalzB 77777",
        "country_code": "NO"
    },
    "weight": 301
}
```

Expect the following response body.

```json
{
    "code": 201,
    "data": {
        "quote": {
            "id": "4d1046a6-647d-4d33-b31c-025c80fdaa02",
            "to": {
                "name": "Hmm Hmmson",
                "email": "hmm.hmmson@example.com",
                "address": "Teststreet 11A, Xcity 55555",
                "country_code": "FR"
            },
            "from": {
                "name": "Wihh a",
                "email": "wihh.a@example.com",
                "address": "Galzstreet 1B, GalzB 77777",
                "country_code": "NO"
            },
            "weight": 301,
            "shipment_cost": 2000
        }
    },
    "success": true
}
```

A request with bad _email_ and _name_ fields could result in the following response body.

```json
{
    "code": 400,
    "error": [
        {
            "field": "name",
            "error": "Key: 'NewQuote.to.name' Error:Field validation for 'name' failed on the 'personname' tag"
        },
        {
            "field": "email",
            "error": "email must be a valid email address"
        }
    ],
    "success": false
}
```

Admittedly, the response message for the _name_ field isn't very nice and user friendly, but I didn't have time to fix that.

## Project Structure

A lot of the boilerplate code and the project structure is inspired by [ardanlabs](https://github.com/ardanlabs/service/). Another big inspiration for how I write my code is [Mat Ryer](https://github.com/matryer).

You find the main entrypoint to the application running the HTTP server at `cmd/quote-api/main.go`. You also have `cmd/quote-admin/main.go`, which allows you to perform migratations and seed the database. All the business logic is found under `internal/business/`. Code that is not related to the business logic, but also not meant to be shared, is found under `internal/foundation/`.

## Design Decisions (discussion)

I have made quite a few design decision that brings both pros and cons. 

1. How to test the application. I have made the decision to _unit_ test all handlers in detail and then provided some integration tests to make sure the endpoints don't fail when using the real quote service implementation. I'm not fully satisfied with this design decision and am instead contemplating whether it would've been better to only have integration tests for the endpoints.
2. When unit testing in general, is it preferable to be a bit more repetetive or to write a lot of helper functions? I tended to be a bit more repetitive so that you don't have to jump around a lot when reading the tests, but it also increases the number of lines of code. 
3. How to deal with the country code to region mapping? This logic can be found under `internal/business/region/`. I'm not satisfied with this solution. The purpose of the package seems odd and hard-coding every country code doesn't seem right.
4. The quote logic can be found under `internal/business/data/quote/`. There were a lot design decision made here. For example, is it this package's responsibility to calculate the shipment cost, or do I just specify an interface for calculating shipment cost? The same reasoning goes with the database calls. I ended up implementing both in the quote package. In the database case I tie the business logic with the choice of database (SQL), but in this app I thought the other alternative would increase the complexity of the code.
5. Regarding the PostgreSQL schema, I decided to keep everything in one flat table instead of normalizing into a _customers_ and _quotes_ table, for example. I think not normalizing makes sense for several reasons, but would love to head your input.
6. How to validate customer addresses. I made sure to specify that customer addresses are between 1 and 100 characters, but I was a bit confused about the specific format. You gave `Vasagatan 5B, Göteborg 41124` as an example, but can't addresses have very different formats depending on the country? I decided to not validate the specific format of an address.
7. Making all HTTP responses following a standard form with the `code`, `success` and `data` fields made the API arguably more user friendly, but introduced quite a lot of bloat in the unit and integration tests for the endpoints.

There are also design decisions about the project structure in general that are quite interesting to talk about. Would love to here your input!


## Licensing

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

Packages licensed under Apache License, Version 2.0:

* [conf](https://github.com/ardanlabs/conf)
* [darwin](https://github.com/dimiro1/darwin)
* [dockertest](https://github.com/ory/dockertest)