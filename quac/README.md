# QUickly Amalgamate Concepts 

---

The intention behind this program is to provide a space for quickly:
 - entering ideas from the command line
 - organizing ideas with tags and timestamps
 - scanning in handwritten notes for later transcription
 - viewing and developing your ideas further

### Installation From Source

1. make sure you [have Go installed][1] and [put $GOPATH/bin in your $PATH][2]
2. run `go get github.com/rigelrozanski/thranch`
3. run `go install` while navigated to this directory
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

The scan functionality is provided to quickly scan in a sheet of paper with
notes written in multiple orientations (currently only, rightside-up or upside
down).  The scan image is then reoriented and seperated out into smaller,
untranscribed, untagged images which can then be quickly tagged and later
transcribed with the `qu wc` command. 


### Using SPLIT

 - the SPLIT keyword takes the most recent above tags AS WELL AS any new provided tags "SPLIT newtag1,newtag2"

### Using the browser

the tag browser can be accessed through `qu ls`, Once launched the following commands can be used:
 `q` - quit
 `h` - go to previous list
 `j` - move down list
 `k` - move up list
 `l` - find associated tags to the current list item (as well as highlighted tags) 
 `Ctrl-l` - highlight the previous tag and drill to associated tags 
 `f` - find associated files with current highlighted tags 
 `Enter` - Either open the highlighted file, or all files associated with the highlighted tags

### License

Quick Ideas is released under the Apache 2.0 license.
