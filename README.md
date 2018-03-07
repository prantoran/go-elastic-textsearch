### go-elastic-textsearch
- cloned from https://github.com/olivere/elastic-with-docker

#### To run:
- $ docker-compose up


#### Creating vendor:
- install govendor
- $ govendor init
- $ govendor list
- $ govendor add +external // add all external packages
- $ govendor fetch packagepath // add a specific package


#### If 'no such images' error:
- $ docker-compose ps
- $ docker-compose rm // remove all old images
- rebuild again

if still persistent
- $ docker-compose down


#### Removing dangling docker images:
- $ docker rmi -f $(docker images -f dangling=true -q)
- warning: do not simply remove images with <noname> tags, because other images may be dependent

#### Removing all containers
- $ docker rm $(docker ps -a -q)

#### Removing all images
- $ docker rmi $(docker images -q)

#### Removing images with a specific string pattern
- $ docker images -a | grep "pattern" | awk '{print $3}' | xargs docker rmi

#### Host port: Container port
- hostport:containerport
- service to service communications use containerport
- hostport allows a service to be accesible outside the swarm as well


https://www.elastic.co/guide/en/elasticsearch/reference/5.6/nested.html

#### Curl Commands
- alias curlpost='curl -H "Content-type: application/json" -X POST -d'
- curlpost "@body.json" 127.0.0.1:6969/setmap
- curl 127.0.0.1:6969/indexexists/laws
- curl 127.0.0.1:6969/deleteindex/laws
- curl 127.0.0.1:9200/laws
- curlpost "@body1.json" 127.0.0.1:6969/insert/single/laws/details
- curl 127.0.0.1:6969/get/single/laws/details/380


