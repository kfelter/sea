package docgen

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func GenDoc(path string) error {
	err := filepath.Walk(path, filewalker)
	if err != nil {
		return fmt.Errorf("error walking file path %s: %v", path, err)
	}
	return nil
}

func filewalker(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	return printVars(path, info)
}

func printVars(path string, info fs.FileInfo) error {
	if shouldNotParse(path, info) {
		// fmt.Printf("skipping %s\n", path)
		return nil
	}
	// fmt.Printf("parsing %s\n", path)
	boxes, err := FindCalls(path)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Boxes:", boxes)
	if len(boxes) > 0 {
		printHeader(path)
	}
	for _, b := range boxes {
		printMD(b)
	}
	return nil
}

func printHeader(path string) {
	fmt.Printf("## %s\n", path)
	fmt.Printf("| NAME | DEFAULT | USAGE | TYPE |\n| --- | --- | --- | --- |\n")
}

func printMD(b Box) {
	if len(b.args) < 1 {
		log.Println("bad box", b.args)
		return
	}
	if b.args[0] == `"sea.Load"` {
		if len(b.args) < 4 {
			log.Println("bad box", b.args)
			return
		}
		fmt.Printf("| `%s` | `nil` | `%s` | `%s` |\n", b.args[1], b.args[2], b.args[3])
		return
	}
	if len(b.args) < 5 {
		log.Println("bad box", b.args)
		return
	}
	fmt.Printf("| `%s` | `%s` | `%s` | `%s` |\n", b.args[1], b.args[2], b.args[3], b.args[4])
}

func shouldNotParse(path string, info fs.FileInfo) bool {
	return strings.Contains(path, "vendor") ||
		info.IsDir() ||
		!strings.HasSuffix(info.Name(), ".go")
}

// type v struct {
// 	name        string
// 	seaType     string
// 	seaDefault  string
// 	required    string
// 	description string
// }

// func getVars(path string) ([]v, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
// 	vars := []v{}
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		if line := scanner.Text(); strings.Contains(line, "sea.Load") {
// 			vars = append(vars, pv(line))
// 		}
// 	}
// 	return vars, scanner.Err()
// }

// func pv(line string) v {
// 	return v{getName(line), getType(line), getDefault(line), getRequired(line), getDescription(line)}
// }

// func getName(l string) string {
// 	ss := strings.Split(l, "sea.Load")
// 	l = ss[1]
// 	ss = strings.Split(l, `"`)
// 	return ss[1]
// }

// func getType(l string) string {

// }

// func writeOut(path string, vs []v) error {

// }
