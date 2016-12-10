ifndef UNIT
	UNIT := `juju status --format json | jq -r '.applications["chatchat"].units|keys|.[]' | head -1`
endif
ifndef CLIENT
	CLIENT := $(shell hostname)
endif

.PHONY: all
all: build

.PHONY: build
build:
	go install github.com/mattyw/chatchat
	cp $(GOBIN)/chatchat files/chatchat
	charm build -l debug

.PHONY: clean
clean:
	$(RM) -r builds deps trusty

.PHONY: deploy
deploy: build
	juju deploy $(shell pwd)/trusty/chatchat chatchat --series trusty

.PHONY: upgrade
upgrade: build
	juju upgrade-charm --path $(shell pwd)/trusty/chatchat chatchat --force-units
	juju resolved $(UNIT)
