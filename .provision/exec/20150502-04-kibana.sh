#!/bin/bash

KIBANA_DEST="/home/vagrant/kibana"
KIBANA_TMP_PATH="/tmp/kibana.tar.gz"

curl -s -L https://download.elastic.co/kibana/kibana/kibana-4.0.2-linux-x64.tar.gz > $KIBANA_TMP_PATH

tar -xzf $KIBANA_TMP_PATH -C /home/vagrant
mv /home/vagrant/kibana-4.0.2-linux-x64 $KIBANA_DEST
chown vagrant:vagrant -R $KIBANA_DEST

rm $KIBANA_TMP_PATH