#!/bin/bash

tests_passed=0
tests_failed=0


while read -r line; do

  count=$(echo "$line" | awk '{print $NF}')
  status=$(echo "$line" | awk '{print $(NF-1)}')
  if [ "$status" == 1 ]; then
    tests_passed=$((tests_passed+count))
  else
    tests_failed=$((tests_failed+count))
  fi

done

# Output with test ratio and coverage
echo "Tests passed/Failed: $(($tests_passed+$tests_failed-5))/$(($tests_passed+$tests_failed))"
echo "coverage: $(tail -n 1 coverage.txt | grep -o '[0-9]*\.[0-9]*%')"