#!/bin/bash

for i in $(cat $1); do
    (cd $2 && curl -O -L $i)
done
