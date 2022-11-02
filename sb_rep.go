package main

import (
	"fmt"
)

func SbReport(output string, mounted Mounted) InstResponse {

	error_list := []string{}
	var inst_res InstResponse
	output_str := BuildSBRep(mounted)

	WriteDotRaw("output/output.dot", output_str)
	graph := GenerateGraph(output)
	if graph == nil {
		error_list = append(error_list, "Error generating graph")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	return InstResponse{Iserror: false, Errors: error_list, Result: "rep", Img: graph}
}

func BuildSBRep(mounted Mounted) string {

	str_ret := "digraph G {\n"
	str_ret += "graph[bgcolor=\"#141D26\" margin=0]\n"
	str_ret += "rankdir=\"TB\";\n"
	str_ret += "node [shape=plaintext fontname= \"Ubuntu\" fontsize=\"14\"];\n"
	str_ret += "edge [style=\"invis\"];\n\n"

	// Lectura del superbloque
	var part_start int32
	if mounted.Part_type == Primary {
		part_start = mounted.Part.Part_start
	} else if mounted.Part_type == Logic {
		part_start = mounted.Logica.Part_start
	} else {
		fmt.Println("bm_report::build_bm_block_report:: incorrect parrtiton type.")
	}
	sb, err := ReadSB(mounted.Path, part_start)
	if err != "" {
		fmt.Println("error reading sb.")
	}

	str_ret += "\"SB Report\" [ margin=\"1\" label = <\n"
	str_ret += "<TABLE BGCOLOR=\"#009999\" BORDER=\"2\" COLOR=\"BLACK\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
	str_ret += "<TR>\n"
	str_ret += "<TD BGCOLOR=\"#d03939\" COLSPAN=\"2\">REPORTE DE SUPERBLOQUE</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD WIDTH=\"150\" BGCOLOR=\"#ff4660\"><B>Nombre</B></TD>\n"
	str_ret += "<TD BGCOLOR=\"#ff4660\"><B>Valor</B></TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_filesystem_type</TD>\n"
	str_ret += "<TD>"
	str_ret += "Ext2"
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_inodes_count</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_inodes_count)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_blocks_count</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_blocks_count)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_free_inodes_count</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_free_inodes_count)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_free_blocks_count</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_free_blocks_count)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_mtime</TD>\n"
	str_ret += "<TD>"
	str_ret += mounted.Tmounted
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_umtime</TD>\n"
	str_ret += "<TD>00:00:00"
	//str_ret += ctime(&sb.s_umtime)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_mnt_count</TD>\n"
	str_ret += "<TD>"
	str_ret += "1" // fmt.Sprint(sb.S_mnt_count)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_magic</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_magic)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_inode_size</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_inode_size)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_block_size</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_block_size)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_first_ino</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_first_ino)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_first_blo</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_first_blo)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_bm_inode_start</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_bm_inode_start)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_bm_block_start</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_bm_block_start)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_inode_start</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_inode_start)
	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"

	str_ret += "<TR>\n"
	str_ret += "<TD ALIGN=\"left\">s_block_start</TD>\n"
	str_ret += "<TD>"
	str_ret += fmt.Sprint(sb.S_block_start)

	str_ret += "</TD>\n"
	str_ret += "</TR>\n\n"
	str_ret += "</TABLE>>];\n\n}"

	return str_ret
}
