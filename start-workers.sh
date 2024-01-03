docker compose up --build -d

sleep 10

docker exec cuternetes-worker0-1 sh -c "air -c .air.toml"
docker exec cuternetes-worker1-1 sh -c "air -c .air.toml"
docker exec cuternetes-worker2-1 sh -c "air -c .air.toml"

