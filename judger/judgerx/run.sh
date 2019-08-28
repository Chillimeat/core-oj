#!/usr/bin/env sh

# nohup ./judger &
SOCK=/judger_tools/socks/${NAME}.sock

echo hello judger!


# chmod 711 /

# groupadd core-qwq
# useradd -g core-qwq orz-qwq

# $SOCK
# touch $SOCK
# chgrp core-qwq $SOCK
# chmod 770 $SOCK

ban() {
    setfacl -R -m u:1001:--- $1
    return 0
}

ban /judger_tools
ban /checker_tools
ban /codes
ban /usr
ban /root
ban /home
# $ban /codes
ban /problems
ban /sbin
ban /etc


touch /presu.txt

./judger -addr $SOCK


