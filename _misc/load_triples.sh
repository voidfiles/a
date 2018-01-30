#!/bin/bash

for i in $(ls $1); do
    $2 load \
  	 -c $3/conf/cayley.yml \
     --dbpath $4 \
  	 -i ${1}/${i}
done
