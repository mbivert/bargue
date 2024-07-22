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

.PHONY: test-site
test-site: input/ site-pp thumbnails
	@echo Building test site from input/ to output/...
	@echo '"file:///home/mb/gits/bargue/output"' > input/db/baseurl.json
	@${dtmpl} input/ output/

.PHONY: site
site: input/ site-pp thumbnails
	@echo Building site from input/ to output/...
	@echo '"https://bargue.mbivert.com"' > input/db/baseurl.json
	@${dtmpl} input/ output/

bargue-pp: bargue-pp.go
	@echo Building $@...
	@go build $<

.PHONY: site-pp
site-pp: bargue-pp input/
	@echo Preprocessing...
	@./bargue-pp input/

.PHONY: thumbnails
thumbnails: mkthumbnails.sh
	@echo Making thumbnails/smaller images...
	@sh mkthumbnails.sh
