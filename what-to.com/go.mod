module what-to.com

go 1.22.1

require (
    github.com/gorilla/mux v1.8.1
)

replace what-to.com/cmd/server => ./cmd/server

replace what-to.com/internal/router => ../internal/router

require github.com/gorilla/mux v1.8.1 // indirect
