module what-to.com

go 1.22

require github.com/gorilla/mux v1.8.1

replace what-to.com/cmd/server => ./cmd/server

replace what-to.com/internal/router => ../internal/router

replace what-to.com/internal/entity => ../internal/entity

replace what-to.com/internal/logger => ../internal/logger

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
