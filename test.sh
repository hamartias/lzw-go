#!/bin/bash

for i in $( ls test_data ); do
  echo testing: test_data/$i
  go build cmd/main.go
  ./main -c test_data/$i test_data/$i.lzw
  ./main -d test_data/$i.lzw test_data/$i.out
  diff test_data/$i.out test_data/$i
done

rm test_data/*.lzw
rm test_data/*.out
rm ./main
