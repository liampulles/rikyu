#!/usr/bin/env bash
FULLPATH="$( realpath ./openapi/reference/rikyud.v1.yaml )"
#docker run -p 9119:8000 -v $FULLPATH:/api.yaml danielgtaylor/apisprout /api.yaml
docker run --init --rm -it -p 9119:4010 -v "$FULLPATH":/api.yaml -P stoplight/prism:3 mock -h 0.0.0.0 /api.yaml