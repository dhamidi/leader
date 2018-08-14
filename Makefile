.PHONY: install
install: pack
	go install .

.PHONY: pack
pack: $(GOPATH)/bin/packr
	packr

$(GOPATH)/bin/packr:
	go get -u github.com/gobuffalo/packr/...
