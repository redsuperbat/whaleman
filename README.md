<p align="center">
  <img src="/assets/whaleman.png" width="414.05px" height="355.05px" >
</p>


## 🐬 Description 
Whaleman subscribes to a number of docker-compose files in a github repo and automagically update the docker containers running when the docker-compose files change. It adheres to the gitops way, by providing an easy way to manage your docker-compose files in github.

Whaleman is ment to manage itself as well as any number of docker image on a node. It's useful if you have a server at home and just want the server to update the docker cluster when the docker-compose images change.

## 🛥️ Setup 

Whaleman can be run as a binary or as a docker image. Since it's ment to be used in conjunction with docker-compose the suggested way to run Whaleman is with docker-compose in a github repo.

The suggested way is to create a private github repo with all your compose files for a specific node, as well as the compose file for Whaleman. 

### Docker compose
```yaml
version: "3"
services:
  whaleman:
    image: maxrsb/whaleman:latest
    restart: unless-stopped
    environment:
      - GH_PAT=<personal access token>
      - POLLING_INTERVAL_MIN=<number of minutes between to poll>
    ports:
      - 8090:8090
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /home/usr/whaleman:/var/lib/whaleman
```

Grab the raw url for the whaleman compose file and run an instance of Whaleman with the Whaleman compose file as the target.

```shell
docker run -e GH_PAT=<pat> -e GH_COMPOSE_FILES=<url to docker-compose whaleman manifest> -p 8090:8090 --name prewhale maxrsb/whaleman
```

Then curl whaleman so it syncs once

```shell
curl localhost:8090
```

Whaleman will then grab the manifest and spin up another instance of itself watching the manifest which was used to create itself with. Neat huh? 🐳

## 🌟 Upcoming features

- [ ] Whaleman should not kill itself when changes are made to it's own manifest.
- [ ] Whaleman should make sure what is defined in the manifests are running in docker.
