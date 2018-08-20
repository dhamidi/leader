.PHONY: build
build:
	gox -os="linux darwin" -arch="amd64" -output="leader.{{.OS}}.{{.Arch}}" -ldflags "-X main.Revision=`git rev-parse --short HEAD` -X main.Release=`git tag --points-at HEAD | head -1`" -verbose ./...


.PHONY: install
install: pack
	go install -ldflags "-X main.Revision=`git rev-parse --short HEAD`"

.PHONY: pack
pack: $(GOPATH)/bin/packr
	packr

$(GOPATH)/bin/packr:
	go get -u github.com/gobuffalo/packr/...
