postgres:
	docker compose up -d
createdb:
	docker exec -it $$(docker ps -f name=postgres -q) createdb --username=root --owner=root simplebank
dropdb:
	docker exec -it $$(docker ps -f name=postgres -q) dropdb --username=root simplebank
.PHONY:postgres createdb dropdb
