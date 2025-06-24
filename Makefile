PROTO_COMPOSE_FILE := ./development/docker-compose.proto.yml

gen-proto:
	docker compose -f ${PROTO_COMPOSE_FILE} up generate_pb_go --build
