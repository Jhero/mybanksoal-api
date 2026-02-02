@echo off
echo Generating Swagger 2.0 spec...
swagger generate spec -o ./docs/swagger.json --scan-models

echo Converting to OpenAPI 3.0...
go run cmd/spec-converter/main.go

echo Done! Documentation available at api.yml
