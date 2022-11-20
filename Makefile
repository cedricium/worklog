.PHONY: all clean install

all: clean worklog

worklog:
	@go build ./cmd/worklog

clean:
	@rm -f worklog

install:
	@go build -o ${GOPATH}/bin/worklog ./cmd/worklog