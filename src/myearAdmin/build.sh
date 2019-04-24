#!/bin/bash
PROJ_NAME=myearAdmin

GO_PROJ_DIR=`cd ../../; pwd`
PROJ_DIR=$GO_PROJ_DIR/src/$PROJ_NAME
BUILD_DIR=$PROJ_DIR/build
CONF_DIR=$PROJ_DIR/conf
SCRIPT_DIR=$PROJ_DIR/script

export GOPATH=$GO_PROJ_DIR

go install

if [ $? -eq 0 ];then
    mkdir -p $BUILD_DIR/$PROJ_NAME/logs
    cp -rfp $GO_PROJ_DIR/bin/$PROJ_NAME $BUILD_DIR/$PROJ_NAME
    cp -rfp $CONF_DIR $BUILD_DIR/$PROJ_NAME
    cp -rfp $SCRIPT_DIR/* $BUILD_DIR/$PROJ_NAME
    cd $BUILD_DIR && tar -czvf $PROJ_NAME.tar.gz $PROJ_NAME
else
    echo "编译失败"
fi



