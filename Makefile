APPNAME := gct

# 'all' is the default target
all: clean build setup-specs run-specs

clean:
	rm -f $(APPNAME)
	find . -type f -name "*coverprofile*" -delete
	find . -type f -name "*coverage*" -delete
	rm -rf vendor/
	rm -rf tmp/

# This make file only builds for linux OS. Modify GOOS here for local builds.
build:
	go build -o $(APPNAME)

run-linters:
	gofmt -w -s .

# should be used only once in the local machine
setup-specs:
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.3.0
	go install github.com/onsi/gomega/...

# for running the unit tests in the local machine as well as in the build machine
run-specs:
	export ACK_GINKGO_DEPRECATIONS=2.16.0 && \
		ginkgo -r --randomize-all --randomize-suites --cover --race --trace
	go vet

test-coverage:
	go build -o $(APPNAME)
	export ACK_GINKGO_DEPRECATIONS=2.16.0 && \
		ginkgo -r --randomize-all --randomize-suites --cover -coverprofile=coverage.out --race --trace
	go tool cover -html=coverage.out

run:
	./$(APPNAME)

.PHONY: clean all
