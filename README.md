# Helicopter Server

Helicopter Server is a project for a causal tree access management system. The server handles the storage and retrieval of nodes in a causal tree, allowing you to easily manage and manipulate the data. In addition to gRPC, Helicopter Server also provides a RESTful API for interacting with the server. The API documentation is provided through Swagger.

## Features

- AddNode: Add a new node to the causal tree with the specified parent and content.
- GetNodes: Retrieve a list of child nodes in a subtree under a specified root node.

## Prerequisites

- Go 1.19 compiler
- Protoc compiler
- Go and gRPC plugins for Protoc

## Try with examples!

We've built cli-messenger as an example of using our server. Also we run ChatGPT chatbot, which listens on chatGPT thread. You can try asking it some prompts yourself!

![](https://github.com/InfJoker/helicopter/blob/main/assets/chat.gif)

To build them just run:

```bash
make examples
```

By manipulating contents of the data you can create different apps on top of our causal tree data model. Just use your imagination.

Also you can access our api through swagger: http://ds.sphericalpotatoinvacuum.xyz:8288/swagger/index.html. Or use grpc on ds.sphericalpotatoinvacuum.xyz:1228, refer to [proto/helicopter.proto](https://github.com/InfJoker/helicopter/blob/main/proto/helicopter.proto).

## Getting Started

To get started with the Helicopter Server, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/yourusername/helicopter-server.git
cd helicopter-server
```

2. Install the required dependencies:

```bash
go mod download
```

3. Build the server binary:

```bash
make build
```

4. Start the server:

```bash
./bin/helicopter
```


## Usage

After starting the server, you can interact with it using any gRPC client to call the available RPC methods. For example, you can call the `GetNodes` and `AddNode` methods to retrieve and add nodes in the causal tree.

## Testing

To run the tests for the project, use the following command:

```bash
make test
```

To generate a test coverage report, use the following command:

```bash
make coverage
```

## Contributing

To contribute to the Helicopter Server, please follow the usual GitHub workflow:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Implement your changes and commit them to your branch.
4. Push your changes to your fork.
5. Create a pull request against the main repository.

## License

This project is licensed under the [MIT License](LICENSE).
