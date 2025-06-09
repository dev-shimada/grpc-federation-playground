# gRPC message sample
## Project Structure

- `proto/message/`: Contains the protobuf definitions and generated code.
- `buf.yaml`: Configuration for Buf, a tool for managing Protobuf files.
- `buf.gen.yaml`: Configuration for generating code from Protobuf files.

## Setup

### Prerequisites

- Docker

### Generating Protobuf Code

To generate the protobuf code using Buf, run:
```bash
buf generate
```

## grpcurl
```console
# grpcurl -plaintext localhost:8081 list
/message.v1.MessageService/Get
/message.v1.MessageService/PingPong
/message.v1.MessageService/Post
```
```console
# grpcurl -plaintext localhost:8081 list message.v1.Message
message.Message.PingPong
message.Message.Subscribe
```
```console
# grpcurl -plaintext -d '{"user": "my self"}' localhost:8081 message.v1.Message.Subscribe
{
  "from": "my self",
  "message": "0: Hello, my self!"
}
{
  "from": "my self",
  "message": "1: Hello, my self!"
}
{
  "from": "my self",
  "message": "2: Hello, my self!"
}
```
```console
# grpcurl -plaintext -d '{"user": "my self", "message": "hello!"}' localhost:8081 message.v1.Message.PingPong
{
  "from": "my self",
  "message": "hello!"
}
```
```console
# grpcurl -plaintext localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"service": "/message.v1.Message/PingPong"}' localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"service": "/message.v1.Message/Subscribe"}' localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
