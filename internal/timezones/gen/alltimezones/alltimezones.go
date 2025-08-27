package main

import (
	"fmt"
	"strings"

	"github.com/ninedraft/daily-bacon/internal/timezones"
)

func main() {
	fmt.Println(strings.Join(timezones.All(), "\n"))
}
