
docker stop $(docker ps -f name=judger -q)
docker stop $(docker ps -f name=compiler -q)
