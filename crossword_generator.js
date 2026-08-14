// crossword_generator.js — JavaScript версия

const fs = require('fs');

class MiniCrossword {
    constructor(size = 5, seed = null) {
        this.size = size;
        this.grid = Array.from({ length: size }, () => Array(size).fill(' '));
        this.solution = Array.from({ length: size }, () => Array(size).fill(' '));
        this.words = ['АРБУЗ', 'БУКЕТ', 'ГАМА', 'ДОМ', 'ЕЛЬ', 'ЁЖИК', 'ЖУК', 'ЗОНТ', 'ИГЛА'];
        this.placedWords = [];
        this.wordPositions = {};
        if (seed !== null) {
            this._seed = seed;
        }
    }

    _rand() {
        if (this._seed !== undefined) {
            this._seed = (this._seed * 9301 + 49297) % 233280;
            return this._seed / 233280;
        }
        return Math.random();
    }

    _shuffle(arr) {
        for (let i = arr.length - 1; i > 0; i--) {
            const j = Math.floor(this._rand() * (i + 1));
            [arr[i], arr[j]] = [arr[j], arr[i]];
        }
        return arr;
    }

    generate() {
        const sortedWords = [...this.words].sort((a, b) => b.length - a.length);
        for (const word of sortedWords) {
            if (this._placeWord(word)) {
                this.placedWords.push(word);
            }
        }

        for (let i = 0; i < this.size; i++) {
            for (let j = 0; j < this.size; j++) {
                this.solution[i][j] = this.grid[i][j];
            }
        }

        const letters = 'АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ';
        for (let i = 0; i < this.size; i++) {
            for (let j = 0; j < this.size; j++) {
                if (this.grid[i][j] === ' ') {
                    this.grid[i][j] = letters[Math.floor(this._rand() * letters.length)];
                }
            }
        }
        return this.grid;
    }

    _placeWord(word) {
        const dirs = [[0, 1], [1, 0]];
        for (let attempt = 0; attempt < 50; attempt++) {
            const row = Math.floor(this._rand() * this.size);
            const col = Math.floor(this._rand() * this.size);
            const dir = dirs[Math.floor(this._rand() * dirs.length)];
            if (this._canPlace(word, row, col, dir)) {
                this._placeWordInGrid(word, row, col, dir);
                this.wordPositions[word] = [row, col, dir];
                return true;
            }
        }
        return false;
    }

    _canPlace(word, row, col, dir) {
        const [dr, dc] = dir;
        if (row + dr * (word.length - 1) >= this.size || col + dc * (word.length - 1) >= this.size) {
            return false;
        }
        for (let i = 0; i < word.length; i++) {
            const r = row + dr * i;
            const c = col + dc * i;
            if (this.grid[r][c] !== ' ' && this.grid[r][c] !== word[i]) {
                return false;
            }
        }
        return true;
    }

    _placeWordInGrid(word, row, col, dir) {
        const [dr, dc] = dir;
        for (let i = 0; i < word.length; i++) {
            const r = row + dr * i;
            const c = col + dc * i;
            this.grid[r][c] = word[i];
        }
    }

    _getNumber(row, col) {
        let num = 1;
        for (const word of Object.keys(this.wordPositions)) {
            const [r, c, dir] = this.wordPositions[word];
            if (r === row && c === col) {
                return num;
            }
            num++;
        }
        return 0;
    }

    printCrossword(showSolution = false) {
        const grid = showSolution ? this.solution : this.grid;
        console.log('┌─────┬─────┬─────┬─────┬─────┐');
        for (let i = 0; i < this.size; i++) {
            let line = '│';
            for (let j = 0; j < this.size; j++) {
                if (grid[i][j] === ' ') {
                    line += '   ';
                } else {
                    const num = this._getNumber(i, j);
                    if (num > 0 && !showSolution) {
                        line += `${String(num).padStart(2)} `;
                    } else {
                        line += ` ${grid[i][j]} `;
                    }
                }
                if (j < this.size - 1) line += '│';
            }
            line += '│';
            console.log(line);
            if (i < this.size - 1) {
                console.log('├─────┼─────┼─────┼─────┼─────┤');
            }
        }
        console.log('└─────┴─────┴─────┴─────┴─────┘');
    }

    printClues() {
        const words = Object.keys(this.wordPositions);
        console.log('\nСлова по горизонтали:');
        for (const word of words) {
            const [r, c, dir] = this.wordPositions[word];
            if (dir[0] === 0 && dir[1] === 1) {
                console.log(`${words.indexOf(word) + 1}. ${word}`);
            }
        }
        console.log('\nСлова по вертикали:');
        for (const word of words) {
            const [r, c, dir] = this.wordPositions[word];
            if (dir[0] === 1 && dir[1] === 0) {
                console.log(`${words.indexOf(word) + 1}. ${word}`);
            }
        }
    }

    saveJSON(filename = 'crossword.json') {
        const data = {
            size: this.size,
            grid: this.grid,
            solution: this.solution,
            words: this.placedWords,
            positions: this.wordPositions
        };
        fs.writeFileSync(filename, JSON.stringify(data, null, 2));
        console.log(`💾 Сохранено: ${filename}`);
    }

    saveHTML(filename = 'crossword.html') {
        let html = `<!DOCTYPE html>
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
`;
        for (let i = 0; i < this.size; i++) {
            html += '<tr>';
            for (let j = 0; j < this.size; j++) {
                const cell = this.grid[i][j];
                const cls = cell === ' ' ? 'empty' : 'filled';
                const num = this._getNumber(i, j);
                const display = cell === ' ' ? '&nbsp;' : cell;
                const numHtml = num > 0 ? `<span class="number">${num}</span>` : '';
                html += `<td class="${cls}">${numHtml}${display}</td>`;
            }
            html += '</tr>';
        }
        html += `</table>
<p style="text-align:center;">💡 Подсказки доступны в консоли</p>
</body></html>`;
        fs.writeFileSync(filename, html);
        console.log(`💾 Сохранено: ${filename}`);
    }
}

function main() {
    const readline = require('readline').createInterface({
        input: process.stdin,
        output: process.stdout
    });

    console.log('🧩 Mini Crossword Generator (JavaScript)');

    readline.question('Размер кроссворда (5 или 7): ', (sizeInput) => {
        let size = 5;
        if (sizeInput && [5, 7].includes(parseInt(sizeInput))) {
            size = parseInt(sizeInput);
        }

        readline.question('Seed (число, опционально): ', (seedInput) => {
            const seed = seedInput ? parseInt(seedInput) : null;

            const gen = new MiniCrossword(size, seed);
            gen.generate();

            console.log('\nКроссворд:');
            gen.printCrossword();

            readline.question('\nПоказать решение? (y/n): ', (show) => {
                if (show.toLowerCase() === 'y') {
                    console.log('\nРешение:');
                    gen.printCrossword(true);
                }

                gen.printClues();
                gen.saveJSON();
                gen.saveHTML();
                readline.close();
            });
        });
    });
}

if (require.main === module) main();
