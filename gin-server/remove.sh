
docker rm $(docker ps -af name=judger -q)
docker rm $(docker ps -af name=compiler -q)
