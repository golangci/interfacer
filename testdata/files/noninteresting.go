package foo

import (
	"os"
)

type EmptyIface interface{}

type UninterestingMethods interface {
	Foo() error
	bar() int
}

type InterestingUnexported interface {
	Foo(f *os.File) error
	bar(f *os.File) int
}
