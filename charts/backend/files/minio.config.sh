#!/bin/env bash

set -e
set -x

. /config.env

echo Configuring alias

echo '{"url":"'$_MINIO_ADDRESS'","accessKey":"'$MINIO_ROOT_USER'","secretKey":"'$MINIO_ROOT_PASSWORD'","api":"S3v4","path":"auto"}' | mc alias i localtemp


echo Configuring policies

ls /policies/*.policy.json | while read POLICYFILE; do
  # https://gist.github.com/magnetikonline/90d6fe30fc247ef110a1
  POLICYNAME=${POLICYFILE##*/}
  POLICYNAME=${POLICYNAME%%.*}
  echo : creating policy $POLICYNAME from $POLICYFILE
  mc admin policy create localtemp $POLICYNAME $POLICYFILE
done

echo Creating replication admin user

. /backend.env
mc admin user add localtemp $ACCESS_KEY $SECRET_KEY
mc admin policy attach localtemp replicationAdmin --user $ACCESS_KEY

#echo Disabling root access
#mc admin config set localtemp api root_access=off
