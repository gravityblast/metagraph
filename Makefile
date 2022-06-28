.PHONY: clean build

BUILD_PATH = ./build
BUILD_NAME = metagraph

build: clean
	go build -o $(BUILD_PATH)/$(BUILD_NAME)

clean:
	rm -rf $(BUILD_PATH)


run:
	$(BUILD_PATH)/$(BUILD_NAME)
