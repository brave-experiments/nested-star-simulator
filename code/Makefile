.PHONY: all test lint clean

binary = nstarsim
godeps = *.go go.mod go.sum

all: test lint $(binary)

test:
	go test -cover ./...

lint:
	golangci-lint run

$(binary): $(godeps)
	go install

clean:
	rm -f $(binary)
