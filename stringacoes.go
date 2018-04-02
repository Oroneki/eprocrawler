package main

import (
	"fmt"
	"regexp"
)

func processoPath(dst string, processo string) string {
	regexNumeroProcessoLimpar := regexp.MustCompile(`\D`)
	num := regexNumeroProcessoLimpar.ReplaceAllString(processo, "")
	return fmt.Sprintf(`%s%s.pdf`, dst, num)
}

// func main() {
// 	fmt.Println(os.TempDir())
// }
