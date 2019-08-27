#!/usr/bin/env sh

# nohup ./judger &
SOCK=/judger_tools/socks/${NAME}.sock

echo hello judger!


# groupadd core-qwq
# useradd -g core-qwq orz-qwq

# $SOCK
# touch $SOCK
# chgrp core-qwq $SOCK
# chmod 770 $SOCK

./judger -addr $SOCK


