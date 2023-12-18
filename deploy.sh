#!/bin/bash

docker compose down

docker rmi gym-file-server
docker rmi gym-ffmpeg-app

docker compose up -d

