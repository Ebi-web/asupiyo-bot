.PHONY: copy-env setup

copy-env:
	cp .env.example .env.local

setup: copy-env
