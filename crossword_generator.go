// crossword_generator.go — Go версия

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type MiniCrossword struct {
	Size          int
	Grid          [][]rune
	Solution      [][]rune
	Words         []string
	PlacedWords   []string
	WordPositions map[string][3]int
	Rand          *rand.Rand
}

func NewMiniCrossword(size int, seed int64) *MiniCrossword {
	var rng *rand.Rand
	if seed != 0 {
		rng = rand.New(rand.NewSource(seed))
	} else {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	grid := make([][]rune, size)
	solution := make([][]rune, size)
	for i := range grid {
		grid[i] = make([]rune, size)
		solution[i] = make([]rune, size)
		for j := range grid[i] {
			grid[i][j] = ' '
			solution[i][j] = ' '
		}
	}

	return &MiniCrossword{
		Size:          size,
		Grid:          grid,
		Solution:      solution,
		Words:         defaultWords(),
		PlacedWords:   []string{},
		WordPositions: make(map[string][3]int),
		Rand:          rng,
	}
}

func defaultWords() []string {
	return []string{"АРБУЗ", "БУКЕТ", "ГАМА", "ДОМ", "ЕЛЬ", "ЁЖИК", "ЖУК", "ЗОНТ", "ИГЛА"}
}

func (m *MiniCrossword) generate() {
	sortedWords := make([]string, len(m.Words))
	copy(sortedWords, m.Words)
	// Сортируем по длине (от длинных к коротким)
	for i := 0; i < len(sortedWords)-1; i++ {
		for j := i + 1; j < len(sortedWords); j++ {
			if len(sortedWords[i]) < len(sortedWords[j]) {
				sortedWords[i], sortedWords[j] = sortedWords[j], sortedWords[i]
			}
		}
	}

	for _, word := range sortedWords {
		if m.placeWord(word) {
			m.PlacedWords = append(m.PlacedWords, word)
		}
	}

	// Сохраняем решение
	for i := 0; i < m.Size; i++ {
		for j := 0; j < m.Size; j++ {
			m.Solution[i][j] = m.Grid[i][j]
		}
	}

	// Заполняем пустые клетки
	letters := []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	for i := 0; i < m.Size; i++ {
		for j := 0; j < m.Size; j++ {
			if m.Grid[i][j] == ' ' {
				m.Grid[i][j] = letters[m.Rand.Intn(len(letters))]
			}
		}
	}
}

func (m *MiniCrossword) placeWord(word string) bool {
	dirs := [][2]int{{0, 1}, {1, 0}}
	runes := []rune(word)
	for attempt := 0; attempt < 50; attempt++ {
		row := m.Rand.Intn(m.Size)
		col := m.Rand.Intn(m.Size)
		dir := dirs[m.Rand.Intn(len(dirs))]
		if m.canPlace(runes, row, col, dir) {
			m.placeWordInGrid(runes, row, col, dir)
			m.WordPositions[word] = [3]int{row, col, dir[0]*2 + dir[1]}
			return true
		}
	}
	return false
}

func (m *MiniCrossword) canPlace(runes []rune, row, col int, dir [2]int) bool {
	dr, dc := dir[0], dir[1]
	if row+dr*(len(runes)-1) >= m.Size || col+dc*(len(runes)-1) >= m.Size {
		return false
	}
	for i, ch := range runes {
		r, c := row+dr*i, col+dc*i
		if m.Grid[r][c] != ' ' && m.Grid[r][c] != ch {
			return false
		}
	}
	return true
}

func (m *MiniCrossword) placeWordInGrid(runes []rune, row, col int, dir [2]int) {
	dr, dc := dir[0], dir[1]
	for i, ch := range runes {
		r, c := row+dr*i, col+dc*i
		m.Grid[r][c] = ch
	}
}

func (m *MiniCrossword) getNumber(row, col int) int {
	num := 1
	for word, pos := range m.WordPositions {
		if pos[0] == row && pos[1] == col {
			return num
		}
		num++
	}
	return 0
}

func (m *MiniCrossword) printCrossword(showSolution bool) {
	grid := m.Grid
	if showSolution {
		grid = m.Solution
	}
	fmt.Println("┌─────┬─────┬─────┬─────┬─────┐")
	for i := 0; i < m.Size; i++ {
		fmt.Print("│")
		for j := 0; j < m.Size; j++ {
			if grid[i][j] == ' ' {
				fmt.Print("   ")
			} else {
				num := m.getNumber(i, j)
				if num > 0 && !showSolution {
					fmt.Printf("%2d ", num)
				} else {
					fmt.Printf(" %c ", grid[i][j])
				}
			}
			if j < m.Size-1 {
				fmt.Print("│")
			}
		}
		fmt.Println("│")
		if i < m.Size-1 {
			fmt.Println("├─────┼─────┼─────┼─────┼─────┤")
		}
	}
	fmt.Println("└─────┴─────┴─────┴─────┴─────┘")
}

func (m *MiniCrossword) printClues() {
	fmt.Println("\nСлова по горизонтали:")
	for _, word := range m.PlacedWords {
		if m.WordPositions[word][2] == 1 {
			fmt.Printf("%d. %s\n", len(m.PlacedWords)-indexOf(m.PlacedWords, word), word)
		}
	}
	fmt.Println("\nСлова по вертикали:")
	for _, word := range m.PlacedWords {
		if m.WordPositions[word][2] == 2 {
			fmt.Printf("%d. %s\n", len(m.PlacedWords)-indexOf(m.PlacedWords, word), word)
		}
	}
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func (m *MiniCrossword) saveJSON(filename string) {
	data := map[string]interface{}{
		"size":     m.Size,
		"grid":     m.Grid,
		"solution": m.Solution,
		"words":    m.PlacedWords,
	}
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(filename, jsonData, 0644)
	fmt.Printf("💾 Сохранено: %s\n", filename)
}

func (m *MiniCrossword) saveHTML(filename string) {
	html := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>Кроссворд</title>
<style>
body { font-family: monospace; background: #1a1a2e; color: #fff; padding: 20px; }
table { border-collapse: collapse; margin: 20px auto; }
td { width: 40px; height: 40px; text-align: center; font-size: 18px; border: 2px solid #444; background: #2a2a4e; }
td.filled { background: #3a3a6e; }
td.empty { background: #1a1a2e; }
.number { font-size: 10px; color: #888; position: relative; top: -15px; left: -15px; }
</style>
</head>
<body>
<h1 style="text-align:center;">🧩 Кроссворд</h1>
<table>
`
	for i := 0; i < m.Size; i++ {
		html += "<tr>"
		for j := 0; j < m.Size; j++ {
			cell := m.Grid[i][j]
			cls := "filled"
			if cell == ' ' {
				cls = "empty"
			}
			num := m.getNumber(i, j)
			display := string(cell)
			if cell == ' ' {
				display = "&nbsp;"
			}
			numHtml := ""
			if num > 0 {
				numHtml = fmt.Sprintf(`<span class="number">%d</span>`, num)
			}
			html += fmt.Sprintf(`<td class="%s">%s%s</td>`, cls, numHtml, display)
		}
		html += "</tr>"
	}
	html += `</table>
<p style="text-align:center;">💡 Подсказки доступны в консоли</p>
</body></html>`
	os.WriteFile(filename, []byte(html), 0644)
	fmt.Printf("💾 Сохранено: %s\n", filename)
}

func main() {
	fmt.Println("🧩 Mini Crossword Generator (Go)")
	size := 5
	fmt.Print("Размер кроссворда (5 или 7): ")
	var input string
	fmt.Scanln(&input)
	if input != "" {
		if s, err := strconv.Atoi(input); err == nil && (s == 5 || s == 7) {
			size = s
		}
	}

	seed := int64(0)
	fmt.Print("Seed (число, опционально): ")
	fmt.Scanln(&input)
	if input != "" {
		if s, err := strconv.ParseInt(input, 10, 64); err == nil {
			seed = s
		}
	}

	gen := NewMiniCrossword(size, seed)
	gen.generate()

	fmt.Println("\nКроссворд:")
	gen.printCrossword(false)

	fmt.Print("\nПоказать решение? (y/n): ")
	fmt.Scanln(&input)
	if input == "y" || input == "Y" {
		fmt.Println("\nРешение:")
		gen.printCrossword(true)
	}

	gen.printClues()
	gen.saveJSON("crossword.json")
	gen.saveHTML("crossword.html")
}
