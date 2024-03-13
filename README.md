```bash
# start database
cd database
bash ./initdb.bash
cd ..

# start backend
cd backend
# go version: 1.21.6
go run ./  # start backend server
# keep backend server running and start a new terminal window for further commands


# start frontend
cd frontend
# bun version: 1.0.14
bun install  # download required node_modules
bun run dev  # start frontend for development
# site is now live at http://localhost:3000
```