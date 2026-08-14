

### 1. `crossword_generator.py` (Python)

```python
# crossword_generator.py — Python версия

import random
import json
from colorama import init, Fore, Style

init(autoreset=True)

class MiniCrossword:
    def __init__(self, size=5, words=None, seed=None):
        self.size = size
        self.grid = [[' ' for _ in range(size)] for _ in range(size)]
        self.solution = [[' ' for _ in range(size)] for _ in range(size)]
        self.words = words or self._default_words()
        self.placed_words = []
        self.word_positions = {}
        self.seed = seed
        if seed is not None:
            random.seed(seed)

    def _default_words(self):
        return [
            "АРБУЗ", "БУКЕТ", "ГАМА", "ДОМ", "ЕЛЬ",
            "ЁЖИК", "ЖУК", "ЗОНТ", "ИГЛА", "ЙОГУРТ"
        ]

    def generate(self):
        """Генерирует кроссворд."""
        # Сортируем слова по длине (от длинных к коротким)
        sorted_words = sorted(self.words, key=len, reverse=True)
        placed = []

        # Пробуем разместить слова
        for word in sorted_words:
            placed_in_grid = self._place_word(word)
            if placed_in_grid:
                placed.append(word)

        # Сохраняем решение
        for i in range(self.size):
            for j in range(self.size):
                self.solution[i][j] = self.grid[i][j]

        # Заполняем пустые клетки случайными буквами
        letters = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
        for i in range(self.size):
            for j in range(self.size):
                if self.grid[i][j] == ' ':
                    self.grid[i][j] = random.choice(letters)

        return self.grid

    def _place_word(self, word):
        """Пытается разместить слово в сетке."""
        directions = [(0, 1), (1, 0)]  # горизонтально, вертикально
        attempts = 50
        for _ in range(attempts):
            row = random.randint(0, self.size - 1)
            col = random.randint(0, self.size - 1)
            direction = random.choice(directions)

            if self._can_place(word, row, col, direction):
                self._place_word_in_grid(word, row, col, direction)
                return True
        return False

    def _can_place(self, word, row, col, direction):
        """Проверяет, можно ли разместить слово."""
        dr, dc = direction
        if row + dr * (len(word) - 1) >= self.size or col + dc * (len(word) - 1) >= self.size:
            return False
        for i, ch in enumerate(word):
            r, c = row + dr * i, col + dc * i
            if self.grid[r][c] != ' ' and self.grid[r][c] != ch:
                return False
        return True

    def _place_word_in_grid(self, word, row, col, direction):
        """Размещает слово в сетке."""
        dr, dc = direction
        for i, ch in enumerate(word):
            r, c = row + dr * i, col + dc * i
            self.grid[r][c] = ch
        self.placed_words.append(word)
        # Сохраняем позицию для нумерации
        self.word_positions[word] = (row, col, direction)

    def print_crossword(self, show_solution=False):
        """Выводит кроссворд в красивом формате."""
        grid = self.solution if show_solution else self.grid
        print("┌" + "───┬".join(["─────" for _ in range(self.size)]) + "───┐")
        for i in range(self.size):
            line = "│"
            for j in range(self.size):
                cell = grid[i][j]
                if cell == ' ':
                    line += "   "
                else:
                    # Проверяем, есть ли номер
                    num = self._get_number(i, j)
                    if num:
                        line += f"{num:2d}"
                    else:
                        line += f" {cell} "
                if j < self.size - 1:
                    line += "│"
            line += "│"
            print(line)
            if i < self.size - 1:
                print("├" + "───┼".join(["─────" for _ in range(self.size)]) + "───┤")
        print("└" + "───┴".join(["─────" for _ in range(self.size)]) + "───┘")

    def _get_number(self, row, col):
        """Возвращает номер слова для клетки, если она является началом."""
        for word, (r, c, direction) in self.word_positions.items():
            if r == row and c == col:
                return len(self.word_positions) - list(self.word_positions.keys()).index(word)
        return None

    def print_clues(self):
        """Выводит подсказки к словам."""
        print("\nСлова по горизонтали:")
        for word in self.placed_words:
            if self.word_positions[word][2] == (0, 1):
                print(f"{len(self.placed_words) - self.placed_words.index(word)}. {word}")
        print("\nСлова по вертикали:")
        for word in self.placed_words:
            if self.word_positions[word][2] == (1, 0):
                print(f"{len(self.placed_words) - self.placed_words.index(word)}. {word}")

    def save_json(self, filename="crossword.json"):
        data = {
            "size": self.size,
            "grid": self.grid,
            "solution": self.solution,
            "words": self.placed_words,
            "positions": {w: {"row": p[0], "col": p[1], "dir": p[2]} for w, p in self.word_positions.items()}
        }
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print(f"💾 Сохранено: {filename}")

    def save_html(self, filename="crossword.html"):
        html = """<!DOCTYPE html>
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
"""
        for i in range(self.size):
            html += "<tr>"
            for j in range(self.size):
                cell = self.grid[i][j]
                cls = "filled" if cell != ' ' else "empty"
                num = self._get_number(i, j)
                display = cell if cell != ' ' else "&nbsp;"
                num_html = f'<span class="number">{num}</span>' if num else ''
                html += f'<td class="{cls}">{num_html}{display}</td>'
            html += "</tr>"
        html += """</table>
<p style="text-align:center;">💡 Подсказки доступны в консоли</p>
</body></html>"""
        with open(filename, 'w', encoding='utf-8') as f:
            f.write(html)
        print(f"💾 Сохранено: {filename}")

def main():
    print("🧩 Mini Crossword Generator (Python)")
    size = 5
    try:
        size = int(input("Размер кроссворда (5 или 7): ") or "5")
        if size not in [5, 7]:
            size = 5
    except:
        size = 5

    seed = input("Seed (число, опционально): ").strip()
    seed = int(seed) if seed else None

    gen = MiniCrossword(size, seed=seed)
    gen.generate()

    print("\nКроссворд:")
    gen.print_crossword()

    show = input("\nПоказать решение? (y/n): ").strip().lower()
    if show == 'y':
        print("\nРешение:")
        gen.print_crossword(show_solution=True)

    gen.print_clues()
    gen.save_json()
    gen.save_html()

if __name__ == "__main__":
    main()
