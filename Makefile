GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get

BINARY_NAME=thub

all: clean build

build:
	$(GO_BUILD) -o $(BINARY_NAME) -v -ldflags="-X main.MenuIcon=command.square.fill"

package: clean build
	rm -rf tHub.app/Contents/MacOS/${BINARY_NAME}
	mv $(BINARY_NAME) tHub.app/Contents/MacOS
	@echo "Done"

clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)

run: clean
	$(GO_BUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

test:
	$(GO_TEST) -v ./...

deps:
	$(GO_GET)