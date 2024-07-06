package rpc

//go:generate find . -type f -name *.proto -exec protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. {} ;
