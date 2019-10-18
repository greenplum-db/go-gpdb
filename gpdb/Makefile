.PHONY: docker
docker:
	env GOOS=linux GOARCH=amd64 go build -o gpdb
	docker build -t pivotaldata/greenplum-db-6.0.1  .