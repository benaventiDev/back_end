package main

func FileReport(output string, disk_path string) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	if disk_path != "/users.txt" {
		error_list = append(error_list, "Not found")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	//PrintPartitions(mbr, disk_path)

	return InstResponse{Iserror: false, Errors: error_list, Result: Users}
}

func BuildFileRep(sb Superbloque, path_file string) string {
	ret_str := ""
	/* Lectura del superbloque */
	if ret_str == "" {
		return ""
	}
	inode := MakeForceFolderFile(path_file)

	for i := 0; i < 15; i++ {
		if inode.I_block[i] != -1 {
			fb := FolderBlock{}
			fb, _ = ReadFolderBlock(LoggedInPartition.Path, inode.I_block[0])
			for j := 0; j < 4; j++ {
				if fb.B_content[j].B_inodo != -1 {
					if Byte12ToString(fb.B_content[j].B_name) == "" {
						ret_str += Byte12ToString(fb.B_content[j].B_name) + " -> "
					}
				}
			}
		}
	}
	return ret_str
}
