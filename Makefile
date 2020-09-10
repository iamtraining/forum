 
.PHONY: postgres adminer migrateup migratedown

postgres:
	docker run -d --name postgresql -e POSTGRES_DB=iamtraining -e POSTGRES_PASSWORD=1111 -p 5433:5433 postgres:12

adminer:
	docker run --rm -ti --network host adminer

migrateup:
	migrate -source file://migrations -database postgres://postgres:1111@localhost/iamtraining?sslmode=disable up

migratedown:
	migrate -source file://migrations -database postgres://postgres:1111@localhost/iamtraining?sslmode=disable down