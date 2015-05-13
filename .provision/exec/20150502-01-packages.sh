#!/bin/sh

add-apt-repository -y ppa:webupd8team/java > /dev/null 2>&1

apt-get update > /dev/null

echo debconf shared/accepted-oracle-license-v1-1 select true | debconf-set-selections
echo debconf shared/accepted-oracle-license-v1-1 seen true | debconf-set-selections

apt-get install -y vim mc htop \
    git mercurial \
    oracle-java8-installer \
    curl \
    binutils bison gcc build-essential > /dev/null
