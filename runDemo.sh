#!/bin/bash -e

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

apc provider register echo --type echosg --job echo-server -port 3000 --url http://user:pass@example.com/ping

echo "#########################"
echo "Creating Service"
echo "#########################"

apc service create echoer-0 --provider echo

echo "#########################"
echo "Creating Capsule to test service binding"
echo "#########################"

apc capsule create echoClient --image linux -ae
echo "Bind Service"
apc service bind echoer-0 -j echoClient
