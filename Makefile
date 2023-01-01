protoc: dependencies
	@$(MAKE) --no-print-directory -C api protoc

dependencies: tools-install

tools-install:
	@$(MAKE) --no-print-directory -C api tools-install
