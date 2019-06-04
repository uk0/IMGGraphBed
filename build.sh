#!/usr/bin/env bash


PASS=emhhbmdq111


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags static ./

go build  -o mac_IMG -tags static ./


# scp -i ~/Desktop/key/pvkey IMGGraphBed  root@10.50.40.10:/opt/
# scp -i ~/Desktop/key/pvkey -r web  root@10.50.40.10:/opt/



  sshpass -p ${PASS}  scp  IMGGraphBed  root@172.17.0.200:/opt/
#  sshpass -p ${PASS}  scp  conf.toml   root@172.17.0.200:/opt/
  sshpass -p ${PASS}  scp  -r web   root@172.17.0.200:/opt/