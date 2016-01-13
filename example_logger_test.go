package main

import (
	"fmt"
	"testing"
)

func TestBaseceLogger(t *testing.T) {
	GsInfo.Println("test logger ok.")
	fmt.Println("test logger ok2.")
}
