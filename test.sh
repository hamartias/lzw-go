#!/bin/bash

for i in $( ls test_data ); do
  echo testing: test_data/$i
  go run lzw.go -c test_data/$i test_data/$i.lzw
  go run lzw.go -d test_data/$i.lzw test_data/$i.out
  diff test_data/$i.out test_data/$i
done

rm test_data/*.lzw
rm test_data/*.out
