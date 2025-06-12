# gRPC user sample
## Project Structure

- `proto/user/`: Contains the protobuf definitions and generated code.
- `buf.yaml`: Configuration for Buf, a tool for managing Protobuf files.
- `buf.gen.yaml`: Configuration for generating code from Protobuf files.

## Setup

### Prerequisites

- Docker

### Generating Ent Code
To generate the Ent code, run:
```bash
go generate ent/generate.go
```

### Generating Protobuf Code

To generate the protobuf code using Buf, run:
```bash
buf generate
```

## grpcurl
```console
# grpcurl -plaintext localhost:8082 list
user.v1.UserService
```
```console
# grpcurl -plaintext localhost:8082 list user.v1.UserService
user.v1.UserService.Get
user.v1.UserService.PingPong
user.v1.UserService.Post
```
```console
# grpcurl -plaintext -d '{"service": "user.v1.UserService"}' localhost:8082 grpc.health.v1.Health.Check
{
  "status": "SERVING"
}
```
```console
# grpcurl -plaintext -d '{"email": "test@example.com", "name": "test_name"}' localhost:8082 user.v1.UserService.PingPong
{
  "userId": "test_id",
  "text": "hello!"
}
```
```console
# grpcurl -plaintext -d '{"email": "test@example.com", "name": "test_name"}' localhost:8082 user.v1.UserService.Post
{
  "id": "019764ad-b9c9-75b9-9663-716d625cc350",
  "email": "test@example.com",
  "name": "test_name"
}
```
```console
# grpcurl -plaintext -d '{"id": "019764ad-b9c9-75b9-9663-716d625cc350"}' localhost:8082 user.v1.UserService.Get
{
  "id": "019764ad-b9c9-75b9-9663-716d625cc350",
  "email": "test@example.com",
  "name": "test_name",
  "createdAt": "2025-06-12T15:06:36Z",
  "updatedAt": "2025-06-12T15:06:36Z"
}
```

## curl
```console
# curl --json '{"service": "user.v1.UserService"}' localhost:8082/grpc.health.v1.Health/Check
{"status":"SERVING_STATUS_SERVING"}
```
```console
# curl --json '{"email": "test@example.com", "name": "test_name"}' localhost:8082/user.v1.UserService/PingPong
{"email":"test@example.com", "name":"test_name"}
```
```console
# curl --json '{"email": "test@example.com", "name": "test_name"}' localhost:8082/user.v1.UserService/Post
{"id":"019764b1-4dc8-7114-928f-6526676648d4", "email":"test@example.com", "name":"test_name"}
```
```console
# curl --json '{"id": "019764b1-4dc8-7114-928f-6526676648d4"}' localhost:8082/user.v1.UserService/Get
{"id":"019764b1-4dc8-7114-928f-6526676648d4", "email":"test@example.com", "name":"test_name", "createdAt":"2025-06-12T15:10:30Z", "updatedAt":"2025-06-12T15:10:30Z"}
```
