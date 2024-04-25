package level_map

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mknycha/rpg/assets"
)

type Map struct {
	Name   string
	Width  int
	Height int

	indices []int
	solids  []bool
	Sprite  *ebiten.Image
}

func NewMap(fileData string, sprite *ebiten.Image, name string) *Map {
	m := &Map{}
	m.Name = name
	m.Sprite = sprite
	file, err := os.Open(fileData)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// resize scanner's capacity to handle lines over 64K
	scanner.Scan()
	firstRow := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	firstRowSplit := removeEmptyStrings(strings.Split(firstRow, " "))
	width, err := strconv.Atoi(firstRowSplit[0])
	if err != nil {
		log.Fatalf("Failed to convert: '%s' into integer: %v", firstRowSplit[0], err)
	}
	height, err := strconv.Atoi(firstRowSplit[1])
	if err != nil {
		log.Fatalf("Failed to convert: '%s' into integer: %v", firstRowSplit[1], err)
	}
	m.Width = width
	m.Height = height
	m.indices = make([]int, 0, width*height)
	m.solids = make([]bool, 0, width*height)
	for scanner.Scan() {
		for i, v := range removeEmptyStrings(strings.Split(scanner.Text(), " ")) {
			if i%2 != 0 {
				m.solids = append(m.solids, v == "1")
			} else {
				num, err := strconv.Atoi(v)
				if err != nil {
					log.Fatalf("Failed to convert: '%s' into integer: %v", v, err)
				}
				m.indices = append(m.indices, num)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return m
}

func removeEmptyStrings(arr []string) []string {
	new := make([]string, 0, len(arr))
	for _, el := range arr {
		if el == "" {
			continue
		}
		new = append(new, el)
	}
	return new
}

func (m *Map) GetIndex(x int, y int) int {
	if (x >= 0 && x < m.Width) && (y >= 0 && y < m.Height) {
		return m.indices[y*m.Width+x]
	}
	return 0
}

func (m *Map) GetSolid(x int, y int) bool {
	if (x >= 0 && x < m.Width) && (y >= 0 && y < m.Height) {
		return m.solids[y*m.Width+x]
	}
	return false
}

func NewMapVillage1() *Map {
	levelTilesAtlas, err := assets.GetAsset("zoria_msx")
	if err != nil {
		log.Fatalf("failed to create map village 1: %v", err)
	}
	return NewMap("./test1.lvl", levelTilesAtlas, "Coders Town")
}
