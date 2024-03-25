#!/bin/bash

rm -rf "*.png"

echo -e "\n🚀 Running Rust version"
cargo build --release 2> /dev/null
./target/release/nanoray

echo -e "\n🚀 Running Go version"
go build
./nanoray

echo -e "\n🚀 Running Node.js version"
node nanoray.js