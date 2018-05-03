#!/bin/bash
rm -rf bin

mkdir bin

xgo --targets=windows/*,linux/386,linux/amd64,linux/arm64,linux/arm-7 -out bin/audittrail-client-v$1 .

rm -rf out
mkdir out

mkdir out/audittrail-client-v$1-windows-x64
mkdir out/audittrail-client-v$1-windows-x86
mkdir out/audittrail-client-v$1-linux-x64
mkdir out/audittrail-client-v$1-linux-x86
mkdir out/audittrail-client-v$1-linux-armv7
mkdir out/audittrail-client-v$1-linux-arm64

cp config.json out/audittrail-client-v$1-windows-x64/config.json
cp config.json out/audittrail-client-v$1-windows-x86/config.json
cp config.json out/audittrail-client-v$1-linux-x64/config.json
cp config.json out/audittrail-client-v$1-linux-x86/config.json
cp config.json out/audittrail-client-v$1-linux-armv7/config.json
cp config.json out/audittrail-client-v$1-linux-arm64/config.json


mv bin/audittrail-client-v$1-windows-4.0-amd64.exe out/audittrail-client-v$1-windows-x64/audittrail-client.exe
mv bin/audittrail-client-v$1-windows-4.0-386.exe out/audittrail-client-v$1-windows-x86/audittrail-client.exe
mv bin/audittrail-client-v$1-linux-amd64 out/audittrail-client-v$1-linux-x64/audittrail-client
mv bin/audittrail-client-v$1-linux-386 out/audittrail-client-v$1-linux-x86/audittrail-client
mv bin/audittrail-client-v$1-linux-arm-7 out/audittrail-client-v$1-linux-armv7/audittrail-client
mv bin/audittrail-client-v$1-linux-arm64 out/audittrail-client-v$1-linux-arm64/audittrail-client


zip -j -9 out/audittrail-client-v$1-windows-x64.zip out/audittrail-client-v$1-windows-x64/*
zip -j -9 out/audittrail-client-v$1-windows-x86.zip out/audittrail-client-v$1-windows-x64/*
zip -j -9 out/audittrail-client-v$1-linux-x64.zip out/audittrail-client-v$1-linux-x64/*
zip -j -9 out/audittrail-client-v$1-linux-x86.zip out/audittrail-client-v$1-linux-x86/*
zip -j -9 out/audittrail-client-v$1-linux-armv7.zip out/audittrail-client-v$1-linux-armv7/*
zip -j -9 out/audittrail-client-v$1-linux-arm64.zip out/audittrail-client-v$1-linux-arm64/*

rm -rf out/audittrail-client-v$1-windows-x64
rm -rf out/audittrail-client-v$1-windows-x86
rm -rf out/audittrail-client-v$1-linux-x64
rm -rf out/audittrail-client-v$1-linux-x86
rm -rf out/audittrail-client-v$1-linux-armv7
rm -rf out/audittrail-client-v$1-linux-arm64

rm -rf bin