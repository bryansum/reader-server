SHELL = /bin/bash

APP = app
JS = $(APP)/js
CSS = $(APP)/css
BUILD = build
DIST = public
BIN = node_modules/.bin

GITHUB_REPOS = mbostock/queue
TARGETS = mbostock/queue/queue.js

.PHONY: all build init watch clean serve

all: init build

init: node_modules $(addprefix $(BUILD)/,$(GITHUB_REPOS))

node_modules: package.json
	npm install

.SECONDARY:

# App
serve: all
	$(BIN)/http-server $(DIST) -p 8000

clean:
	rm -fr $(DIST)/*

build: $(addprefix $(DIST)/,$(GITHUB_REPOS))

watch: node_modules build
	$(BIN)/wach -o "$(APP)/**/*" $(MAKE) build

$(DIST)/%: $(BUILD)/%
	mkdir -p $@
	bin/templatize $*

# Builds

$(BUILD)/%:
	git clone https://github.com/$*.git $@

$(DIST)/tags/%: $(BUILD)/src/%
	mkdir -p $(dir $@)
	ctags -f $@ --fields=+afiKlmnsSzt -R $<
