# rosgo
[![Build Status](https://travis-ci.org/ppg/rosgo.svg?branch=master)](https://travis-ci.org/ppg/rosgo)
[![Coverage Status](https://coveralls.io/repos/github/ppg/rosgo/badge.svg?branch=master)](https://coveralls.io/github/ppg/rosgo?branch=master)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fppg%2Frosgo.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fppg%2Frosgo?ref=badge_shield)

Package Summary
---------------------------------

**rosgo** is pure Go implementation of [ROS](http://www.ros.org/) client library.

- Author: Akiyoshi Ochiai<akio7141 AT gmail DOT com>
- License: Apache License 2.0
- Source: git [https://github.com/akio/rosgo](https://github.com/akio/rosgo)

Status
---------------------------------

**rosgo** is under development to implement all features of [ROS Client Library Requiements](http://www.ros.org/wiki/Implementing%20Client%20Libraries).

At present, following basic functions are provided.

- Parameter API (get/set/search....)
- ROS Slave API (with some exceptions)
- Publisher/Subscriber API (with TCPROS)

Building
---------------------------------

Setup environmet variable:

     export GOPATH=${path/to/rosgo/dir}


Build rosgo library:

     go install ros


Examples programs:

     go install test_listener test_talker test_param

Example executables are placed in `bin` directory.


*In future release, the build system will be integrated with [catkin](http://www.ros.org/wiki/catkin).*


See also
---------------------------------

- [rosgo in ROS Wiki](http://www.ros.org/wiki/rosgo)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fppg%2Frosgo.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fppg%2Frosgo?ref=badge_large)