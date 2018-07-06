

export GOPATH=$(shell pwd):$(shell pwd)/vendor

all: vendor/src/github.com/miekg/pkcs11
	cd src && go build -o ../test-softhsm-v2 && cd -

vendor/src/github.com/miekg/pkcs11:
	GOPATH=$(shell pwd)/vendor go get github.com/miekg/pkcs11

clean:
	rm -fr test-softhsm-v2