version_folder = $(subst /,_,$(lastword $(shell go version)))

all: check-env
	@echo "GOATS - GO Attribute-based Template System"
	@echo ""
	@echo "make:\n    show this help message"
	@echo "make binary:\n    build and install goats executable to $(GOATS_GOPATH)/bin"
	@echo "make clean:\n    clean up the source folder"
	@echo "make fmt:\n    format goats source code"
	@echo "make unittests:\n    run unit tests for goats"
	@echo "make amd64:\n    build Ubuntu deb package for amd64"
	@echo "make i386:\n    build Ubuntu deb package for i386"
	@echo "make armhf:\n    build Ubuntu deb package for armhf"

clean: check-env
	rm $(GOATS_GOPATH)/../pkg/$(version_folder)/goats-html/* -rf
	rm $(GOATS_GOPATH)/../bin/goats -f
	rm -f debian/usr/share/goats-html/bin/*
	rm -rf debian/usr/share/goats-html/go

binary: check-env
	make clean
	go get code.google.com/p/go.net/html
	go get github.com/howeyc/fsnotify
	go install goats-html/cmd/goats
	rm $(GOATS_GOPATH)/../pkg/$(version_folder)/goats-html/* -rf

fmt: check-env
	cd cmd/goats && go fmt .
	cd examples && go fmt .
	cd examples/data && go fmt .
	cd examples/server && go fmt .
	cd goats && go fmt .
	cd goats/runtime && go fmt .
	cd tests && go fmt .
	cd tests/data && go fmt .

unittests: check-env
	../../bin/goats gen --package_root .. --template_dir goats-html/tests/templates/
	go test goats-html/goats
	go test goats-html/tests

deb-package: check-env
	make binary
	mkdir -p debian/usr/share/goats-html/bin/
	cp ../../bin/goats debian/usr/share/goats-html/bin/
	mkdir -p debian/usr/bin
	ln -f -s /usr/share/goats-html/bin/goats debian/usr/bin/goats
	mkdir -p debian/usr/share/goats-html/go/src/goats-html/goats/runtime/
	cp -r goats/runtime debian/usr/share/goats-html/go/src/goats-html/goats/
	rm debian/usr/share/goats-html/go/src/goats-html/goats/runtime/filters_test.go
	chmod -R 0755 debian/usr/share/goats-html
	cd debian && fakeroot dpkg -b . ..

	@echo
	@echo "###############################################"
	@echo
	@echo "After installing the deb file, you need to add /usr/share/goats-html/go in GOPATH. Example:"
	@echo "    export GOPATH=~/go:/usr/share/goats-html/go"
	@echo
	@echo "###############################################"
	@echo

amd64:
	cp debian/DEBIAN/control-amd64 debian/DEBIAN/control
	make deb-package

i386:
	cp debian/DEBIAN/control-i386 debian/DEBIAN/control
	make deb-package

armhf:
	cp debian/DEBIAN/control-armhf debian/DEBIAN/control
	make deb-package

check-env:
ifndef GOATS_GOPATH
	@echo "Before you build GOATS, you need to set the GOATS_GOPATH environment variable like:"
	@echo "    export GOATS_GOPATH=~/go/src"
	@echo ""
	$(error GOATS_GOPATH is undefined)
endif

