#!/bin/bash

# Assumes you have protobuf_inspector installed, see https://github.com/mildsunrise/protobuf-inspector
# Aslso assumes ~/.local/bin isn't on your PATH
# Sed command from https://stackoverflow.com/questions/17998978/removing-colors-from-output/51141872#51141872
~/.local/bin/protobuf_inspector < ./dji_meta.bin | sed 's/\x1B[@A-Z\\\]^_]\|\x1B\[[0-9:;<=>?]*[-!"#$%&'"'"'()*+,.\/]*[][\\@A-Z^_`a-z{|}~]//g' > ./protobuf_inspector_dji_dbmeta.txt

# This creates a very, very large output!
~/.local/bin/protobuf_inspector < ./dji_dbgi.bin | sed 's/\x1B[@A-Z\\\]^_]\|\x1B\[[0-9:;<=>?]*[-!"#$%&'"'"'()*+,.\/]*[][\\@A-Z^_`a-z{|}~]//g' > ./protobuf_inspector_dji_dbgi.txt
