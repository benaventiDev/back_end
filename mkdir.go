package main

import (
	"fmt"
	"strings"
	"unsafe"
)

func Mkdir_(inst Instruction) InstResponse {
	/*error_list := []string{}
	var inst_res InstResponse



	path_val, path_err := CheckPathExists(inst)
	if path_err != "" {
		error_list = append(error_list, path_err)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	fmt.Printf("%v", path_val)

	inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}*/
	return MakeForceFolder(inst)

}

func MakeForceFolder(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse

	path_raw := Get_param_value(inst.Params, Path_)
	if path_raw == nil {
		error_list = append(error_list, "Invalid path value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	path_val := fmt.Sprint(path_raw)

	/*
		if(logged_in_user == NULL){
			show_result_operation("You must be logged in to perform this operation.", failure);
			return -1;
		}*/
	if LoggedInPartition.Part_type == Error_Part_type {
		error_list = append(error_list, "No partition mounted.")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	var part_start int32
	//var part_size int32
	if LoggedInPartition.Part_type == Primary {
		part_start = LoggedInPartition.Part.Part_start
		//part_size = LoggedInPartition.Part.Part_size
	} else if LoggedInPartition.Part_type == Logic {
		part_start = LoggedInPartition.Logica.Part_start
		//part_size = LoggedInPartition.Logica.Part_size
	}

	/*
		folder_to_create := FolderBlock{} // Nuevo bloque carpeta
		folder_content := Content{}       // Nuevo bloque contenido
		new_inode := Inode{}              // Nuevo inodo
	*/
	/* Lectura del superbloque */

	sb, err := ReadSB(LoggedInPartition.Path, part_start)
	if err != "" {
		error_list = append(error_list, "Error reading sb")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	/*
		inode_root, err := ReadInode(LoggedInPartition.Path, sb.S_inode_start)
		if err != "" {
			error_list = append(error_list, "Error reading Inode root")
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}

		inode := inode_root*/
	var current_inode int32 = 0
	folderBlock := FolderBlock{}
	paths := strings.Split(path_val, "/")
	paths = paths[1:]
	inode, err_1 := ReadInode(LoggedInPartition.Path, sb.S_inode_start+current_inode*int32(unsafe.Sizeof(folderBlock)))
	if err_1 != "" {
		error_list = append(error_list, "Error reading Inode root")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	for _, path := range paths {

		found := false
		for i := 0; i < 15; i++ {
			if inode.I_block[i] != -1 {
				folderBlock, err = ReadFolderBlock(LoggedInPartition.Path, sb.S_block_start+int32(inode.I_block[i])*int32(unsafe.Sizeof(folderBlock)))
				if err != "" {
					error_list = append(error_list, "Error reading Inode root")
					inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
					return inst_res
				}
				for j := 0; j < 4; j++ {

					if Byte12ToString(folderBlock.B_content[j].B_name) == path {
						current_inode = folderBlock.B_content[j].B_inodo
						inode, err = ReadInode(LoggedInPartition.Path, sb.S_inode_start+int32(folderBlock.B_content[j].B_inodo)*int32(sb.S_inode_size))
						if err != "" {
							error_list = append(error_list, "Error reading Inode root")
							inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
							return inst_res
						}
						found = true
						break
					}
				}
				if found {
					if inode.Inode_type == Inode_file {
						error_list = append(error_list, "Error: "+path+"is a file, part from path: "+path_val)
						inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
						return inst_res
					} else {
						break
					}
				}
			}
		}

		if !found {

			created := false
			for i := 0; i < 15; i++ {
				if created {
					break
				}
				if inode.I_block[i] != -1 {
					folderBlock, err = ReadFolderBlock(LoggedInPartition.Path, sb.S_block_start+int32(inode.I_block[i])*int32(sb.S_block_size))
					if err != "" {
						error_list = append(error_list, "Error reading Inode root")
						inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
						return inst_res
					}
					for j := 0; j < 4; j++ {
						if folderBlock.B_content[j].B_inodo == -1 {
							//crear un nuevo bloque
							newFolderBlock := FolderBlock{}
							for j := 0; j < 4; j++ {
								newFolderBlock.B_content[j].B_inodo = -1
							}
							newFolderBlock.B_content[0].B_inodo = current_inode
							newFolderBlock.B_content[0].B_name = StringToByte12(".")
							newFolderBlock.B_content[1].B_inodo = sb.S_first_ino
							newFolderBlock.B_content[1].B_name = StringToByte12("..")
							WriteStruct(LoggedInPartition.Path, sb.S_block_start+int32(sb.S_first_blo)*int32(sb.S_block_size), newFolderBlock)
							current_inode = sb.S_first_ino
							//Create new inode here
							inode_folder := Inode{}
							for i := 0; i < 15; i++ {
								inode_folder.I_block[i] = -1
							}
							inode_folder.I_block[0] = sb.S_first_blo
							inode_folder.I_size = 0
							inode_folder.Inode_type = Inode_folder
							inode_folder.I_gid = 1
							inode_folder.I_uid = 1
							inode_folder.I_perm = 777
							//guardar inodo
							WriteStruct(LoggedInPartition.Path, sb.S_inode_start+int32(sb.S_first_ino)*int32(sb.S_inode_size), inode_folder)

							//actualizar bloque actual
							folderBlock.B_content[j].B_inodo = sb.S_first_ino
							folderBlock.B_content[j].B_name = StringToByte12(path)
							WriteStruct(LoggedInPartition.Path, sb.S_block_start+int32(inode.I_block[i])*int32(sb.S_block_size), folderBlock)
							//Actualizar Bitmaps
							bitmapBlock := ReadBitmapBlock(LoggedInPartition.Path, sb.S_bm_block_start, sb.S_blocks_count)
							bitmapInode := ReadBitmapInode(LoggedInPartition.Path, sb.S_bm_inode_start, sb.S_inodes_count)
							bitmapBlock[sb.S_first_blo] = 1
							bitmapInode[sb.S_first_ino] = 1
							WriteStruct(LoggedInPartition.Path, sb.S_bm_inode_start, bitmapInode)
							WriteStruct(LoggedInPartition.Path, sb.S_bm_block_start, bitmapBlock)
							var next_free_inode int32 = -1
							for i := 0; i < len(bitmapInode); i++ {
								if bitmapInode[i] == 0 {
									next_free_inode = int32(i)
									break
								}
							}

							var next_free_block int32 = -1
							for i := 0; i < len(bitmapBlock); i++ {
								if bitmapBlock[i] == 0 {
									next_free_block = int32(i)
									break
								}
							}

							sb.S_first_blo = next_free_block
							sb.S_first_ino = next_free_inode
							sb.S_free_blocks_count--
							sb.S_free_inodes_count--
							WriteStruct(LoggedInPartition.Path, part_start, sb)
							inode = inode_folder
							created = true
							break
						}

					}
				} else { // Crear un nuevo bloque carpeta

					//Crear nuevo bloque carpeta
					bitmapBlock := ReadBitmapBlock(LoggedInPartition.Path, sb.S_bm_block_start, sb.S_blocks_count)
					bitmapInode := ReadBitmapInode(LoggedInPartition.Path, sb.S_bm_inode_start, sb.S_inodes_count)

					folderBlock := FolderBlock{}
					for j := 0; j < 4; j++ {
						folderBlock.B_content[j].B_inodo = -1
					}
					inode.I_block[i] = sb.S_first_blo
					WriteStruct(LoggedInPartition.Path, sb.S_inode_start+int32(current_inode)*int32(sb.S_inode_size), inode)
					//escribir un bloque en direccion primer bnloque libre que apunte a direccion primer inodo libre
					newFolderBlock := FolderBlock{}
					for j := 0; j < 4; j++ {
						newFolderBlock.B_content[j].B_inodo = -1
					}
					newFolderBlock.B_content[0].B_inodo = sb.S_first_ino
					newFolderBlock.B_content[0].B_name = StringToByte12(path)
					WriteStruct(LoggedInPartition.Path, sb.S_block_start+int32(sb.S_first_blo)*int32(sb.S_block_size), newFolderBlock)
					bitmapBlock[sb.S_first_blo] = 1
					bitmapInode[sb.S_first_ino] = 1

					//buscar nuevo blocque libre
					var next_free_block int32 = -1
					for i := 0; i < len(bitmapBlock); i++ {
						if bitmapBlock[i] == 0 {
							next_free_block = int32(i)
							bitmapBlock[i] = 1
							break
						}
					}
					//crerar un bloque para el nuevo inodo
					newFolderBlock = FolderBlock{}
					for j := 0; j < 4; j++ {
						newFolderBlock.B_content[j].B_inodo = -1
					}
					newFolderBlock.B_content[0].B_inodo = current_inode
					newFolderBlock.B_content[0].B_name = StringToByte12(".")
					newFolderBlock.B_content[0].B_inodo = sb.S_first_ino
					newFolderBlock.B_content[0].B_name = StringToByte12("..")
					WriteStruct(LoggedInPartition.Path, sb.S_block_start+int32(next_free_block)*int32(sb.S_block_size), newFolderBlock)
					//crear un nuevo  inodo folder con el pimer libre y asignarle bloque recien creado
					inode_folder := Inode{}
					for i := 0; i < 15; i++ {
						inode_folder.I_block[i] = -1
					}
					inode_folder.I_block[0] = next_free_block
					inode_folder.I_size = 0
					inode_folder.Inode_type = Inode_folder
					inode_folder.I_gid = 1
					inode_folder.I_uid = 1
					inode_folder.I_perm = 777
					//guardar inodo
					WriteStruct(LoggedInPartition.Path, sb.S_inode_start+int32(sb.S_first_ino)*int32(sb.S_inode_size), inode_folder)
					current_inode = sb.S_first_ino
					inode = inode_folder
					//actualizar y guardar sb: count, free start, free count,

					for i := 0; i < len(bitmapBlock); i++ {
						if bitmapBlock[i] == 0 {
							next_free_block = int32(i)
							break
						}
					}

					var next_free_inode int32 = -1
					for i := 0; i < len(bitmapInode); i++ {
						if bitmapInode[i] == 0 {
							next_free_inode = int32(i)
							break
						}
					}

					sb.S_first_blo = next_free_block
					sb.S_first_ino = next_free_inode
					sb.S_free_blocks_count--
					sb.S_free_inodes_count--
					WriteStruct(LoggedInPartition.Path, part_start, sb)
					WriteStruct(LoggedInPartition.Path, sb.S_bm_inode_start, bitmapInode)
					WriteStruct(LoggedInPartition.Path, sb.S_bm_block_start, bitmapBlock)
					created = true
					break

				}
				if created {
					break
				}
			}
		}

	}

	inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Empty result"}
	return inst_res

}

//int make_folder_force(char* path, int print_success_message){

/* CREACIÓN DE CARPETA */
/*
	if(logged_in_user == NULL){
		show_result_operation("You must be logged in to perform this operation.", failure);
		return -1;
	}
	if(logged_in_partition == NULL){
		show_result_operation("No partition mounted.", failure);
		return -1;
	}
	if(path[0]!= '/'){
		show_result_operation("Path must be absolute.", failure);
		return -1;
	}*/

//}
