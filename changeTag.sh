#!/bin/bash
env=""
originEnv=""
if [ $2 != "" ] && [ $2 != "prod" ]
then
    originEnv=$2
    env=$2"_"
fi
sed "s/{SYMPER_IMAGE}/$1/g" k8s/go_deployment.yaml > k8s/new_go_deployment.yaml
sed "s/{ENVIRONMENT}/$originEnv/g;s/{ENVIRONMENT_}/$env/g;s/{POSTGRES_USER}/$3/g;s/{POSTGRES_PASSWORD}/$4/g" k8s/config_env.yaml > k8s/new_config_env.yaml
rm -rf k8s/go_deployment.yaml
rm -rf k8s/config_env.yaml