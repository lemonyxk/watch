```text
install:
go install github.com/lemoyxk/watch@v0.0.0-20220524105950-e29f0e06b0b1

use:
1.create file name .watch in you project
2.write start shell
3.watch --path /you/project/path
4.or you can cd /you/project/path and run watch

watch file like this:
[ignore]
.idea
.git
.watch
vendor
node_modules
main

[command]
#
# you command
# whatever will be executed
# example: ls
#
# dir or file: you command
# variable: $DIR = project dir
# execute when the dir or file has changed
# example: echo $PATH
#          main.go: echo hello
#          .: go run main.go
#          .: go run $DIR/main.go
#          ./src: ls
#          ./test*: echo test
#

go run main.go

[host]
# http://127.0.0.1:12365
# http://127.0.0.1:12365/reload?name=.
# http://127.0.0.1:12365/reload
# http://127.0.0.1:12365/reload?name=$DIR
# http://127.0.0.1:12365/reload?name=./src
# http://127.0.0.1:12365/reload?name=./src*

http://127.0.0.1:12365/reload