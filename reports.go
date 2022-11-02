package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-graphviz"
)

func Rep_(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse

	name_raw := Get_param_value(inst.Params, Name_)
	if name_raw == nil {
		error_list = append(error_list, "Invalid name value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	name_val := fmt.Sprint(name_raw)

	id_raw := Get_param_value(inst.Params, Id_)
	if id_raw == nil {
		error_list = append(error_list, "Invalid Id value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	id_val := fmt.Sprint(id_raw)

	path_raw := Get_param_value(inst.Params, Path_)
	if path_raw == nil {
		error_list = append(error_list, "Invalid path parameter ")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	path_val := fmt.Sprint(path_raw)
	dir, _ := filepath.Split(path_val)
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err_1 := os.MkdirAll(dir, os.ModePerm)
			if err_1 != nil {
				error_list = append(error_list, "Error creating directory: "+dir+": "+err_1.Error())
				inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
				return inst_res
			}
		}
	}

	var mounted Mounted

	for _, mounted_part := range MountedParts {
		if mounted_part.FullID == id_val {
			mounted = mounted_part
		}
	}
	if mounted.FullID == "" {
		error_list = append(error_list, "Partition is not mounted")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	switch strings.ToLower(name_val) {
	case "disk":
		return DiskReport(path_val, mounted.Path)
	case "tree":
		return TreeReport(path_val, mounted)
	case "file":
		return FileReport(path_val, path_val)
	case "sb":
		return SbReport(path_val, mounted)
	default:
		error_list = append(error_list, "Incorrect report inst: "+name_val)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	//fmt.Printf("%v", path_val)
	//fmt.Printf("%v", name_val)
	//fmt.Printf("%v", id_val)

	/**
	ruta_raw := Get_param_value(inst.Params, Ruta_)
	if ruta_raw == nil {
		error_list = append(error_list, "Invalid Ruta value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	ruta_val := fmt.Sprint(ruta_raw)

	*/

	inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
	return inst_res

}

func GenerateGraph(output string) []byte {

	path := "output/output.dot"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Graphviz Error:", err)
		return nil
	}
	graph, err := graphviz.ParseBytes(b)
	if err != nil {
		fmt.Println("Graphviz Error :", err)
	}
	// create your graph
	g := graphviz.New()

	dir, _ := filepath.Split(output)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err_1 := os.MkdirAll(dir, os.ModePerm)
			if err_1 != nil {
				fmt.Println("Graphviz Error 3:", err)
				return nil
			}
		}
	}

	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Graphviz Error 3:", err)
		return nil
	}
	file.Close()

	if err := g.RenderFilename(graph, graphviz.PNG, output); err != nil {
		fmt.Println("Graphviz Error 4:", err)
	}

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		fmt.Println("Graphviz Error 5:", err)
	}

	var bytebuf []byte = buf.Bytes()
	return bytebuf
}
