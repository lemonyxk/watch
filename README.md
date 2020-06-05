```text
install:
go get -u github.com/Lemo-yxk/go-watch

use:
1.create file name .watch in you project
2.write start shell
3.go-watch --path /you/project/path
4.or you can cd /you/project/path and run go-watch

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
echo $PATH
main.go: go run main.go
test.log: cat test.log