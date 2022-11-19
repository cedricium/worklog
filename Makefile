.PHONY: all clean

all: clean worklog

worklog:
	@go build ./cmd/worklog

clean:
	@rm -f worklog
