GO=go

run:
	$(GO) run *.go

build:
	$(GO) build *.go

clean:
	rm -- !(*.go)