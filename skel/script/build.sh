#!/usr/bin/env bash
APP_ENV=$1
case $APP_ENV in 
    'dev')
        echo "environment = $APP_ENV"
        ;;
    'qa')
        echo "environment = $APP_ENV"
        ;;
    'preview')
        echo "environment = $APP_ENV"
        ;;
    'prod')
        echo "environment = $APP_ENV"
        ;;
    *)
        APP_ENV='prod'
        echo "environment = $APP_ENV"
        ;;
esac
PROJECT_HOME=$(cd $(dirname $0) && cd .. && pwd -P)
PKG_HOME="$PROJECT_HOME/output"
#程序名称
APP=__APP_NAME__

GO=$GOROOT/bin/go

$GO version

cd $PROJECT_HOME

function doBuild() {
    pwd
    echo -ne "-> building $1 \t ... "
    make #>/dev/null
    if [ $? -eq 0 ]; then
        echo 'done'
    else
        exit 1
    fi
}

rm -rf $PKG_HOME 2>/dev/null
mkdir -p $PKG_HOME/{bin,conf,log,run}
cp conf/$APP_ENV/$APP.yaml $PKG_HOME/conf
cp script/$APP_ENV.sh $PKG_HOME/bin/

echo 'building started'
doBuild
echo 'building finished'

