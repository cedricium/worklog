.PHONY: all clean

all: clean worklog

worklog:
	@go build -o ${GOPATH}/bin/worklog ./cmd/worklog

clean:
	@rm -f ${GOPATH}/bin/worklog