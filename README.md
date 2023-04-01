# ECE49595 NPM and Github package web-service Project
## Docker Compose Up
- All components have necessary dockerfiles
- To run full docker `docker compose up`
- Web-service will run on port 4200.

Use jar file to create backend gin api code
`java -jar openapi-generator-cli-6.3.0.jar generate -i backend\api\openapi.yaml -g go-gin-server -o backend --package-name openapi --additional-properties=apiPath=openapi`