#!/bin/bash

ELASTIC_TMP_PATH="/tmp/elastic.deb"

curl -s -L https://download.elastic.co/elasticsearch/elasticsearch/elasticsearch-1.5.2.deb > $ELASTIC_TMP_PATH

sudo dpkg -i $ELASTIC_TMP_PATH

update-rc.d elasticsearch defaults 95 10

/etc/init.d/elasticsearch start

rm $ELASTIC_TMP_PATH