SHELL = /bin/bash

APP = app
JS = $(APP)/js
CSS = $(APP)/css
BUILD = build
DIST = public
BIN = node_modules/.bin

.PHONY: all watch clean serve

all: node_modules $(DIST)/$(APP).css $(DIST)/$(APP).js $(DIST)/index.html

node_modules:
	npm install

.SECONDARY:

# App
serve: all
	$(BIN)/http-server $(DIST) -p 8000

clean:
	rm -fr $(DIST)/*

watch: node_modules
	$(BIN)/wach -o "$(APP)/**/*" make all

$(DIST)/$(APP).css: $(CSS)/*.css
	cat $^ > $@

$(DIST)/$(APP).js: $(JS)/*.js
	cat $^ > $@

$(DIST)/index.html: $(BUILD)/mbostock/queue/queue.js $(APP)/*.ejs $(APP)/views/*.ejs
	bin/templatize $< $@

# Repo builds

$(BUILD)/%:
	git clone https://github.com/$*.git $@

$(DIST)/tags/%: $(BUILD)/src/%
	mkdir -p $(dir $@)
	ctags -f $@ --fields=+afiKlmnsSzt -R $<
