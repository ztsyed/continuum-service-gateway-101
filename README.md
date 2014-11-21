#Continuum Service Gateway- 101

Sample code base that demonstrates:

1. Creating a continuum service gateway
2. Creating a service provider 
3. Binding the service to a job 
4. Consuming the service (using curl inside a capsule)

**Please take a look at the `runDemo.sh` for command line instructions**

#### Quick Start

`runDemo.sh` will run all the commands for you. Once completed, you can connect to the capsule with `apc capsule connect echoClient`. Inside the capsule, run `curl $ECHO0_URI`

#### Clean up

If you want to remove all the changes, run `runDemo.sh clean` 

#### Screencast

[![ScreenCast](http://f.cl.ly/items/0c3T1l0M042i2U2C3v3H/Image%202014-11-21%20at%209.32.16%20PM.png)](https://vimeo.com/112553715)

#### High-Level Sub-System Interaction (probably inaccurate)
![System Components](http://cl.ly/image/2E3N3w0p1n0v/Image%202014-11-21%20at%2012.18.47%20PM.png)
#### High-Level Flow  (probably inaccurate)

![enter image description here](http://bit.ly/1xZOJB4)
