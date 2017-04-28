# Protologbeat

## Description


Protologbeat is an application that accepts logs via tcp or udp.  It can act accept plain-text or JSON logs and also act as a syslog destination replacement.  It's intended as a replacement for [udplogbeat]().


Ensure that this folder is at the following location:
`${GOPATH}/github.com/harfordfive`

## Getting Started with Protologbeat

### Configuration Options

- `protologbeat.port` : The UDP port on which the process will listen (Default = 5000)
- `protologbeat.max_message_size` : The maximum accepted message size (Default = 1024)
- `protologbeat.enable_syslog_format_only` : Boolean value indicating if only syslog messages should be accepted. (Default = false)
- `protologbeat.enable_json_validation` : Boolean value indicating if JSON schema validation should be applied for `json` format messages (Default = false)
- `protologbeat.publish_failed_json_invalid` : Boolean value indicating if JSON objects should be sent serialized in the event of a failed validation.  This will add the `_protologbeat_jspf` tag. (Default = false)
- `protologbeat.json_document_type_schema` :  A hash consisting of the Elasticsearch type as the key, and the absolute local schema file path as the value.

### Configuration Example

The following are examples of configuration blocks for the `protologbeat` section.  

1. Sample configuration block for a syslog replacement
_sample/config1.yml

2. Sample configuration block that enforces JSON schema only for indicated Elasticsearch types:
_sample/config2.yml

3. Sample configuration block that enforces JSON schema only for indicated Elasticsearch types:
_sample/config3.yml

4. Sample configuration block that enforces JSON schema for all events:
_sample/config4.yml


JSON schemas can be automatically generated from an object here: http://jsonschema.net/.  You can also view the included sample schemas `_samples/app1_schema.json` and `_samples/app2_schema.json` as examples.

#### Considerations

If you intend on using this as a drop-in replacement to logging with Rsyslog, this method will not persist your data to a file on disk.  
If protologbeat is down for any given reason, messages sent to the configured UDP port will never be processed or sent to your ELK cluster.
If you need 100% guarantee each message will be delivered at least once, this may not be the best solution for you.  
If some potential loss of log events is acceptable for you, than this may be a reasonable solution for you.

*Please note the current date/time is automatically added to each log entry.*

### Sample Clients

Please see the `sample_clients/` directory for examples of clients in various languages.



### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Protologbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Protologbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/harfordfive/protologbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Protologbeat run the command below. This will generate a binary
in the same directory with the name protologbeat.

```
make
```


### Run

To run Protologbeat with debugging output enabled, run:

```
./protologbeat -c protologbeat.yml -e -d "*"
```


### Test

To test Protologbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/protologbeat.template.json and etc/protologbeat.asciidoc

```
make update
```


### Cleanup

To clean  Protologbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Protologbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/harfordfive
cd ${GOPATH}/github.com/harfordfive
git clone https://github.com/harfordfive/protologbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
