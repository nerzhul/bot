#! /bin/bash

# SSH Init
which ssh-agent || ( apt-get update -y && apt-get install openssh-client -y )
eval $(ssh-agent -s)
ssh-add <(echo "$SSH_PRIVATE_KEY")
git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
mkdir -p ~/.ssh && echo -e "Host gitlab.com\n\tStrictHostKeyChecking no\n\n" >> ~/.ssh/config

# Godep & repo init
go get github.com/tools/godep
mkdir -p $GOPATH/src/$REPO_NAME
mv $CI_PROJECT_DIR/* $GOPATH/src/$REPO_NAME/
cd $GOPATH/src/$REPO_NAME
godep restore