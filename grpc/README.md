
protoc --proto_path=sampleapp --go_out=sampleapp --go_opt=paths=source_relative sampleapp.proto


NEW_RELIC_API_KEY=<New Relic Key> go run server.go

NEW_RELIC_API_KEY=<New Relic Key> go run client.go

GRPC Server NewRelic Service-APM name: "newrelic-opentelemetry-poc Server"
GRPC Client NewRelic Service-APM name: "newrelic-opentelemetry-poc Client"
