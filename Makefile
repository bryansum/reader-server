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
	mkdir -p $(DIST)

node_modules: package.json
	npm install

.SECONDARY:

# App
serve: all
	$(BIN)/http-server $(DIST) -p 8000

clean:
	rm -fr $(DIST)/*

build: $(DIST)/$(APP).css $(DIST)/$(APP).js $(DIST)/index.html

watch: node_modules
	$(BIN)/wach -o "$(APP)/**/*" $(MAKE) build

$(DIST)/$(APP).css: $(CSS)/*.css
	cat $^ > $@

$(DIST)/$(APP).js: $(JS)/*.js
	cat $^ > $@

$(DIST)/index.html: $(addprefix $(BUILD)/,$(TARGETS)) $(APP)/*.ejs $(APP)/views/*.ejs
	bin/templatize $< $@

# Repo builds

$(BUILD)/%:
	git clone https://github.com/$*.git $@

$(DIST)/tags/%: $(BUILD)/src/%
	mkdir -p $(dir $@)
	ctags -f $@ --fields=+afiKlmnsSzt -R $<
