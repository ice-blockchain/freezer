OZ_VERSION := 5.0.1
OZ_URL  := https://github.com/OpenZeppelin/openzeppelin-contracts/archive/refs/tags/v$(OZ_VERSION).tar.gz
ICE_URL := https://codeload.github.com/ice-blockchain/erc-20-ice-coins-distribution/zip/refs/heads/master
ICE_CONTRACT := output/erc-20-ice-coins-distribution-master/contracts/ICEToken.sol

# Compile a contract
# $(1) -- output type. One of: bin, abi.
define compile_contract
	@set -e; \
	mkdir -p output; \
	if [ -e $(ICE_CONTRACT) ]; then \
		echo "---> Using $(ICE_CONTRACT) for $(1) generation"; \
		solc @openzeppelin=./output/openzeppelin-contracts-$(OZ_VERSION) --overwrite --$(1) $(ICE_CONTRACT) -o output/$(1); \
	elif [ -e ICEToken.flatten.sol ]; then \
		echo "---> Using local ICEToken.flatten.sol for $(1) generation"; \
		solc --overwrite --$(1) ICEToken.flatten.sol -o output/$(1); \
	else \
		echo "---> No local contract was found"; \
		false; \
	fi
endef

.PHONY: all
all: generate

.PHONY: tools
tools:
	@solc --version && abigen --version

.PHONY: generate
generate: ice_token.go

.PHONY: bindata
bindata: tools
	@$(call compile_contract,bin)

.PHONY: abidata
abidata: tools
	@$(call compile_contract,abi)

ice_token.go: tools bindata abidata
	abigen --bin=output/bin/ICEToken.bin --abi=output/abi/ICEToken.abi --pkg=coindistribution --out=$@

download:
	mkdir -p output
	wget -O- $(OZ_URL) | tar -zxvf- -C output
	wget -O output/ice.zip $(ICE_URL)?token=$(TOKEN)
	unzip -u output/ice.zip -d output

refresh: tools
	$(MAKE) download
	$(MAKE) generate

clean:
	rm -rf output
