package main

import (
	_ "embed"
	"encoding/json"
	"image/color"
	"strconv"
	"strings"
)

//go:embed data/color_map.json
var json_colors []byte

var color_id map[string]int // map namespaced block IDs to color IDs
// VVV initialize color_id (this runs before main)
var _ error = json.Unmarshal(json_colors, &color_id)

//go:embed data/rgb_map.json
var json_rgb_map []byte

var rgb_map color.Palette // map color IDs to rgb values
// VVV initialize rgb_map (this runs before main)
var _ error = load_rgb_map()

// NEVER CALL!!!
// load_rgb_map() parses the embedded json in json_rgb_map
// and stores it in (global variable) rgb_map
func load_rgb_map() error {
	var raw map[string]string // json data
	err := json.Unmarshal(json_rgb_map, &raw)

	rgb_map = make(color.Palette, len(raw)+1) // give rgb_map a length

	// each iteration => one color
	for i, v := range raw {
		index, _ := strconv.Atoi(i)

		var rgb [3]uint8 // rgb values
		// parse each value from a string
		// example: "216:127:51"
		for x, y := range strings.Split(v, ":") {
			k, _ := strconv.Atoi(y)
			rgb[x] = uint8(k)
		}

		// populate rgb_map with a new color.RGBA object
		rgb_map[index] = color.RGBA{rgb[0], rgb[1], rgb[2], 255}
	}

	return err
}

// check if a block is transparent using color_id
func is_transparent(block string) bool {
	return color_id[block] == 0
}
