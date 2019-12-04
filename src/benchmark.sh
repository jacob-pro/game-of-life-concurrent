#!/usr/bin/env bash

go test -c -o gameoflife.test
rm -r benchmarks
mkdir benchmarks

benchtime=10x
benchmark=512x512x8

for impl in serial parallel halo parallelshared rust
do
    echo implementation ${impl} on our solution
    ./gameoflife.test -test.run XXX -test.bench /${benchmark} -test.benchtime ${benchtime} -i ${impl} > benchmarks/${impl}.txt
done

echo baseline solution
./baseline.test -test.run XXX -test.bench /${benchmark} -test.benchtime ${benchtime} > benchmarks/baseline.txt
