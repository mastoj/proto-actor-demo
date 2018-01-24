#!/bin/bash
protoc ./messages/messages.proto -I=. --csharp_out=./fsharp-node/messages --csharp_opt=file_extension=.g.cs --grpc_out . --plugin=protoc-gen-grpc=/Users/tomasjansson/.nuget/packages/grpc.tools/1.6.1/tools/macosx_x64/grpc_csharp_plugin

protoc --gogoslick_out=./go-node ./messages/messages.proto