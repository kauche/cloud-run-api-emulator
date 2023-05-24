BUILD_TOOLS := cd ./tools && go build -o

BIN_DIR := ./bin

XO := $(abspath $(BIN_DIR)/xo)

.PHONY: xo
xo:
	@$(BUILD_TOOLS) $(XO) github.com/xo/xo
