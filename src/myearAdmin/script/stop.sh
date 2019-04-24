#!/bin/bash
ps aux | grep -w myearAdmin | grep -v "grep" | awk '{print $2}' |xargs kill -9
