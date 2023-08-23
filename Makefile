#!/bin/bash

update_dependency_injection:
	@if !command -v wire >/dev/null 2>&1 ; then \
		echo "Go Wire is not installed. Installing..."; \
		go install github.com/google/wire/cmd/wire@latest; \
	fi

	@echo "Updating dependency injection";
	@wire di/wire.go;

update_mocks:
	@if !command -v mockgen >/dev/null 2>&1 ; then \
		echo "Go Mock is not installed. Installing..."; \
		go install go.uber.org/mock/mockgen@latest; \
	fi

	@echo "Updating mocks";
	mockgen --source=internal/domain/repository/room_repository.go --destination=test/mock/room_repository.go --package=mock;
	mockgen --source=internal/domain/repository/message_repository.go --destination=test/mock/message_repository.go --package=mock;
	mockgen --source=internal/domain/gateway/message_sender_gateway.go --destination=test/mock/message_sender_gateway.go --package=mock;
	mockgen --source=internal/domain/gateway/message_receiver_gateway.go --destination=test/mock/message_receiver_gateway.go --package=mock;

