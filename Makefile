CMDS := $(shell find cmd -mindepth 1 -maxdepth 1 -type d)
BINARIES := $(foreach dir,$(CMDS),$(notdir $(dir)))

.PHONY: all clean

all: $(BINARIES)

$(BINARIES):
	@echo Building $@
	@go build -o $@ ./cmd/$@

clean:
	rm -f $(BINARIES)
