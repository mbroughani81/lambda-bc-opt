curl -X POST "http://10.10.0.1:3233/api/v1/namespaces/_/actions/visitor-counter?blocking=true&result=true" \
    -H "Authorization: Basic $(echo -n '23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP' | base64 -w 0)"
