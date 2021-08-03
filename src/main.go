package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

const filepath = "region/r.0.0.mca"

func main() {
	fmt.Println("Starting")

	chunks := parse_chunks_from_region(filepath)

	for _, v := range chunks {
		fmt.Println(v)
	}
}

func split_bytes(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func bytes_to_int(input []byte) int {
	return int(binary.BigEndian.Uint32(input))
}

type chunk_meta struct {
	offset int // chunk data offset (as sectors) in region file
	length int // number of sectors
	time   int // last modification time of a chunk in epoch seconds
	x, z   int // x and z chunk coordinates inside the region file
}

func parse_chunks_from_region(region string) []chunk_meta {
	buf, _ := os.ReadFile(filepath) // fully read a region file => []byte
	locations := split_bytes(buf[:4096], 4)
	time_data := split_bytes(buf[4096:8192], 4)

	ret := make([]chunk_meta, 1024)

	for i := 0; i < len(ret); i++ {
		pos := append([]byte{0}, locations[i][:3]...)
		length := locations[i][3]

		offset := bytes_to_int(pos)
		chunk_time := bytes_to_int(time_data[i])

		x, z := calculate_chunk_pos(i)

		ret[i] = chunk_meta{offset: offset, length: int(length), time: chunk_time, x: x, z: z}
	}

	return ret
}

func calculate_chunk_pos(index int) (int, int) {
	return index % 32, index / 32
}
