#!/usr/bin/env bash
set -m

export MINIO_BROWSER_REDIRECT_URL=http://127.0.0.1:9002/
export MINIO_VOLUMES="/mnt/data"
export MINIO_IDENTITY_PLUGIN_URL="http://127.0.0.1:8593/auth"
export MINIO_IDENTITY_PLUGIN_ROLE_POLICY="ethauthPolicy"
export MINIO_IDENTITY_PLUGIN_TOKEN="Bearer TOKEN"
export MINIO_IDENTITY_PLUGIN_ROLE_ID="ethauth"
export MINIO_IDENTITY_PLUGIN_COMMENT="External Identity Management using ethauth"

export MINIO_ROOT_USER=$(tr -dc "A-Za-z0-9" <dev/urandom | head -c20)
export MINIO_ROOT_PASSWORD=$(tr -dc "A-Za-z0-9" <dev/urandom | head -c20)
export MINIO_ROOT_ACCESS=on

minio "$@" &
sleep 1

echo '{"url":"http://127.0.0.1:9000","accessKey":"'$MINIO_ROOT_USER'","secretKey":"'$MINIO_ROOT_PASSWORD'","api":"S3v4","path":"auto"}' | mc alias i localtemp

echo '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["arn:aws:s3:::${jwt:preferred_username}","arn:aws:s3:::${jwt:preferred_username}/*"]}]}' | mc admin policy create localtemp ethauthPolicy /dev/stdin || echo $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD

mc admin config set localtemp api root_access=off

mc alias rm localtemp

fg




