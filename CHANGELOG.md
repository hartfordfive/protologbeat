
### Version 0.2.0
- Added support for receiving GELF format messages by setting `enable_gelf: true`
- Fixed issue where (more commonly when getting high volumes of messages) the  `processMessage` goroutine might have it's byte buffer modified before it actually executes with the original payload/message. (Related to original PR [#8](https://github.com/hartfordfive/protologbeat/pull/8), credit to [vcostet](https://github.com/vcostet))


### Version 0.1.1
- Added Dockerfile and seperate `protologbeat-docker.yml` config file to be used by docker image
- Updated default `protologbeat.yml` to have bare-minimum config values
- Added `build-bin.sh` build script to simplify compiling the binary for the most common platforms