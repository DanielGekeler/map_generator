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
	fmt.Println(chunks)
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

type chunk struct {
	offset int // chunk data offset (as sectors) in region file
	length int // number of sectors
	time   int // last modification time of a chunk in epoch seconds
}

func parse_chunks_from_region(region string) []chunk {
	buf, _ := os.ReadFile(filepath) // fully read a region file => []byte
	raw_locations := buf[:4096]
	locations := make([][2]int, len(raw_locations)/4)

	for index, element := range split_bytes(raw_locations, 4) {
		pos_bytes := append([]byte{0}, element[:3]...)
		length := element[len(element)-1]

		x := bytes_to_int(pos_bytes)
		//fmt.Println(index, "\t", x, "\t", x*4096, "\t", length)

		locations[index] = [2]int{x, int(length)}
	}

	raw_time_data := buf[4096:8192]
	time_data := make([]int, len(raw_time_data)/4)

	for index, element := range split_bytes(raw_time_data, 4) {
		time_data[index] = bytes_to_int(element)
	}

	ret := make([]chunk, len(time_data))
	for i := 0; i < 1024; i++ {
		ret[i] = chunk{locations[i][0], locations[i][1], time_data[i]}
	}
	return ret

	//fmt.Println(time_data, len(time_data))
}
