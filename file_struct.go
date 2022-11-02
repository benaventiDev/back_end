package main

type Partition struct {
	Part_status    int32 // 0 = inactivo, 1 = activo
	Partition_type int32 // primaria, extendida, logica
	Fit_type       int32 // las opciones por default
	Part_start     int32
	Part_size      int32
	Part_name      [64]byte
}

// TODO Mbr_fecha_creacion
type MBR struct {
	Mbr_tamano int32
	//Mbr_fecha_creacion time.Time
	Mbr_dsk_signature int32
	Disk_fit          int32 //TODO see if it can be removed
	Mbr_partition     [4]Partition
}

type EBR struct {
	Part_status int32
	Part_fit    int32 //b,f,w
	Part_start  int32
	Part_size   int32
	Part_next   int32 //próximo EBR || -1
	Part_name   [64]byte
}

type SpaceFit struct {
	Start     int32
	Available int32
	Prev_ebr  EBR
}

// TODO time_t
type Superbloque struct {
	S_filesystem_type   int32
	S_inodes_count      int32 // Guarda el número total de inodos
	S_blocks_count      int32 // Guarda el número total de bloques
	S_free_blocks_count int32 // Contiene el número de bloques libres
	S_free_inodes_count int32 // Contiene el número de inodos libres
	//time_t s_mtime             // Última fecha en el que el sistema fue montado
	//time_t s_umtime            // Última fecha en que el sistema fue desmontado

	S_mnt_count      int32 // Indica cuantas veces se ha montado el sistema
	S_magic          int32 // Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
	S_inode_size     int32 // Tamaño del inodo
	S_block_size     int32 // Tamaño del bloque
	S_first_ino      int32 // Primer inodo libre
	S_first_blo      int32 // Primer bloque libre
	S_bm_inode_start int32 // Guardará el inicio del bitmap de inodos
	S_bm_block_start int32 // Guardará el inicio del bitmap de bloques
	S_inode_start    int32 // Guardará el inicio de la tabla de inodos
	S_block_start    int32 // Guardará el inicio de la tabla de bloques
}

type Inode struct {
	I_uid  int32 // UID del usuario propietario del archivo o carpeta
	I_gid  int32 // GID del grupo al que pertenece el archivo o carpeta
	I_size int32 // Tamaño del archivo en bytes
	//time_t i_atime; // Última fecha en que se leyó el inodo sin modificarlo
	//time_t i_ctime; // Fecha en la que se creó el inodo
	//time_t i_mtime; // Úlitma fecha en la que se modificó el inodo
	/* Array en los que los primeros 12 registros son bloques directos.
	   El 13 será el número del bloque simple indirecto.
	   El 14 será el número del bloque doble indirecto.
	   El 15 será el número del bloque triple indirecto.
	   Si no son utilizados tendrá el valor -1 */
	I_block    [15]int32 // = {1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}; // Hace referencia al bloque apunta
	Inode_type int32     // Indica si es archivo o carpeta. 1=archivo, 0=carpeta
	I_perm     int32     // Guardará los permisos del archivo o carpeta a nivel de bits
}

type Content struct {
	B_name  [12]byte // = ""; // Nombre de la carpeta o archivo
	B_inodo int32    // = -1;     // Apuntador hacia un inodo asociado al archivo o carpeta
}

type FolderBlock struct {
	B_content [4]Content // Array con el contenido de la carpeta
}

type FilesBlock struct {
	B_content [64]byte // = ""; // Array con el contenido del archivo
}

type BlockPointers struct {
	B_pointers [16]int32 // = {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1};
}

/*


type struct{
    int s_filesystem_type;
    int s_inodes_count;      // Guarda el número total de inodos
    int s_blocks_count;      // Guarda el número total de bloques
    int s_free_blocks_count; // Contiene el número de bloques libres
    int s_free_inodes_count; // Contiene el número de inodos libres
    time_t s_mtime;          // Última fecha en el que el sistema fue montado
    time_t s_umtime;         // Última fecha en que el sistema fue desmontado
    int s_mnt_count;         // Indica cuantas veces se ha montado el sistema
    int s_magic;             // Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
    int s_inode_size;        // Tamaño del inodo
    int s_block_size;        // Tamaño del bloque
    int s_first_ino;         // Primer inodo libre
    int s_first_blo;         // Primer bloque libre
    int s_bm_inode_start;    // Guardará el inicio del bitmap de inodos
    int s_bm_block_start;    // Guardará el inicio del bitmap de bloques
    int s_inode_start;       // Guardará el inicio de la tabla de inodos
    int s_block_start;       // Guardará el inicio de la tabla de bloques
} Superbloque;

type struct{
    int i_uid;      // UID del usuario propietario del archivo o carpeta
    int i_gid;      // GID del grupo al que pertenece el archivo o carpeta
    int i_size;     // Tamaño del archivo en bytes
    time_t i_atime; // Última fecha en que se leyó el inodo sin modificarlo
    time_t i_ctime; // Fecha en la que se creó el inodo
    time_t i_mtime; // Úlitma fecha en la que se modificó el inodo
    //Array en los que los primeros 12 registros son bloques directos.
    //El 13 será el número del bloque simple indirecto.
    //El 14 será el número del bloque doble indirecto.
    //El 15 será el número del bloque triple indirecto.
    //Si no son utilizados tendrá el valor -1
    int i_block[15];// = {1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}; // Hace referencia al bloque apunta
    inode_type i_type;                                                                   // Indica si es archivo o carpeta. 1=archivo, 0=carpeta
    int i_perm;                                                                    // Guardará los permisos del archivo o carpeta a nivel de bits
}InodesTable;

type struct{
    char b_name[12];// = ""; // Nombre de la carpeta o archivo
    int b_inodo;// = -1;     // Apuntador hacia un inodo asociado al archivo o carpeta
}Content;

type struct{
    Content b_content[4]; // Array con el contenido de la carpeta
}FolderBlock;

type struct{
    char b_content[64];// = ""; // Array con el contenido del archivo
}FilesBlock;

type struct{
    int b_pointers[16];// = {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1};
}BlockPointers;

type struct{
    char lettter_id;
    int number_index[30];
}disk_id;

type struct{
    char carnet[3];
    //char lettter_id;
    disk_id* disk_;
    int number_id;
}partition_identifier;

type struct{
    char command[100];
}Journaling;

type struct{
    int GID;
    char tipo;
    char* nombre;
}group;

type struct{
    int UID;
    char tipo;
    int GID;
    char* nombre;
    char* contrasena;
}user;


type struct{
    char path[220];
    partition_identifier id;
    partition_type type;
    partition particion;
    EBR logica;
}mounted;

*/
