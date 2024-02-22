
postgres-up:
	docker run --name postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=quotation_db -p 5433:5432 -d postgres:latest


