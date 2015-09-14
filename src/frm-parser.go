package main

import (
	"os"
	"fmt"
	"bufio"
	"bytes"
	//"io"
)

func main() {
	var fp *os.File
	var err error
	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		fmt.Printf(">> read file: %s\n", os.Args[1])
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}
	buf_reader := bufio.NewReader(fp)
	key_info_length := parse_frm_header(buf_reader)
	parse_frm_key_info(buf_reader, key_info_length)
}

func parse_frm_header (buf_reader *bufio.Reader) (uint) {
	var c byte
	var x uint
	var err error
	var buf []byte
	header_buffer := make([]byte, 0x1000)
	_, err = buf_reader.Read(header_buffer)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(header_buffer)
	/* Header Section */
	fmt.Printf(">>> Header Section\n")
	buf = read_buffer(reader, 2)
	if buf[0] != 0xfe {
		panic("not matched (1)")
	}
	if buf[1] != 0x01 {
		panic("not matched (2)")
	}
	buf = read_buffer(reader, 1)
	fmt.Printf("FRM_VER+3+test(create_info->varchar) = %d\n", buf[0])
	buf = read_buffer(reader, 1)
	fmt.Printf("legacy_db_type = %d\n", buf[0])
	buf = read_buffer(reader, 1)
	if buf[0] != 0x03 {
		panic("???")
	}
	buf = read_buffer(reader, 1)
	if buf[0] != 0x00 {
		panic("zero")
	}
	buf = read_buffer(reader, 2)
	x = bytearray2int(buf, 2)
	fmt.Printf("IO_SIZE = %d\n", x)
	buf = read_buffer(reader, 2)
	x = bytearray2int(buf, 2)
	if x != 0x0001 {
		fmt.Printf("x = %04X", x)
		panic("???")
	}
	buf = read_buffer(reader, 4)
	length := bytearray2int(buf, 4)
	fmt.Printf("length(key_length+rec_length+extra_size) = %08X [4KiB page aligned]\n", length)
	buf = read_buffer(reader, 2)
	tmp_key_length := bytearray2int(buf, 2)
	fmt.Printf("tmp_key_length = %04X\n", tmp_key_length)
	buf = read_buffer(reader, 2)
	rec_length := bytearray2int(buf, 2)
	fmt.Printf("rec_length = %04X\n", rec_length)
	buf = read_buffer(reader, 4)
	max_rows := bytearray2int(buf, 4)
	fmt.Printf("max_rows = %08X\n", max_rows)
	buf = read_buffer(reader, 4)
	min_rows := bytearray2int(buf, 4)
	fmt.Printf("min_rows = %08X\n", min_rows)
	// skip 1 byte
	_, err = reader.ReadByte()
	c, err = reader.ReadByte()
	if c != 0x02 {
		fmt.Printf("pack-fields: %02X\n", c)
		panic("use long pack-fields")
	}
	buf = read_buffer(reader, 2)
	key_info_length := bytearray2int(buf, 2)
	fmt.Printf("key_info_length = %04X\n", key_info_length)
	buf = read_buffer(reader, 2)
	table_options := bytearray2int(buf, 2)
	fmt.Printf("table_options = %04X\n", table_options)
	buf = read_buffer(reader, 1)
	if buf[0] != 0x00 {
		panic("always")
	}
	buf = read_buffer(reader, 1)
	if buf[0] != 0x05 {
		panic("frm version 5")
	}
	buf = read_buffer(reader, 4)
	avg_row_length := bytearray2int(buf, 4)
	fmt.Printf("avg_row_length = %08X\n", avg_row_length)
	buf = read_buffer(reader, 1)
	fmt.Printf("default_table_charset = %02X\n", buf[0])
	buf = read_buffer(reader, 1)
	if buf[0] != 0x00 {
		panic("always")
	}
	buf = read_buffer(reader, 1)
	fmt.Printf("row_type = %02X\n", buf[0])
	buf = read_buffer(reader, 6)
	raid_support := bytearray2int(buf, 6)
	if raid_support != 0 {
		panic("raid support bit is zero")
	}
	buf = read_buffer(reader, 4)
	key_length := bytearray2int(buf, 4)
	fmt.Printf("key_length = %08X\n", key_length)
	buf = read_buffer(reader, 4)
	mysql_version_id := bytearray2int(buf, 4)
	fmt.Printf("mysql_version_id = %08X (%d)\n", mysql_version_id, mysql_version_id)
	buf = read_buffer(reader, 4)
	extra_size := bytearray2int(buf, 4)
	fmt.Printf("extra_size = %08X\n", extra_size)
	buf = read_buffer(reader, 2)
	extra_rec_buf_length := bytearray2int(buf, 2)
	fmt.Printf("extra_rec_buf_length = %02X\n", extra_rec_buf_length)
	buf = read_buffer(reader, 1)
	default_part_db_file := bytearray2int(buf, 1)
	fmt.Printf("default_part_db_file = %02X\n", default_part_db_file)
	buf = read_buffer(reader, 2)
	key_block_size := bytearray2int(buf, 2)
	fmt.Printf("key_block_size = %04X\n", key_block_size)
	fmt.Printf("<<< Header Section\n")
	return key_info_length
}

func parse_frm_key_info (buf_reader *bufio.Reader, key_info_length uint) {
	/* File Key Information Section starts from 0x1000 */
	//var buf []byte
	var err error
	var keys, key_parts uint
	fmt.Printf(">>> File Key Information Section\n")
	file_key_information_buffer := make([]byte, key_info_length)
	_, err = buf_reader.Read(file_key_information_buffer)
	if err != nil {
		panic(err)
	}
	/*
	reader := bytes.NewReader(file_key_information_buffer)
	buf = read_buffer(reader, 1)
	fmt.Printf("always_00_when_no_index = %02X\n", buf[0])
        */
	if file_key_information_buffer[0] & 0x80 != 0 {
		keys = uint(file_key_information_buffer[1] << 7) | uint(file_key_information_buffer[0] & 0x7F)
		key_parts = bytearray2int(file_key_information_buffer[2:3], 2)
	} else {
		keys = uint(file_key_information_buffer[0])
		key_parts = uint(file_key_information_buffer[1])
	}
	fmt.Printf("keys = %d key_parts = %d\n", keys, key_parts)
	fmt.Printf("<<< File Key Information Section\n")
}

func read_buffer(reader *bytes.Reader, b int) ([]byte) {
	var err error
	buf := make([]byte, b)
	_, err = reader.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func bytearray2int(buf []byte, size uint) (uint) {
	var x uint
	for i := 0; i < len(buf); i++ {
		k := uint(buf[i])
		//x = x << 8 + k
		x += k << uint(8 * i)
	}
	return x
}
