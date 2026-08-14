// crossword_generator.rs — Rust версия

use rand::seq::SliceRandom;
use rand::thread_rng;
use std::collections::HashMap;
use std::fs;
use std::io::{self, Write};

struct MiniCrossword {
    size: usize,
    grid: Vec<Vec<char>>,
    solution: Vec<Vec<char>>,
    words: Vec<String>,
    placed_words: Vec<String>,
    word_positions: HashMap<String, (usize, usize, (isize, isize))>,
    rng: rand::rngs::ThreadRng,
}

impl MiniCrossword {
    fn new(size: usize) -> Self {
        let grid = vec![vec![' '; size]; size];
        let solution = vec![vec![' '; size]; size];
        let words = vec![
            "АРБУЗ".to_string(), "БУКЕТ".to_string(), "ГАМА".to_string(),
            "ДОМ".to_string(), "ЕЛЬ".to_string(), "ЁЖИК".to_string(),
            "ЖУК".to_string(), "ЗОНТ".to_string(), "ИГЛА".to_string()
        ];
        MiniCrossword {
            size,
            grid,
            solution,
            words,
            placed_words: Vec::new(),
            word_positions: HashMap::new(),
            rng: thread_rng(),
        }
    }

    fn generate(&mut self) {
        let mut sorted_words = self.words.clone();
        sorted_words.sort_by(|a, b| b.len().cmp(&a.len()));

        for word in sorted_words {
            if self.place_word(&word) {
                self.placed_words.push(word);
            }
        }

        for i in 0..self.size {
            for j in 0..self.size {
                self.solution[i][j] = self.grid[i][j];
            }
        }

        let letters = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ";
        let letters_vec: Vec<char> = letters.chars().collect();
        for i in 0..self.size {
            for j in 0..self.size {
                if self.grid[i][j] == ' ' {
                    self.grid[i][j] = *letters_vec.choose(&mut self.rng).unwrap();
                }
            }
        }
    }

    fn place_word(&mut self, word: &str) -> bool {
        let dirs = [(0, 1), (1, 0)];
        let chars: Vec<char> = word.chars().collect();
        for _ in 0..50 {
            let row = self.rng.gen_range(0..self.size);
            let col = self.rng.gen_range(0..self.size);
            let dir = *dirs.choose(&mut self.rng).unwrap();
            if self.can_place(&chars, row, col, dir) {
                self.place_word_in_grid(&chars, row, col, dir);
                self.word_positions.insert(word.to_string(), (row, col, dir));
                return true;
            }
        }
        false
    }

    fn can_place(&self, chars: &[char], row: usize, col: usize, dir: (isize, isize)) -> bool {
        let (dr, dc) = dir;
        if (row as isize + dr * (chars.len() as isize - 1)) >= self.size as isize ||
           (col as isize + dc * (chars.len() as isize - 1)) >= self.size as isize {
            return false;
        }
        for (i, &ch) in chars.iter().enumerate() {
            let r = (row as isize + dr * i as isize) as usize;
            let c = (col as isize + dc * i as isize) as usize;
            if self.grid[r][c] != ' ' && self.grid[r][c] != ch {
                return false;
            }
        }
        true
    }

    fn place_word_in_grid(&mut self, chars: &[char], row: usize, col: usize, dir: (isize, isize)) {
        let (dr, dc) = dir;
        for (i, &ch) in chars.iter().enumerate() {
            let r = (row as isize + dr * i as isize) as usize;
            let c = (col as isize + dc * i as isize) as usize;
            self.grid[r][c] = ch;
        }
    }

    fn get_number(&self, row: usize, col: usize) -> usize {
        let mut num = 1;
        for (_, &(r, c, _)) in self.word_positions.iter() {
            if r == row && c == col {
                return num;
            }
            num += 1;
        }
        0
    }

    fn print_crossword(&self, show_solution: bool) {
        let grid = if show_solution { &self.solution } else { &self.grid };
        println!("┌─────┬─────┬─────┬─────┬─────┐");
        for i in 0..self.size {
            print!("│");
            for j in 0..self.size {
                if grid[i][j] == ' ' {
                    print!("   ");
                } else {
                    let num = self.get_number(i, j);
                    if num > 0 && !show_solution {
                        print!("{:2} ", num);
                    } else {
                        print!(" {} ", grid[i][j]);
                    }
                }
                if j < self.size - 1 {
                    print!("│");
                }
            }
            println!("│");
            if i < self.size - 1 {
                println!("├─────┼─────┼─────┼─────┼─────┤");
            }
        }
        println!("└─────┴─────┴─────┴─────┴─────┘");
    }

    fn print_clues(&self) {
        println!("\nСлова по горизонтали:");
        for (i, word) in self.placed_words.iter().enumerate() {
            if let Some(&(_, _, dir)) = self.word_positions.get(word) {
                if dir == (0, 1) {
                    println!("{}. {}", i + 1, word);
                }
            }
        }
        println!("\nСлова по вертикали:");
        for (i, word) in self.placed_words.iter().enumerate() {
            if let Some(&(_, _, dir)) = self.word_positions.get(word) {
                if dir == (1, 0) {
                    println!("{}. {}", i + 1, word);
                }
            }
        }
    }

    fn save_json(&self, filename: &str) {
        let data = serde_json::json!({
            "size": self.size,
            "grid": self.grid,
            "solution": self.solution,
            "words": self.placed_words,
        });
        let json = serde_json::to_string_pretty(&data).unwrap();
        fs::write(filename, json).unwrap();
        println!("💾 Сохранено: {}", filename);
    }

    fn save_html(&self, filename: &str) {
        let mut html = String::from(r#"<!DOCTYPE html>
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
"#);
        for i in 0..self.size {
            html.push_str("<tr>");
            for j in 0..self.size {
                let cell = self.grid[i][j];
                let cls = if cell == ' ' { "empty" } else { "filled" };
                let num = self.get_number(i, j);
                let display = if cell == ' ' { "&nbsp;" } else { cell.to_string() };
                let num_html = if num > 0 { format!("<span class=\"number\">{}</span>", num) } else { "".to_string() };
                html.push_str(&format!("<td class=\"{}\">{}{}</td>", cls, num_html, display));
            }
            html.push_str("</tr>\n");
        }
        html.push_str(r#"</table>
<p style="text-align:center;">💡 Подсказки доступны в консоли</p>
</body>
</html>"#);
        fs::write(filename, html).unwrap();
        println!("💾 Сохранено: {}", filename);
    }
}

fn main() {
    println!("🧩 Mini Crossword Generator (Rust)");
    print!("Размер кроссворда (5 или 7): ");
    io::stdout().flush().unwrap();
    let mut input = String::new();
    io::stdin().read_line(&mut input).unwrap();
    let size = input.trim().parse::<usize>().unwrap_or(5);
    let size = if size == 5 || size == 7 { size } else { 5 };

    let mut gen = MiniCrossword::new(size);
    gen.generate();

    println!("\nКроссворд:");
    gen.print_crossword(false);

    print!("\nПоказать решение? (y/n): ");
    io::stdout().flush().unwrap();
    input.clear();
    io::stdin().read_line(&mut input).unwrap();
    if input.trim().to_lowercase() == "y" {
        println!("\nРешение:");
        gen.print_crossword(true);
    }

    gen.print_clues();
    gen.save_json("crossword.json");
    gen.save_html("crossword.html");
}
