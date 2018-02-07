path := $(shell pwd)

red := \033[1;31m
green := \033[1;32m
yellow := \033[1;33m
purple := \033[1;35m
blue := \033[1;36m
nocolor := \033[0m

define print_title
	printf "$(yellow)\n=============================[ $(1) ]==============================\n\n$(nocolor)"
endef


test:
	@ $(call print_title,DONE)


.PHONY: test
