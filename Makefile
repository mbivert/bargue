dtmpl = dtmpl

.PHONY: all
all: site

.PHONY: site
site: site-pp input/
	@echo Building site from input/ to output/...
	@${dtmpl} input/ output/

bargue-pp: bargue-pp.go
	@echo Building $@...
	@go build $<

.PHONY: site-pp
site-pp: bargue-pp input/
	@echo Preprocessing...
	@./bargue-pp input/