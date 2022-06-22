.PHONY: clean build

BUILD_PATH = ./build
BUILD_NAME = metagraph

clean:
	rm -rf $(BUILD_PATH)

build:
	go build -o $(BUILD_PATH)/$(BUILD_NAME)

run:
	$(BUILD_PATH)/$(BUILD_NAME)
