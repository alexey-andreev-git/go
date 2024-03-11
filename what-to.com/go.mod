module what-to.com

go 1.22

require github.com/gorilla/mux v1.8.1

replace what-to.com/cmd/server => ./cmd/server

replace what-to.com/internal/router => ../internal/router

replace what-to.com/internal/entity => ../internal/entity

replace what-to.com/internal/logger => ../internal/logger

require (
	github.com/lib/pq v1.10.9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
