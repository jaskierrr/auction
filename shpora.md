protoc -I=api/proto --go_out=. --go-grpc_out=. api/proto/auction.proto

protoc -I. --go_out=. --go-grpc_out=. --grpc-gateway_out . --openapiv2_out . api/proto/auction.proto

docker-compose up --build -d

docker-compose up -d

docker-compose ps

docker exec -it card-project sh
