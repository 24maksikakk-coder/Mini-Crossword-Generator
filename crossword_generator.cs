// crossword_generator.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;

class MiniCrossword {
    private int size;
    private char[][] grid;
    private char[][] solution;
    private List<string> words;
    private List<string> placedWords;
    private Dictionary<string, int[]> wordPositions;
    private Random rand;

    public MiniCrossword(int size, int seed = 0) {
        this.size = size;
        this.grid = new char[size][];
        this.solution = new char[size][];
        this.words = new List<string> { "АРБУЗ", "БУКЕТ", "ГАМА", "ДОМ", "ЕЛЬ", "ЁЖИК", "ЖУК", "ЗОНТ", "ИГЛА" };
        this.placedWords = new List<string>();
        this.wordPositions = new Dictionary<string, int[]>();
        this.rand = seed != 0 ? new Random(seed) : new Random();
        for (int i = 0; i < size; i++) {
            grid[i] = new char[size];
            solution[i] = new char[size];
            for (int j = 0; j < size; j++) {
                grid[i][j] = ' ';
                solution[i][j] = ' ';
            }
        }
    }

    public void Generate() {
        var sortedWords = words.OrderByDescending(w => w.Length).ToList();

        foreach (var word in sortedWords) {
            if (PlaceWord(word)) {
                placedWords.Add(word);
            }
        }

        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                solution[i][j] = grid[i][j];
            }
        }

        string letters = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ";
        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                if (grid[i][j] == ' ') {
                    grid[i][j] = letters[rand.Next(letters.Length)];
                }
            }
        }
    }

    private bool PlaceWord(string word) {
        int[][] dirs = { new int[] {0, 1}, new int[] {1, 0} };
        for (int attempt = 0; attempt < 50; attempt++) {
            int row = rand.Next(size);
            int col = rand.Next(size);
            var dir = dirs[rand.Next(2)];
            if (CanPlace(word, row, col, dir)) {
                PlaceWordInGrid(word, row, col, dir);
                wordPositions[word] = new int[] { row, col, dir[0], dir[1] };
                return true;
            }
        }
        return false;
    }

    private bool CanPlace(string word, int row, int col, int[] dir) {
        int dr = dir[0], dc = dir[1];
        if (row + dr * (word.Length - 1) >= size || col + dc * (word.Length - 1) >= size) {
            return false;
        }
        for (int i = 0; i < word.Length; i++) {
            int r = row + dr * i;
            int c = col + dc * i;
            if (grid[r][c] != ' ' && grid[r][c] != word[i]) {
                return false;
            }
        }
        return true;
    }

    private void PlaceWordInGrid(string word, int row, int col, int[] dir) {
        int dr = dir[0], dc = dir[1];
        for (int i = 0; i < word.Length; i++) {
            int r = row + dr * i;
            int c = col + dc * i;
            grid[r][c] = word[i];
        }
    }

    private int GetNumber(int row, int col) {
        int num = 1;
        foreach (var entry in wordPositions) {
            var pos = entry.Value;
            if (pos[0] == row && pos[1] == col) {
                return num;
            }
            num++;
        }
        return 0;
    }

    public void PrintCrossword(bool showSolution) {
        var display = showSolution ? solution : grid;
        Console.WriteLine("┌─────┬─────┬─────┬─────┬─────┐");
        for (int i = 0; i < size; i++) {
            Console.Write("│");
            for (int j = 0; j < size; j++) {
                if (display[i][j] == ' ') {
                    Console.Write("   ");
                } else {
                    int num = GetNumber(i, j);
                    if (num > 0 && !showSolution) {
                        Console.Write($"{num,2} ");
                    } else {
                        Console.Write($" {display[i][j]} ");
                    }
                }
                if (j < size - 1) Console.Write("│");
            }
            Console.WriteLine("│");
            if (i < size - 1) {
                Console.WriteLine("├─────┼─────┼─────┼─────┼─────┤");
            }
        }
        Console.WriteLine("└─────┴─────┴─────┴─────┴─────┘");
    }

    public void PrintClues() {
        Console.WriteLine("\nСлова по горизонтали:");
        foreach (var word in placedWords) {
            var pos = wordPositions[word];
            if (pos[2] == 0 && pos[3] == 1) {
                Console.WriteLine($"{placedWords.IndexOf(word) + 1}. {word}");
            }
        }
        Console.WriteLine("\nСлова по вертикали:");
        foreach (var word in placedWords) {
            var pos = wordPositions[word];
            if (pos[2] == 1 && pos[3] == 0) {
                Console.WriteLine($"{placedWords.IndexOf(word) + 1}. {word}");
            }
        }
    }

    public void SaveJSON(string filename) {
        var data = new { size, grid, solution, words = placedWords };
        string json = JsonSerializer.Serialize(data, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine($"💾 Сохранено: {filename}");
    }

    public void SaveHTML(string filename) {
        var sb = new StringBuilder();
        sb.AppendLine("<!DOCTYPE html>");
        sb.AppendLine("<html><head><meta charset=\"UTF-8\"><title>Кроссворд</title>");
        sb.AppendLine("<style>");
        sb.AppendLine("body { font-family: monospace; background: #1a1a2e; color: #fff; padding: 20px; }");
        sb.AppendLine("table { border-collapse: collapse; margin: 20px auto; }");
        sb.AppendLine("td { width: 40px; height: 40px; text-align: center; font-size: 18px; border: 2px solid #444; background: #2a2a4e; }");
        sb.AppendLine("td.filled { background: #3a3a6e; }");
        sb.AppendLine("td.empty { background: #1a1a2e; }");
        sb.AppendLine(".number { font-size: 10px; color: #888; position: relative; top: -15px; left: -15px; }");
        sb.AppendLine("</style></head><body>");
        sb.AppendLine("<h1 style=\"text-align:center;\">🧩 Кроссворд</h1><table>");
        for (int i = 0; i < size; i++) {
            sb.Append("<tr>");
            for (int j = 0; j < size; j++) {
                char cell = grid[i][j];
                string cls = cell == ' ' ? "empty" : "filled";
                int num = GetNumber(i, j);
                string display = cell == ' ' ? "&nbsp;" : cell.ToString();
                string numHtml = num > 0 ? $"<span class=\"number\">{num}</span>" : "";
                sb.Append($"<td class=\"{cls}\">{numHtml}{display}</td>");
            }
            sb.AppendLine("</tr>");
        }
        sb.AppendLine("</table>");
        sb.AppendLine("<p style=\"text-align:center;\">💡 Подсказки доступны в консоли</p>");
        sb.AppendLine("</body></html>");
        File.WriteAllText(filename, sb.ToString());
        Console.WriteLine($"💾 Сохранено: {filename}");
    }

    public static void Main() {
        Console.WriteLine("🧩 Mini Crossword Generator (C#)");
        Console.Write("Размер кроссворда (5 или 7): ");
        string input = Console.ReadLine();
        int size = 5;
        if (!string.IsNullOrEmpty(input) && int.TryParse(input, out int s) && (s == 5 || s == 7)) {
            size = s;
        }

        Console.Write("Seed (число, опционально): ");
        input = Console.ReadLine();
        int seed = 0;
        if (!string.IsNullOrEmpty(input)) {
            int.TryParse(input, out seed);
        }

        var gen = new MiniCrossword(size, seed);
        gen.Generate();

        Console.WriteLine("\nКроссворд:");
        gen.PrintCrossword(false);

        Console.Write("\nПоказать решение? (y/n): ");
        input = Console.ReadLine();
        if (input?.ToLower() == "y") {
            Console.WriteLine("\nРешение:");
            gen.PrintCrossword(true);
        }

        gen.PrintClues();
        gen.SaveJSON("crossword.json");
        gen.SaveHTML("crossword.html");
    }
}
