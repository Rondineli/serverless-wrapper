TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=example.com
NAMESPACE=svl
NAME=wrapper
BINARY=terraform-provider-${NAME}
VERSION=0.3
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	#go test -i $(TEST) || exit 1                                                   
	#echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    
	go test -v -cover ./internal/provider/

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   
