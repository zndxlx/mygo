#!/bin/bash
ps aux | grep -w opencast | grep -v "grep" | awk '{print $2}' |xargs kill -USR2 