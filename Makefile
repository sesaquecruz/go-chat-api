#!/bin/bash

update_dependency_injection:
	@if !command -v wire >/dev/null 2>&1 ; then \
		echo "Go Wire is not installed. Installing..."; \
		go install github.com/google/wire/cmd/wire@latest; \
	fi

	@echo "Updating dependency injection";
	@wire di/wire.go;

update_mocks:
	@if !command -v mockery >/dev/null 2>&1 ; then \
		echo "Mockery is not installed. Installing..."; \
		go install github.com/vektra/mockery/v2@v2.33.2; \
	fi

	@echo "Updating mocks";
	@mockery;
