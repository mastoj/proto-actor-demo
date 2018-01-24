#!/bin/bash

docker-compose up --build --force-recreate --scale fsharp-worker=0 --scale go-worker=1
