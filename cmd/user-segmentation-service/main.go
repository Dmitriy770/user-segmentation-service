package main

import (
	"fmt"

	"github.com/Dmitriy770/user-segmentation-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
