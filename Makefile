SHELL = /bin/bash

.PHONY: all

all: node_modules tags/queue-1.0.7/tags highlight/queue-1.0.7/queue.html

node_modules:
	npm install

.SECONDARY:

zip/queue-%.zip:
	mkdir -p $(dir $@)
	curl "https://github.com/mbostock/queue/archive/v$*.zip" -L -o $@.download
	mv $@.download $@

src/queue-%: zip/queue-%.zip
	mkdir -p $(dir $@)
	unzip -d src $<

highlight/queue-%/queue.html: src/queue-%/queue.js
	mkdir -p $(dir $@)
	bin/highlight $< $@

tags/queue-%/tags: src/queue-%/queue.js
	mkdir -p $(dir $@)
	ctags -f $@ --fields=+afiKlmnsSzt -R $<
