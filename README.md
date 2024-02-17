```bash
# start database
cd database
docker compose -f docker-compose.yml up -d

# start backend
cd backend
go get ./  # download required packages
go run ./  # start backend server

# start frontend
cd frontend
bun install  # download required node_modules
bun run dev  # start frontend for development
```