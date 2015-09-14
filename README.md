# mysql-debugger

## frm-parser

frm file parser

This program dumps just the internal value from frm file header.

For file for mat information, see the MySQL documentation https://dev.mysql.com/doc/internals/en/frm-file-format.html

### run

```bash
./bin/frm-parser /path/to/frm-file.frm
```

### output

```
>> read file: testdata/live_check.frm
>>> Header Section
FRM_VER+3+test(create_info->varchar) = 10
legacy_db_type = 12
IO_SIZE = 4096
length(key_length+rec_length+extra_size) = 00003000 [4KiB page aligned]
tmp_key_length = 0169
rec_length = 091A
max_rows = 00000000
min_rows = 00000000
key_info_length = 0021
table_options = 0009
avg_row_length = 00000000
default_table_charset = 21
row_type = 00
key_length = 00000169
mysql_version_id = 0000C3F1 (50161)
extra_size = 00000010
extra_rec_buf_length = 00
default_part_db_file = 00
key_block_size = 0000
<<< Header Section
>>> File Key Information Section
keys = 1 key_parts = 1
<<< File Key Information Section
```

## Caveats

File key information structure is some complicated and not yet parsed well.
You can grab how to read this file by reading `sql/table.cc` in MySQL source tree.
