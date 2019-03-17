#!/usr/bin/env bash
APP_ENV=$1
case $APP_ENV in 
    'dev')
        ;;
    'qa')
        ;;
    'prod')
        ;;
    *)
        APP_ENV='prod'
        ;;
esac
echo "environment = $APP_ENV"
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
cp scripts/conf/$APP_ENV.conf $PKG_HOME/conf/$APP.conf

echo 'building started'
doBuild
echo 'building finished'

