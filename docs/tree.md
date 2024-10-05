.
|--- api
|   |--- openapi
|   |   |--- README.md
|   |   |--- user.yaml
|   |--- proto
|   |   |--- user
|   |   |   |--- v1
|   |   |   |   |--- user_service.proto
|   |   |   |   |--- user_type.proto
|   |   |--- buf.gen.yaml
|   |   |--- buf.yaml
|   |   |--- README.md
|--- cmd
|   |--- main.go
|   |--- README.md
|--- docs
|   |--- README.md
|   |--- tree.md
|--- internal
|   |--- adapter
|   |   |--- driven
|   |   |   |--- db
|   |   |   |   |--- migration
|   |   |   |   |   |--- 000001_create_users_table.down.sql
|   |   |   |   |   |--- 000001_create_users_table.up.sql
|   |   |   |   |   |--- init.sql
|   |   |   |   |   |--- README.md
|   |   |   |   |--- repository
|   |   |   |   |   |--- README.md
|   |   |   |   |   |--- user.go
|   |   |   |   |   |--- user_mock.go
|   |   |   |   |--- db_handler.go
|   |   |   |   |--- postgres_transaction.go
|   |   |   |   |--- postgres_transaction_mock.go
|   |   |   |--- passowrd.go
|   |   |--- driver
|   |   |   |--- grpc
|   |   |   |   |--- proto
|   |   |   |   |   |--- user
|   |   |   |   |   |   |--- v1
|   |   |   |   |   |   |   |--- user_service.pb.go
|   |   |   |   |   |   |   |--- user_service_grpc.pb.go
|   |   |   |   |   |   |   |--- user_type.pb.go
|   |   |   |   |--- README.md
|   |   |   |   |--- user_service.go
|   |   |   |--- http
|   |   |   |   |--- README.md
|   |   |   |--- model
|   |   |   |   |--- README.md
|   |   |   |--- service
|   |   |   |   |--- README.md
|   |   |   |   |--- user_service.go
|   |   |   |   |--- user_service_test.go
|   |   |   |--- README.md
|   |--- domain
|   |   |--- error
|   |   |   |--- user.go
|   |   |--- model
|   |   |   |--- user.go
|   |   |--- README.md
|   |--- port
|   |   |--- driven
|   |   |   |--- db
|   |   |   |   |--- repository
|   |   |   |   |   |--- README.md
|   |   |   |   |   |--- user.go
|   |   |   |   |--- db_handler.go
|   |   |   |   |--- db_transaction.go
|   |   |   |--- passowrd.go
|   |   |--- driver
|   |   |   |--- model
|   |   |   |   |--- README.md
|   |   |   |   |--- user.go
|   |   |   |--- service
|   |   |   |   |--- README.md
|   |   |   |   |--- user.go
|--- pkg
|   |--- README.md
|--- script
|   |--- create_tree.ps1
|   |--- fill_tree.ps1
|--- test
|   |--- README.md
|--- .env
|--- .env.example
|--- .gitignore
|--- docker-compose.yml
|--- Dockerfile
|--- go.mod
|--- go.sum
|--- Makefile
|--- README.md
