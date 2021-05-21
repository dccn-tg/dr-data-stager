#!/bin/bash

./gen_irods_environment.sh > irods_environments.sh &&
	./data-stager-worker $@
