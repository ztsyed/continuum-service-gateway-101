#!/bin/bash

if [ $1 == "clean" ] ; then

	echo "#########################"
	echo "UnBind Service to the job"
	echo "#########################"

	apc service unbind echo0 -j echoClient

	echo "#########################"
	echo "Deleting Capsule echoClient"
	echo "#########################"

	apc capsule delete echoClient

	echo "#########################"
	echo "Deleting Service echo0"
	echo "#########################"

	apc service delete echo0

	echo "##############################"
	echo "Unregistering Provider echo"
	echo "##############################"

	apc provider delete echo

	echo "#####################################"
	echo "Deleting Service Provider echo-server"
	echo "#####################################"

	apc app delete echo-server

	echo "#################################"
	echo "Deleting Service Gateway echo-sg"
	echo "#################################"


	apc app delete echo-sg

else
	
	echo "#########################"
	echo "Creating Service Gateway"
	echo "#########################"

	cd echo-sg
	apc app create echo-sg --disable-routes
	apc gateway promote echo-sg --type echosg

	echo "#########################"
	echo "Creating Service Provider"
	echo "#########################"

	cd ../echo-server
	apc app create echo-server --disable-routes
	apc app update echo-server --port-add 3000
	apc app start echo-server

	echo "#########################"
	echo "Registering Service Provider"
	echo "#########################"

	apc provider register echo --type echosg --job echo-server -port 3000 --url http://can-be-anything.com/echo

	echo "#########################"
	echo "Creating Service echo0"
	echo "#########################"

	apc service create echo0 --provider echo

	echo "########################################"
	echo "Creating Capsule to test service binding"
	echo "########################################"

	apc capsule create echoClient --image linux -ae

	echo "#########################"
	echo "Bind Service to the job"
	echo "#########################"

	apc service bind echo0 -j echoClient

fi
