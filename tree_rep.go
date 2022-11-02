package main

import (
	"fmt"
	"strings"
	"unsafe"
)

func TreeReport(output string, mounted Mounted) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	output_str := "digraph G {\ngraph[bgcolor=\"#ffffff\" margin=0]\nrankdir=\"LR\";\n"
	output_str += "node [shape=plaintext fontname= \"Ubuntu\"];\nedge [arrowhead=\"normal\" penwidth=3];\n\n"

	var part_start int32
	if mounted.Part_type == Primary {
		part_start = mounted.Part.Part_start
	} else if mounted.Part_type == Logic {
		part_start = mounted.Logica.Part_start
	}

	sb, err := ReadSB(mounted.Path, part_start)
	if err != "" {
		error_list = append(error_list, "Error generating graph")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	output_str += PrintTree(sb, mounted.Path, 0)
	output_str += "}\n"

	WriteDotRaw("output/output.dot", output_str)
	graph := GenerateGraph(output)
	if graph == nil {
		error_list = append(error_list, "Error generating graph")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	return InstResponse{Iserror: false, Errors: error_list, Result: "rep", Img: graph}

}

func PrintTree(sb Superbloque, path string, inode_index int32) string {
	fmt.Printf("PrintTree inode_index:%v\n", inode_index)
	inode := Inode{}
	inode, err := ReadInode(path, sb.S_inode_start+inode_index*int32(unsafe.Sizeof(inode)))
	if err != "" {
		fmt.Println("PrintTree: Error reading inode")
	}

	str_ret := "\"INODE_"
	str_ret += fmt.Sprint(inode_index)
	str_ret += "\" [ fontsize=\"17\"  label = <\n"
	str_ret += "<TABLE BGCOLOR=\"#009999\" BORDER=\"2\" COLOR=\"BLACK\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	str_ret += "<TR>\n"
	str_ret += "<TD BGCOLOR=\"#B8860B\" COLSPAN=\"2\">Inodo "
	str_ret += fmt.Sprint(inode_index)
	str_ret += " </TD>\n"
	str_ret += "</TR>\n"
	str_ret += "<TR>\n"
	str_ret += "<TD WIDTH=\"130\" BGCOLOR=\"#708090\"><B>Tipo</B></TD>\n"
	str_ret += "<TD WIDTH=\"70\" BGCOLOR=\"#797d7f\"> "
	if inode.Inode_type == Inode_folder {
		str_ret += "Folder"
	} else if inode.Inode_type == Inode_file {
		str_ret += "File"
	} else {
		str_ret += "Error"
	}

	str_ret += " </TD>\n"
	str_ret += "</TR>\n"
	str_ret += "<TR>\n"
	str_ret += "<TD WIDTH=\"130\" BGCOLOR=\"#708090\"><B>Tamaño</B></TD>\n"
	str_ret += "<TD WIDTH=\"70\" BGCOLOR=\"#797d7f\"> "
	str_ret += fmt.Sprint(inode.I_size)
	str_ret += " </TD>\n"
	str_ret += "</TR>\n"
	str_ret += "<TR>\n"
	str_ret += "<TD WIDTH=\"130\" BGCOLOR=\"#708090\"><B>Permisos</B></TD>\n"
	str_ret += "<TD WIDTH=\"70\" BGCOLOR=\"#797d7f\"> "

	str_ret += fmt.Sprint(inode.I_perm)
	str_ret += " </TD>\n"
	str_ret += "</TR>\n\n"

	for i := 0; i < 15; i++ {

		str_ret += "<TR>\n"
		str_ret += "<TD WIDTH=\"130\" BGCOLOR=\"#708090\"><B>i_block["
		str_ret += fmt.Sprint(i)
		str_ret += "]</B></TD>\n"
		str_ret += "<TD PORT=\"PI_"
		str_ret += fmt.Sprint(i)
		str_ret += "\" BGCOLOR=\"#797d7f\">"
		str_ret += fmt.Sprint(inode.I_block[i])
		str_ret += "</TD>\n"
		str_ret += "</TR>\n\n"
	}
	str_ret += "</TABLE>>];\n\n"

	for i := 0; i < 15; i++ {
		if inode.I_block[i] != -1 {
			str_ret += "\"INODE_"
			str_ret += fmt.Sprint(inode_index)
			str_ret += "\":\"PI_"
			str_ret += fmt.Sprint(i)
			str_ret += "\" -> \"BLOCK_"
			str_ret += fmt.Sprint(inode.I_block[i])
			str_ret += "\";\n\n"
		}
	}

	if inode.Inode_type == Inode_folder {
		for i := 0; i < 15; i++ {
			if inode.I_block[i] != -1 {
				str_ret += Print_folder_block(path, sb, inode.I_block[i], inode_index)
			}
		}

	} else if inode.Inode_type == Inode_file {
		for i := 0; i < 15; i++ {
			if inode.I_block[i] != -1 {
				str_ret += Print_file_block(path, sb, inode.I_block[i])
			}
		}

	} else {
		str_ret += "Error"
	}
	return str_ret
}

func Print_folder_block(path string, sb Superbloque, block_index int32, inode_index int32) string {
	fmt.Printf("Print_folder_block\n")
	var fb FolderBlock
	fb, err := ReadFolderBlock(path, sb.S_block_start+block_index*int32(unsafe.Sizeof(fb)))
	if err != "" {
		fmt.Println("Print_folder_block: Error reading FolderBlock")
	}
	str_ret := "\"BLOCK_"
	str_ret += fmt.Sprint(block_index)
	str_ret += "\" [ fontsize=\"17\" label = <\n"
	str_ret += "<TABLE BGCOLOR=\"#009999\" BORDER=\"2\" COLOR=\"BLACK\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	str_ret += "<TR>\n"
	str_ret += "<TD BGCOLOR=\"#B8860B\" COLSPAN=\"2\">Bloque de carpeta: "
	str_ret += fmt.Sprint(block_index)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n"
	str_ret += "<TR>\n"
	str_ret += "<TD WIDTH=\"130\" BGCOLOR=\"#708090\"><B>b_name</B></TD>\n"
	str_ret += "<TD WIDTH=\"70\" BGCOLOR=\"#708090\"><B>b_inodo</B></TD>\n"
	str_ret += "</TR>\n\n"

	for i := 0; i < 4; i++ {
		str_ret += "<TR>\n"
		str_ret += "<TD "

		if fb.B_content[i].B_inodo == inode_index { //&& Byte12ToString(fb.B_content[i].B_name) != "." && Byte12ToString(fb.B_content[i].B_name) != ".." {
			str_ret += "PORT=\"PB_"
			str_ret += fmt.Sprint(i)
			str_ret += "\" "
		}
		str_ret += "ALIGN=\"left\">   "
		str_ret += Byte12ToString(fb.B_content[i].B_name)
		str_ret += "</TD>\n"
		str_ret += "<TD "

		if fb.B_content[i].B_inodo != inode_index { //&& Byte12ToString(fb.B_content[i].B_name) != "." && Byte12ToString(fb.B_content[i].B_name) != ".." {
			str_ret += "PORT=\"PB_"
			str_ret += fmt.Sprint(i)
			str_ret += "\" "
		}

		str_ret += ">"
		str_ret += fmt.Sprint(fb.B_content[i].B_inodo)
		str_ret += "</TD>\n"
		str_ret += "</TR>\n\n"
	}

	str_ret += "</TABLE>>];\n\n"

	for i := 0; i < 4; i++ {
		if fb.B_content[i].B_inodo != -1 && Byte12ToString(fb.B_content[i].B_name) != "." && Byte12ToString(fb.B_content[i].B_name) != ".." {
			str_ret += "\"BLOCK_"
			str_ret += fmt.Sprint(block_index)
			str_ret += "\":\"PB_"
			str_ret += fmt.Sprint(i)
			str_ret += "\" -> \"INODE_"
			str_ret += fmt.Sprint(fb.B_content[i].B_inodo)
			str_ret += "\";\n\n"
		}
	}

	for i := 0; i < 4; i++ {
		if fb.B_content[i].B_inodo != -1 && Byte12ToString(fb.B_content[i].B_name) != "." && Byte12ToString(fb.B_content[i].B_name) != ".." {
			str_ret += PrintTree(sb, path, fb.B_content[i].B_inodo)

		}
	}
	return str_ret

}

func Print_file_block(path string, sb Superbloque, block_index int32) string {
	fmt.Printf("Print_file_block\n")
	var fb FilesBlock
	fb, err := ReadFileBlock(path, sb.S_block_start+block_index*int32(unsafe.Sizeof(fb)))
	if err != "" {
		fmt.Println("Print_file_block: Error reading FolderBlock")
	}

	str_ret := "\"BLOCK_"
	str_ret += fmt.Sprint(block_index)
	str_ret += "\" [ fontsize=\"17\" label = <\n"
	str_ret += "<TABLE BGCOLOR=\"#009999\"  BORDER=\"2\" COLOR=\"BLACK\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	str_ret += "<TR>\n"
	str_ret += "<TD WIDTH=\"190\" BGCOLOR=\"#708090\">Bloque de archivo "
	str_ret += fmt.Sprint(block_index)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"
	str_ret += "<TR>\n"
	str_ret += "<TD>\n"

	aux_str := Byte64ToString(fb.B_content)
	aux_str = strings.ReplaceAll(aux_str, "\n", "<BR ALIGN=\"LEFT\"/>\n")
	str_ret += aux_str

	str_ret += "\n"
	str_ret += "</TD>\n</TR>\n\n</TABLE>>];\n\n"
	return str_ret

}
