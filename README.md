<p align="center">
  <img src="/assets/whaleman.png" width="414.05px" height="355.05px" >
</p>


## Description
Whaleman subscribes to a number of docker-compose files in a github repo and automagically update the docker containers running when the docker-compose files change. It adheres to the gitops way, by providing an easy way to manage your docker-compose files in github.

Whaleman is ment to manage itself as well as any number of docker image on a node. It's useful if you have a server at home and just want the server to update the docker cluster when the docker-compose images change.

## Setup

Whaleman can be run as a binary or as a docker image. Since it's ment to be used in conjunction with docker-compose the suggested way to run Whaleman is with docker-compose in a github repo.





