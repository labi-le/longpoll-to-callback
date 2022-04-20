.DEFAULT: run

PROJ_NAME = lp-to-cb

MAIN_PATH = main.go
BUILD_PATH = build/package/

export CGO_ENABLED = 0

run:
	@go run $(MAIN_PATH)

build-prod:
	@go build -ldflags "-w" -a -v -o $(BUILD_PATH)$(PROJ_NAME) $(MAIN_PATH)

build-dev:
	@go build -v -o $(BUILD_PATH)$(PROJ_NAME) $(MAIN_PATH)
