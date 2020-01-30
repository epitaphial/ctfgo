#!/bin/bash
go env -w GOPROXY=https://goproxy.io,direct
bee run