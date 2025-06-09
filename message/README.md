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
# grpcurl -plaintext localhost:8081 list message.v1.MessageService
message.v1.MessageService.Get
message.v1.MessageService.PingPong
message.v1.MessageService.Post
```
```console
# grpcurl -plaintext -d '{"user_id": "test_id", "text": "hello!"}' localhost:8081 message.v1.MessageService.PingPong
{
  "userId": "test_id",
  "text": "hello!"
}
```
```console
# grpcurl -plaintext localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"service": "/message.v1.MessageService/PingPong"}' localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"service": "/message.v1.MessageService/Get"}' localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"service": "/message.v1.MessageService/Post"}' localhost:8081 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
