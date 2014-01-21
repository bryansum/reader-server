SHELL = /bin/bash

APP=app
JS=$(APP)/js
CSS=$(APP)/css
DIST=public
BIN=node_modules/.bin

.PHONY: all watch clean serve

all: node_modules $(DIST)/$(APP).css $(DIST)/$(APP).js $(DIST)/index.html $(DIST)/repo/mbostock/queue $(DIST)/highlight/mbostock/queue $(DIST)/tags/mbostock/queue

node_modules:
	npm install

.SECONDARY:

# App
serve: all
	$(BIN)/http-server $(DIST) -p 8000

clean:
	rm -f $(DIST)/$(APP).css $(DIST)/$(APP).js

watch: node_modules
	$(BIN)/wach -o "$(APP)/**/*" make all

$(DIST)/$(APP).css: $(CSS)/$(APP).css
	cat $^ > $@

$(DIST)/$(APP).js: $(JS)/$(APP).js
	cat $^ > $@

$(DIST)/index.html: $(APP)/index.html
	cat $^ > $@

# Repo builds

$(DIST)/repo/%:
	mkdir -p $(dir $@)
	git clone https://github.com/$*.git $@

$(DIST)/highlight/%: $(DIST)/repo/%
	mkdir -p $@
	bin/highlight $</queue.js $@/queue.js

$(DIST)/tags/%: $(DIST)/repo/%
	mkdir -p $(dir $@)
	ctags -f $@ --fields=+afiKlmnsSzt -R $<
