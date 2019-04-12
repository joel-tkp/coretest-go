#!/bin/bash

# url target
targetUrl="GET http://localhost:8000/api/v1/user?paginated=no"

# continuous attack
# jq -ncM '{method: "GET", url: $targetUrl, body: "Punch!" | @base64, header: {"Content-Type": ["text/plain"]}}' | vegeta attack -format=json -rate=100 | vegeta encode

# reported test
echo $targetUrl | vegeta attack -rate=100/s -duration=1m | vegeta encode > results.json
vegeta report results.json
cat results.json | vegeta report -type='hist[0,2ms,4ms,6ms]'

# realtime chart (using jaggr and jplot)
#echo $targetUrl | \
#    vegeta attack -rate 5000 -duration 10m | vegeta encode | \
#    jaggr @count=rps \
#          hist\[100,200,300,400,500\]:code \
#          p25,p50,p95:latency \
#          sum:bytes_in \
#          sum:bytes_out | \
#    jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 \
#          latency.p95+latency.p50+latency.p25 \
#          bytes_in.sum+bytes_out.sum

