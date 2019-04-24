#!/bin/bash

GO_PROJ_DIR=`cd ../../; pwd`
PROJ=opencast
OPENCAST_DIR=$GO_PROJ_DIR/src/$PROJ
BUILD_DIR=$OPENCAST_DIR/build
CONF_DIR=$OPENCAST_DIR/conf


export GOPATH=$GO_PROJ_DIR

go install

if [ $? -eq 0 ];then
    cp -rfp $GO_PROJ_DIR/bin/$PROJ $BUILD_DIR/$PROJ
    cp -rfp $CONF_DIR/* $BUILD_DIR/$PROJ

    cd $BUILD_DIR && tar -czvf $PROJ.tar.gz $PROJ
else
    echo "编译失败"
fi



