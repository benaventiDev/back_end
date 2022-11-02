package main

import (
	"fmt"
	"strings"
	"unsafe"
)

var Users string = "1, G, root\n1, U, root, root, 123\n"

func Mkfs_(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	id_raw := Get_param_value(inst.Params, Id_)
	if id_raw == nil {
		error_list = append(error_list, "Invalid id value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	type_raw := Get_param_value(inst.Params, Type_)
	var type_str string
	if type_raw != nil {
		type_str = strings.ToLower(fmt.Sprint(type_raw))
		if type_str != "full" {
			error_list = append(error_list, "Invalid type value")
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
	}

	id_val := fmt.Sprint(id_raw)

	found := false
	var mounted_part Mounted
	for _, mounted_p := range MountedParts {
		if strings.Trim(mounted_p.FullID, "\000") == id_val {
			mounted_part = mounted_p
			found = true
		}
	}
	if !found {
		error_list = append(error_list, "Error, Not found partition with id: "+id_val)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	var part_size int32
	var part_start int32
	if mounted_part.Part_type == Primary {
		part_size = mounted_part.Part.Part_size
		part_start = mounted_part.Part.Part_start
	} else if mounted_part.Part_type == Logic {
		part_size = mounted_part.Logica.Part_size
		part_start = mounted_part.Logica.Part_start
	} else {
		error_list = append(error_list, "Error, Partition with unrecognized type. Name: "+id_val)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	super_bloque := Superbloque{}
	bm_start_inodes := part_start + int32(unsafe.Sizeof(super_bloque))
	//files_block := FilesBlock{}
	inode := Inode{}
	for i := 0; i < 15; i++ {
		inode.I_block[i] = -1
	}

	n := ((part_size - int32(unsafe.Sizeof(super_bloque))) / int32(1+3+unsafe.Sizeof(inode)+3*64))

	super_bloque.S_filesystem_type = 0
	super_bloque.S_inodes_count = n
	super_bloque.S_blocks_count = 3 * n
	super_bloque.S_free_inodes_count = super_bloque.S_inodes_count - 2
	super_bloque.S_free_blocks_count = super_bloque.S_blocks_count - 2
	//super_bloque.S_mtime = getCurrentTime();
	super_bloque.S_mnt_count = 1
	super_bloque.S_magic = 0xEF53 //61267;
	super_bloque.S_inode_size = int32(unsafe.Sizeof(inode))
	super_bloque.S_block_size = 64 //int32(unsafe.Sizeof(files_block))
	super_bloque.S_first_ino = 2
	super_bloque.S_first_blo = 2
	super_bloque.S_bm_inode_start = bm_start_inodes
	super_bloque.S_bm_block_start = super_bloque.S_bm_inode_start + n
	super_bloque.S_inode_start = super_bloque.S_bm_block_start + 3*n
	super_bloque.S_block_start = super_bloque.S_inode_start + n*int32(unsafe.Sizeof(inode))

	/*fmt.Printf("n:%v\n", n)
	fmt.Printf("super_bloque.S_bm_inode_start:%v\n", super_bloque.S_bm_inode_start)
	fmt.Printf("super_bloque.S_bm_block_start:%v\n", super_bloque.S_bm_block_start)
	fmt.Printf("super_bloque.S_inode_start:%v\n", super_bloque.S_inode_start)
	fmt.Printf("super_bloque.S_block_start:%v\n", super_bloque.S_block_start)*/

	// Bitmap Inodes
	var bitmap_inodes = make([]byte, n)
	bitmap_inodes[0] = 1
	bitmap_inodes[1] = 1

	for i := 2; int32(i) < n; i++ {
		bitmap_inodes[i] = 0
	}

	// Bitmaps blocks
	bitmap_blocks_size := 3 * n
	var bitmap_blocks = make([]byte, bitmap_blocks_size)

	bitmap_blocks[0] = 1
	bitmap_blocks[1] = 1
	for i := 2; int32(i) < 3*n; i++ {
		bitmap_blocks[i] = 0
	}

	WriteStruct(mounted_part.Path, part_start, super_bloque)
	WriteStruct(mounted_part.Path, super_bloque.S_bm_inode_start, bitmap_inodes)
	WriteStruct(mounted_part.Path, super_bloque.S_bm_block_start, bitmap_blocks)

	//root folder

	root_folder_block := FolderBlock{}
	root_content := Content{}
	inode_folder := Inode{}

	root_content.B_inodo = 0
	root_content.B_name = StringToByte12(".")
	root_folder_block.B_content[0] = root_content

	root_content.B_inodo = 0
	root_content.B_name = StringToByte12("..")
	root_folder_block.B_content[1] = root_content

	root_content.B_name = StringToByte12("users.txt")
	root_content.B_inodo = 1
	root_folder_block.B_content[2] = root_content

	root_content.B_name = StringToByte12("")
	root_content.B_inodo = -1
	root_folder_block.B_content[3] = root_content
	WriteStruct(mounted_part.Path, super_bloque.S_block_start, root_folder_block)

	//fmt.Printf("From Mkfs: %v", Byte12ToString(root_folder_block.B_content[1].B_name))
	//group group_ = {1, 'G', "root"};
	//user user_ = {1,'U', 1, "root", "123"};
	/* char str[100] = "";
	add_group_text(str, 1, 'G', "root");
	add_user_text(str, 1, 'U', 1, "root", "123");*/
	str := "1, G, root\n1, U, root, root, 123\n"
	for i := 0; i < 15; i++ {
		inode_folder.I_block[i] = -1
	}

	inode_folder.I_block[0] = 0
	inode_folder.I_size = int32(len(str))
	inode_folder.Inode_type = Inode_folder
	inode_folder.I_gid = 1
	inode_folder.I_uid = 1
	inode_folder.I_perm = 777
	//inode_folder.i_ctime = getCurrentTime();
	//inode_folder.i_mtime = inode_folder.i_ctime;
	//inode_folder.i_atime = inode_folder.i_ctime;

	// Setting users.txt
	users_file := FilesBlock{}
	users_inode := Inode{}
	for i := 0; i < 15; i++ {
		users_inode.I_block[i] = -1
	}
	users_inode.I_block[0] = 1
	users_inode.I_size = int32(len(str))
	users_inode.Inode_type = Inode_file
	users_inode.I_uid = 1
	users_inode.I_gid = 1
	users_inode.I_perm = 700
	//users_inode.I_ctime = getCurrentTime();
	//users_inode.I_mtime = users_inode.i_ctime;
	//users_inode.I_atime = users_inode.i_ctime;

	users_file.B_content = StringToByte64(str)

	//Write folder and user file Inode
	WriteStruct(mounted_part.Path, super_bloque.S_inode_start, inode_folder)
	WriteStruct(mounted_part.Path, super_bloque.S_inode_start+int32(unsafe.Sizeof(inode_folder)), users_inode)
	WriteStruct(mounted_part.Path, super_bloque.S_block_start+int32(unsafe.Sizeof(users_file)), users_file)

	/*
		fb := FilesBlock{}
		fb, err1 := ReadFileBlock(mounted_part.Path, super_bloque.S_block_start+int32(unsafe.Sizeof(fb)))
		if err1 != "" {
			fmt.Println("Error reading folder block")
		}
		fmt.Println("\nFolderBlock")
		fmt.Printf("Content: %v\n", Byte64ToString(fb.B_content))*/

	inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Succesfully formatted partition " + mounted_part.Name + "with id:" + mounted_part.FullID}
	return inst_res

}
