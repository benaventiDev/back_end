package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

func DiskReport(output string, disk_path string) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	mbr := ReadMBR(disk_path)
	//PrintPartitions(mbr, disk_path)
	output_str := "digraph G {\ngraph[bgcolor=\"#141D26\" margin=0]\nrankdir=\"TB\";\nnode [shape=plaintext fontname= \"Ubuntu\" fontsize=\"14\"];\nedge [style=\"invis\"];\n\n"
	output_str += PrintPartitions(mbr, disk_path)
	WriteDotRaw("output/output.dot", output_str)
	graph := GenerateGraph(output)
	if graph == nil {
		error_list = append(error_list, "Error generating graph")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	return InstResponse{Iserror: false, Errors: error_list, Result: "rep", Img: graph}
}

func PrintPartitions(mbr MBR, disk_path string) string {
	colors := []string{"#229954", "#2471a3", "#5d6d7e", "#8e44ad", "#b7950b"}
	str_ret := "\"Disk report\" [ label = <" + "<TABLE CELLBORDER=\"2\" BGCOLOR=\"BLACK\" BORDER=\"2\" COLOR=\"BLACK\"  CELLSPACING=\"0\">"
	str_ret += "\n\n"
	str_ret += "<TR>"
	str_ret += "<TD BGCOLOR=\"#1c2833\" COLSPAN=\"100\">"
	str_ret += "<FONT POINT-SIZE=\"20\" COLOR=\"#f2f3f4\">Disk Report of: \t"
	str_ret += disk_path
	str_ret += "</FONT>"
	str_ret += "</TD>"
	str_ret += "</TR>"
	str_ret += "<TR>"
	str_ret += "\n\n"
	str_ret += "<TD HEIGHT=\"150\" WIDTH=\"75\" BGCOLOR=\"#c0392b\">MBR</TD>"
	str_ret += "\n\n"

	var last_start int32 = int32(unsafe.Sizeof(mbr))
	for i := 0; i < 4; i++ {
		part := mbr.Mbr_partition[i]
		if part.Part_status == 0 {
			continue
		}

		if part.Part_start-int32(last_start) > 0 {
			str_ret += "<TD HEIGHT=\"160\" WIDTH=\"150"
			str_ret += "\" BGCOLOR=\"#b7950b"
			str_ret += "\">FREE<br/><br/><b>"
			str_ret += GetPorcentaje(part.Part_start-int32(last_start), mbr.Mbr_tamano)
			str_ret += "% of the disk.</b></TD>\n"
		}
		last_start = part.Part_start + part.Part_size

		str_ret += "<TD HEIGHT=\"160\" WIDTH=\""
		str_ret += "150\" BGCOLOR=\""
		if part.Partition_type == Primary {
			str_ret += colors[0]
		} else if part.Partition_type == Extended {
			str_ret += colors[2]
		}
		str_ret += "\">"

		if part.Partition_type == Extended {
			str_ret += "<TABLE ALIGN=\"LEFT\">"
			str_ret += "<TR><TD>"
		}

		str_ret += Byte64ToString(part.Part_name)
		str_ret += "<br/>"
		str_ret += "<br/>"
		if part.Partition_type == Primary {
			str_ret += "Primary"
		} else if part.Partition_type == Extended {
			str_ret += "Extended"
		}
		str_ret += "<br/>"
		str_ret += "<br/>"

		part_size := part.Part_size
		if part_size < 1000 {
			str_ret += fmt.Sprint(part_size)
			str_ret += " Bytes"
		} else if part_size < 999999 {
			part_size = part_size / 1000
			str_ret += "<b>"
			str_ret += fmt.Sprint(part_size)
			str_ret += " KB</b>"
		} else {
			part_size = part_size / 1000000
			str_ret += "<b>"
			str_ret += fmt.Sprint(part_size)
			str_ret += " MB</b>"
		}
		str_ret += "<br/>"
		str_ret += "<b>"
		str_ret += GetPorcentaje(part.Part_size, mbr.Mbr_tamano)
		str_ret += "% of the disk.</b>"
		str_ret += "<br/>"
		str_ret += "</TD>"
		str_ret += "\n"

		if part.Partition_type == Extended {
			str_ret += "</TR><TR>\n<TD>\n<TABLE BORDER=\"0\">\n<TR>\n"
			ebr, _ := GetEBR(disk_path, part.Part_start)
			if ebr.Part_size != 0 {
				str_ret += PrintLogicPartitions(ebr, mbr.Mbr_tamano, disk_path, colors[i], part.Part_start+part.Part_size)
			}

			str_ret += "</TR>\n</TABLE>\n</TD>\n\n</TR>\n</TABLE>\n</TD>"
		}
	}
	if mbr.Mbr_tamano-last_start > 0 {
		str_ret += "<TD HEIGHT=\"160\" WIDTH=\"150"
		str_ret += "\" BGCOLOR=\"#b7950b"
		str_ret += "\">FREE<br/><br/><b>"
		str_ret += GetPorcentaje(mbr.Mbr_tamano-last_start, mbr.Mbr_tamano)

		str_ret += "% of the disk.</b></TD>\n"
	}
	str_ret += "\n</TR>\n\n</TABLE>>];\n}"
	return str_ret
}

func PrintLogicPartitions(ebr EBR, mbr_tamano int32, disk_path string, _color string, extended_limit int32) string {
	init_position := ebr.Part_start
	str_ret := ""
	for ebr.Part_next != -1 {
		if ebr.Part_start-init_position > 0 {
			str_ret += "<TD><TABLE BORDER=\"1\"><TR>\n<TD HEIGHT=\"160\" WIDTH=\"150"
			str_ret += "\" BGCOLOR=\"#b7950b"
			str_ret += "\">FREE<br/><br/><b>"
			str_ret += GetPorcentaje(ebr.Part_start-init_position, mbr_tamano)
			str_ret += "% of the disk.</b></TD>\n</TR></TABLE></TD>"
		}
		if ebr.Part_size != 0 {
			init_position = ebr.Part_start + ebr.Part_size
			str_ret += "<TD><TABLE BORDER=\"1\"><TR>\n<TD HEIGHT=\"160\" WIDTH=\"40\" BGCOLOR=\""
			str_ret += _color
			str_ret += "\">"
			str_ret += "EBR"
			str_ret += "</TD>"
			str_ret += "<TD HEIGHT=\"160\" WIDTH=\"150\" BGCOLOR=\""
			str_ret += _color
			str_ret += "\">"
			str_ret += "<br/>"
			str_ret += Byte64ToString(ebr.Part_name)
			str_ret += "<br/>"
			str_ret += "<br/>"
			str_ret += "Lógica"
			str_ret += "<br/>"
			str_ret += "<br/>"

			if ebr.Part_size < 1000 {
				str_ret += fmt.Sprint(ebr.Part_size)
				str_ret += " Bytes"
			} else if ebr.Part_size < 999999 {
				str_ret += fmt.Sprint(ebr.Part_size)
				str_ret += " KB"
			} else {
				str_ret += fmt.Sprint(ebr.Part_size / 1000000)
				str_ret += " MB"
			}
			str_ret += "<br/>"
			str_ret += "<b>"

			str_ret += GetPorcentaje(ebr.Part_size, mbr_tamano)

			str_ret += "% of the disk</b>\n<br/>\n</TD></TR></TABLE></TD>\n\n"
		}

		ebr, _ = GetEBR(disk_path, ebr.Part_next)

	}

	if ebr.Part_start-init_position > 0 {
		str_ret += "<TD><TABLE BORDER=\"1\"><TR>\n<TD HEIGHT=\"160\" WIDTH=\"150"
		str_ret += "\" BGCOLOR=\"#b7950b"
		str_ret += "\">FREE<br/><br/><b>"
		str_ret += GetPorcentaje(ebr.Part_start-init_position, mbr_tamano)
		str_ret += "% of the disk.</b></TD>\n</TR></TABLE></TD>"
	}

	if ebr.Part_size != 0 {
		init_position = ebr.Part_start + ebr.Part_size
		str_ret += "<TD><TABLE BORDER=\"1\"><TR>\n<TD HEIGHT=\"160\" WIDTH=\"40\" BGCOLOR=\""
		str_ret += _color
		str_ret += "\">"
		str_ret += "EBR"
		str_ret += "</TD>"
		str_ret += "<TD HEIGHT=\"160\" WIDTH=\"120\" BGCOLOR=\""
		str_ret += _color
		str_ret += "\">"
		str_ret += Byte64ToString(ebr.Part_name)
		str_ret += "<br/>"
		str_ret += "<br/>"
		str_ret += "Lógica"
		str_ret += "<br/>"
		str_ret += "<br/>"
		str_ret += "<b>"

		str_ret += GetPorcentaje(ebr.Part_size, mbr_tamano)

		str_ret += "% of the disk</b>\n<br/>\n</TD></TR></TABLE></TD>\n\n"

	}

	if extended_limit-init_position > 0 {
		str_ret += "<TD><TABLE BORDER=\"1\"><TR>\n<TD HEIGHT=\"160\" WIDTH=\"150"
		str_ret += "\" BGCOLOR=\"#b7950b"
		str_ret += "\">FREE<br/><br/><b>"
		str_ret += GetPorcentaje(extended_limit-init_position, mbr_tamano)

		str_ret += "% of the disk.</b></TD>\n</TR></TABLE></TD>"
	}
	return str_ret

}

func WriteDotRaw(dot_file string, content string) {
	CreatePath(dot_file)
	f, err := os.Create(dot_file)

	if err != nil {
		fmt.Printf("WriteDotRaw %v", err)
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		fmt.Printf("WriteDotRaw %v", err2)
	}
}

func PrintPartitionsRaw(mbr MBR, disk_path string) {
	fmt.Println("Partitions:")
	for i := 0; i < 4; i++ {
		fmt.Print("Partition ", i)
		fmt.Print(". Status: ", mbr.Mbr_partition[i].Part_status)
		if mbr.Mbr_partition[i].Partition_type == Extended {
			fmt.Print(". Type: ", "Extendida")
		} else if mbr.Mbr_partition[i].Partition_type == Primary {
			fmt.Print(". Type: ", "Primaria")
		}

		fmt.Print(". Fit: ", mbr.Mbr_partition[i].Fit_type)
		fmt.Print(". Start: ", mbr.Mbr_partition[i].Part_start)
		fmt.Print(". Size: ", mbr.Mbr_partition[i].Part_size)
		fmt.Print(". Name: ", mbr.Mbr_partition[i].Part_name)
		fmt.Println()
		if mbr.Mbr_partition[i].Partition_type == Extended {
			fmt.Println("***************Logicas***************")
			ebr, err_1 := GetEBR(disk_path, mbr.Mbr_partition[i].Part_start)
			if err_1 == "" {
				if ebr.Part_size != 0 {

					for ebr.Part_next != -1 {
						PrintEBRRaw(ebr)
						ebr, err_1 = GetEBR(disk_path, ebr.Part_next)
						if err_1 != "" {
							fmt.Printf("Error2 reading ebr: %s", err_1)
							break
						}
					}
					PrintEBRRaw(ebr)
				}
			} else {
				fmt.Printf("Error1 reading ebr: %s", err_1)
			}

			fmt.Println("*************************************")

		}
	}

}

func GetPorcentaje(_free int32, _mbr_tamano int32) string {
	tmp := (float64)(_free) / (float64)(_mbr_tamano) * 100.0
	//strconv.Itoa(part_size)
	return strconv.Itoa(int(tmp))
}

func PrintEBRRaw(ebr EBR) {
	fmt.Print("\tLogic Partition ** ")
	fmt.Print(". Status: ", ebr.Part_status)
	fmt.Print(". Type: ", "Logica")
	fmt.Print(". Fit: ", ebr.Part_fit)
	fmt.Print(". Start: ", ebr.Part_start)
	fmt.Print(". Size: ", ebr.Part_size)
	fmt.Print(". Name: ", ebr.Part_name)
	fmt.Println()
}
