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

### Details on use of `qu scan`

The scan functionality is provided to quickly scan in a sheet of 
paper with notes written in multiple orientations. The user is expected
to circle each orientation (noon, 3am, 6am, 9am) with different colour
markers. Additionally the markers are to be using for caquacration:

 - The top-left 1/2 inch squared is reserved for the caquacration
   - Caquacration markers should be drawn as four stacked horizontal lines representing:
     - Noon
     - Quarter-Past
     - Half-Past
     - Quarter-To

### Using SPLIT

 - the SPLIT keyword takes the most recent above tags AS WELL AS any new provided tags "SPLIT newtag1,newtag2"

### License

Quick Ideas is released under the Apache 2.0 license.
