# User Flight Tracking Microservice

Simple microservice API that can help understand and track how the flight path of a particular person can be queried. The API must accept a request that includes a list of flights, which are defined by an origin and destination airport code.

## Installation and Setup

1. Clone the repository:
```
git clone https://github.com/svillad/user-flight-tracking.git
cd user-flight-tracking
```

2. Build the microservice:
```
make build
```

3. Run the microservice:
```
make run
```

The microservice will be available at `http://localhost:8080`.

## Endpoints

Endpoint to retrieve the flight path information.

- **URL:** `/calculate`
- **Method:** `POST`
- **Content-Type:** `application/json`

#### Request Body

Example:

```
{
  "flights": [
    ["IND", "EWR"],
    ["SFO", "ATL"],
    ["GSO", "IND"],
    ["ATL", "GSO"]
  ]
}
```

#### Response Body

Example:

```
{
  "start": "SFO",
  "end": "EWR",
  "path": [
    "SFO",
    "ATL",
    "GSO",
    "IND",
    "EWR"
  ]
}
```

#### Response Codes

- `200 OK`: Successful response with the flight path information.
- `400 Bad Request`: Invalid request body or missing required fields.
- `404 Not Found`: Flight path not found or invalid airports.
- `405 Method Not Allowed`: when you use an invalid method in the mirocservice

## Directory Structure

- `router/`: Contains the router configuration using `github.com/gorilla/mux`.
- `controllers/`: Handles the HTTP requests and responses.
- `translators/`: Converts data between different formats or structures.
- `mediators/`: Implements the business logic and coordinates between different components.
- `dto/`: Data transfer objects used for communication between components.
- `gateways/`: Handles external service interactions.
- `models/`: Defines the data models used in the microservice.

## Testing

The microservice includes unit tests for the following components:

- `controllers/`
- `mediators/`
- `gateways/`
- `translators/`

To run the tests, use the following command:
```
make test
```

## Generating Mocks

Mocks for the different components can be generated automatically using the `go generate` command. The mock implementation details are specified in the `gen.go` file.

To generate mocks, use the following command:
```
make generate
```


## Makefile Commands

- `make build`: Build the microservice.
- `make run`: Run the microservice.
- `make lint`: Run the linter for code linting.
- `make test`: Run the unit tests.
- `make cover`: Generate test coverage report.
- `make generate`: Generate mocks.
- `make regenerate`: Clean existing mocks and regenerate them.
- `make clean-mock`: Clean generated mock files.
- `make vendor`: Create a vendor directory for dependencies.
- `make remove-vendor`: Remove the vendor directory.


For more information on available Makefile commands, refer to the Makefile in the project root.

## Docker

### Building the Docker Image

To build the Docker image for the microservice, run the following command:
```
docker build -t user-flight-tracking:1.0 .
```

This command will build the Docker image with the tag user-flight-tracking:1.0 using the Dockerfile in the current directory.

### Running the Docker Container

To run the microservice as a Docker container, use the following command:
```
docker run -p 8080:8080 user-flight-tracking:1.0
```

This command will start a Docker container based on the user-flight-tracking:1.0 image and map port 8080 of the container to port 8080 of the host machine. You can access the microservice at http://localhost:8080.

### Running the Docker Container in Interactive Mode

If you need to access the Docker container in interactive mode for debugging or running commands inside the container, you can use the following command:
```
docker run -p 8080:8080 -it user-flight-tracking:1.0 /bin/bash
```

This command will start the Docker container in interactive mode with a Bash shell. You can then execute commands or perform any necessary debugging tasks within the container.

Please note that these commands assume you have Docker installed and properly configured on your system.

## Contributing

Contributions are welcome!.
