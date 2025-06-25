# protoc-gen-http-client

A Protocol Buffers compiler plugin that generates HTTP client code for gRPC services. This plugin automatically creates HTTP clients from your `.proto` files, making it easy to call gRPC services over HTTP without manually writing client code.

## Features

- **Automatic HTTP Client Generation**: Generates Go HTTP clients from gRPC service definitions
- **gRPC-Gateway Integration**: Works seamlessly with gRPC-Gateway HTTP annotations
- **Type-Safe**: Full type safety with Protocol Buffers message types
- **Error Handling**: Proper gRPC status code mapping and error handling
- **JSON Encoding/Decoding**: Automatic JSON serialization and deserialization
- **Path Parameter Injection**: Support for path parameters in HTTP URLs

## Installation

### Prerequisites

- Go 1.24.4 or later
- Protocol Buffers compiler (`protoc`)
- gRPC-Gateway annotations

### Install the Plugin

```bash
go install github.com/Open-Myth/protoc-gen-http-client@latest
```

### Dependencies

The plugin requires the following Go modules:
- `github.com/dave/jennifer` - Code generation
- `github.com/grpc-ecosystem/grpc-gateway/v2` - gRPC-Gateway support
- `github.com/json-iterator/go` - JSON processing
- `google.golang.org/grpc` - gRPC support
- `google.golang.org/protobuf` - Protocol Buffers

## Usage

### 1. Define Your Service

Create a `.proto` file with gRPC service definitions and HTTP annotations:

```protobuf
syntax = "proto3";
option go_package = "/greeting";
package chat;

import "google/api/annotations.proto";

service GreetingService {
  rpc SayHello(SayHelloRequest)
      returns (SayHelloResponse) {
    option (google.api.http) = {
      post : "/v1/hello",
      body: "*"
    };
  };
}

message SayHelloRequest {
  string ping = 1;
}

message SayHelloResponse { 
  string pong = 1; 
}
```

### 2. Generate HTTP Client Code

Use the `protoc` compiler with the `--http-client_out` flag:

```bash
protoc \
  proto/*.proto \
  -I=/usr/local/include \
  --proto_path=proto \
  --go_out=:pb \
  --go-grpc_out=:pb \
  --grpc-gateway_out=:pb \
  --http-client_out=:pb
```

### 3. Use the Generated Client

The plugin generates a file named `{service}_http-client.pb.go` containing the HTTP client:

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "time"

    api "github.com/Open-Myth/protoc-gen-http-client/pb/greeting"
)

func main() {
    // Create a new HTTP client
    client := api.NewGreetingServiceHTTPClient("http://localhost:8080")
    
    // Make HTTP calls
    resp, err := client.SayHello(context.Background(), &api.SayHelloRequest{
        Ping: fmt.Sprintf("Hello %d", rand.Int31n(100)),
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp)
}
```

## Generated Code Structure

The plugin generates the following structure for each service:

```go
// HTTPClient is a http client for the GreetingService service
type GreetingServiceHTTPClient struct {
    BaseURL string
}

func NewGreetingServiceHTTPClient(baseURL string) *GreetingServiceHTTPClient {
    return &GreetingServiceHTTPClient{BaseURL: baseURL}
}

// SayHello is a http call method for the GreetingService service
func (c *GreetingServiceHTTPClient) SayHello(ctx context.Context, reqData *SayHelloRequest) (*SayHelloResponse, error) {
    // Implementation with proper error handling and JSON encoding/decoding
}
```

## Docker Support

For easier development, you can use the provided Docker setup:

```bash
# Generate code using Docker
make gen-proto
```

This uses the Docker Compose configuration in `development/docker-compose.proto.yml` to run the protoc compiler with all necessary plugins.

## Examples

See the `examples/` directory for complete working examples:

- `examples/client/` - HTTP client usage example
- `examples/server/` - gRPC server with HTTP gateway

## Development

### Building from Source

```bash
git clone https://github.com/Open-Myth/protoc-gen-http-client.git
cd protoc-gen-http-client
go build -o protoc-gen-http-client
```

### Running Tests

```bash
go test ./...
```

### Project Structure

```
.
├── main.go                 # Plugin entry point
├── generate.go             # Main generation logic
├── generate_test.go        # Tests for generation
├── generator/              # Code generation utilities
├── util/                   # HTTP client utilities
├── proto/                  # Example .proto files
├── pb/                     # Generated Go code
├── examples/               # Usage examples
└── development/            # Development tools and Docker setup
```

## How It Works

1. **Service Analysis**: The plugin analyzes gRPC service definitions and extracts HTTP annotations
2. **Code Generation**: Uses the `jennifer` library to generate Go code
3. **HTTP Client Creation**: Creates HTTP clients with proper error handling and JSON encoding/decoding
4. **Path Parameter Support**: Injects path parameters from request messages into HTTP URLs

## Error Handling

The generated clients handle errors by:
- Converting HTTP status codes to gRPC status codes
- Proper JSON error response parsing
- Maintaining gRPC error semantics

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built on top of Protocol Buffers and gRPC
- Uses gRPC-Gateway for HTTP annotations
- Inspired by the need for simple HTTP clients for gRPC services 