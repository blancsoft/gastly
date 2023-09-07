#!/usr/bin/env sh

if [ ! -f "coverage.out" ]; then
  echo "Error: coverage.out is missing"
  exit 1
fi

totalCoverage=`go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'`
echo "Required coverage threshold   : $TESTCOVERAGE_THRESHOLD %"
echo "Current test coverage         : $totalCoverage %"

if (( $(echo "$totalCoverage $TESTCOVERAGE_THRESHOLD" | awk '{print ($1 > $2)}') )); then
    echo "OK"
else
    echo "Current test coverage is below threshold."
    exit 1
fi
