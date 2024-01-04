controller:
	cd src && go run ./cmd/controller/api-server && cd ..
let:
	docker compose up -d --build
kill:
	docker kill cuternetes-worker0-1 cuternetes-worker1-1 cuternetes-worker2-1