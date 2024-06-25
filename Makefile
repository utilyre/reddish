.PHONY: protogen

protogen:
	protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. rpc/**/*.proto
