.PHONY: generate
generate:
	mkdir -p pkg/grpc_geo_provider
	protoc --proto_path vendor.protogen \
	--proto_path api \
	--go_out=pkg/grpc_geo_provider \
	--go-grpc_out=pkg/grpc_geo_provider \
	--validate_out lang=go:pkg/grpc_geo_provider \
	api/geo_provider.proto
		make move

.PHONY: move
move:
	mv pkg/grpc_geo_provider/github.com/Artenso/Geo-Service/geo/pkg/grpc_geo_provider/* pkg/grpc_geo_provider &&\
	rm -rf pkg/grpc_geo_provider/github.com/ 

.PHONY: vendor-proto
vendor-proto:
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google/protobuf ]; then \
			git clone https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
			mkdir -p  vendor.protogen/google/protobuf &&\
			mv vendor.protogen/protobuf/src/google/protobuf/*.proto vendor.protogen/google/protobuf &&\
			rm -rf vendor.protogen/protobuf ;\
		fi