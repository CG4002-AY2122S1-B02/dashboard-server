# dashboard-server

### Start Postgres Instance
Create a Database called 'cg4002' with a user that can access that database called 'g2' 

### Start Postgres database
``` sudo service postgresql start ```

### Run Ultra96 test
```` python3 scripts/laptop_server_ultra96.py ````

### Run Dashboard Server
```` ./bin/dashboard_server_linux````

#### Compile for MACOS
```GOOS=darwin GOARCH=amd64 go build -o bin/dashboard_server_macos cmd/api/run.go```