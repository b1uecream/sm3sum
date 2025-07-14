BINARY = sm3sum
SRC = sm3sum.go
VERSION := $(shell grep -Po 'version = \"\K[^\"]+' $(SRC))

PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: all build clean install uninstall test crossbuild distclean dist

all: build

build:
	go build -o $(BINARY) $(SRC)

install: build
	install -m 755 $(BINARY) /usr/local/bin/$(BINARY)

uninstall:
	rm -f /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)
	rm -f *.sm3 *.txt *.1
	go mod tidy

test:
	echo "hello world" > test.txt
	./$(BINARY) test.txt > test.sm3
	./$(BINARY) -c test.sm3
	./$(BINARY) -c test.sm3 --status
	@echo "Corrupt test (expect failure):"
	@sed 's/./0/' test.sm3 > bad.sm3 && ! ./$(BINARY) -c bad.sm3 --status
	./$(BINARY) --tag test.txt
	./$(BINARY) --zero test.txt | hexdump -C

version:
	@echo "Version: $(VERSION)"

man:
	gzip -c sm3sum.1 > sm3sum.1.gz

completion:
	@echo "Installing bash completion..."
	install -Dm644 sm3sum_completion /etc/bash_completion.d/sm3sum
	@echo "Installing zsh completion..."
	install -Dm644 _sm3sum /usr/share/zsh/site-functions/_sm3sum

crossbuild:
	@echo "Starting cross compilation for platforms: $(PLATFORMS)"
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		OUTPUT=$(BINARY)-$${OS}-$${ARCH}; \
		if [ "$${OS}" = "windows" ]; then OUTPUT=$${OUTPUT}.exe; fi; \
		echo "Building $$OUTPUT ..."; \
		GOOS=$${OS} GOARCH=$${ARCH} go build -o $$OUTPUT $(SRC); \
	done
	@echo "Cross compilation done."

distclean: clean
	rm -rf dist

dist: distclean crossbuild
	mkdir -p dist
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		BINNAME=$(BINARY)-$${OS}-$${ARCH}; \
		if [ "$${OS}" = "windows" ]; then BINNAME=$${BINNAME}.exe; fi; \
		PKGNAME=$(BINARY)-$${OS}-$${ARCH}-v$(VERSION); \
		mkdir -p dist/$$PKGNAME; \
		cp $$BINNAME dist/$$PKGNAME/; \
		echo "$$PKGNAME" > dist/$$PKGNAME/README.txt; \
		tar -czf dist/$$PKGNAME.tar.gz -C dist $$PKGNAME; \
		rm -rf dist/$$PKGNAME; \
		echo "Created dist/$$PKGNAME.tar.gz"; \
	done
