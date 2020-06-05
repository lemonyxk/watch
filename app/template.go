package app

var Template = `
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
`
