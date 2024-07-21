dtmpl = dtmpl

.PHONY: all
all: site

.PHONY: help
help:
	@echo Available targets:
	@echo "	site/all   : build the website"
	@echo "	site-pp    : run bargue-pp on the site"
	@echo "	bargue-pp  : build bargue-pp"
	@echo "	check-deps : check for dependencies"

.PHONY: check-deps
check-deps:
	@echo Checking dependencies...
	@sh check-deps.sh
	@echo All dependencies found

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