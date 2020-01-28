# QUickly Amalgamate Concepts 

---

The intention behind this program is to provide a space for quickly:
 - entering ideas from the command line
 - organizing ideas with tags and timestamps
 - scanning in handwritten notes for later transcription
 - viewing and developing your ideas further

### Installation From Source

1. make sure you [have Go installed][1] and [put $GOPATH/bin in your $PATH][2]
2. run `go get github.com/rigelrozanski/qi`
3. run `go install` from the repo directory
4. set the `qi` working directory in a config file in `~/qi_config.txt` (see `example_config.txt`)

[1]: https://golang.org/doc/install
[2]: https://github.com/tendermint/tendermint/wiki/Setting-GOPATH 

### File Structure

filestructure:
               ./ideas/a,123456,YYYYMMDD,eYYYYMMDD,cYYYYMMDD,c432978,c543098...,tag1,tag2,tag3...
               ./qi
               ./log
               ./config
               ./working_files
               ./working_content
123456 = id
c123456 = consumes-id
YYYYMMDD = creation date
eYYYYMMDD = last edited date
cYYYYMMDD = consumed date

### License

Quick Ideas is released under the Apache 2.0 license.
