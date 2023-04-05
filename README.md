# Helicopter Server

Helicopter Server is a project that provides access to a causal tree through
REST API and gRPC services. Users can add new nodes and get lists of nodes in
a subtree with a given root that has an lseq bigger than the given lseq. The
server uses a database with Lamport sequence numbers to manage the tree
structure.

## Table Of Contents
- [Helicopter Server](#helicopter-server)
  - [Table Of Contents](#table-of-contents)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [Alternative installation](#alternative-installation)
    - [Testing](#testing)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Running the server](#running-the-server)
  - [Examples](#examples)
  - [Contributing](#contributing)
  - [License](#license)

## Features

- AddNode: Add a new node to the causal tree with the specified parent and content.
- GetNodes: Retrieve a list of child nodes in a subtree under a specified root node.

## Prerequisites

- Go 1.19 compiler
- Protoc compiler
- Go and gRPC plugins for Protoc

## Installation

1. Clone this repository
   
```bash
git clone https://github.com/InfJoker/helicopter.git
cd helicopter
```

2. Install the `protoc` compiler, check out [this link](https://grpc.io/docs/protoc-installation/)
   for instructions on installing it for linux or MacOS

3. Install the `protoc` compiler plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

4. Update the `PATH`

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

5. Build the project

```bash
make build
```

The build binary will be in the `bin` folder.

### Alternative installation

If you have docker installed you can build a docker container for the helicopter
server using [this Dockerfile](build/server/Dockerfile):

```bash
docker build -t helicopter:latest -f build/server/Dockerfile .
```

### Testing

To run the tests for the project, use the following command:

```bash
make test
```

To generate a test coverage report, use the following command:

```bash
make coverage
```

## Configuration

The example server configuration is stored in the [configs/config.yml](configs/config.yml)
file. This file contains settings for the REST and gRPC servers, database address,
and the location of the OpenAPI spec template.

Customize the file as needed. If the config file provides the address of the goldb
server, the server will use it as storage. If you wish to use the non-persistent
in-memory storage, remove the `lseqdb` entry from the config file.


## Usage

### Running the server

To start the `helicopter` server run the binary providing the path to the config
file:

```bash
./bin/helicopter -config configs/config.yaml
```

This will launch the REST and gRPC servers using the configuration provided in
`configs/config.yaml`.

After starting the server, you can interact with it using any gRPC client to call the available RPC methods. For example,
you can call the `GetNodes` and `AddNode` methods to retrieve and add nodes in the causal tree. Refer to 
[web/openapi_template.yaml](web/openapi_template.yaml) and 
[proto/helicopter.proto](proto/helicopter.proto) for REST and gRPC specs,
respectively.

## Examples

Two example applications are provided in the examples directory. To build them
run the `make examples` command. The build binaries will be in the `bin` directory.

1. cli-messenger: A CLI messenger that interacts with the server. It has two
   modes: "write" for writing messages, and "listen" for listening to a specific
   thread.

![](https://github.com/InfJoker/helicopter/blob/main/assets/chat.gif)

2. chatgpt-bot: A chatbot that connects to the `chatGPT` thread and uses OpenAI
   API to interact with users. Users can prompt the chat bot by prefixing their
   message with `AskChat:`. We already run an instance of this bot, so you can 
   try interacting with it yourself by connecting to the `chatGPT` thread!

By manipulating contents of the data you can create different apps on top of our causal tree data model. Just use your imagination.

You can access our api through swagger: http://ds.sphericalpotatoinvacuum.xyz:8288/swagger/index.html
or use the grpc server at ds.sphericalpotatoinvacuum.xyz:1228.

## Contributing

To contribute to the Helicopter Server, please follow the usual GitHub workflow:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Implement your changes and commit them to your branch.
4. Push your changes to your fork.
5. Create a pull request against the main repository.

## License

This project is licensed under the [MIT License](LICENSE).
