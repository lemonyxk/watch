package app

var Template = []string{
	"[ignore]",
	".idea",
	".git",
	"vendor",
	"node_modules",
	"main",

	"[start]",
	"go build -o main main.go",
	"./main",
}
