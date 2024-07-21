dtmpl = dtmpl

.PHONY: all
all: site

.PHONY: site
site: input/
	@echo Building site from input/ to output/...
	@${dtmpl} input/ output/
