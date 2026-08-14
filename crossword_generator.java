// crossword_generator.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.util.*;

public class crossword_generator {
    private int size;
    private char[][] grid;
    private char[][] solution;
    private List<String> words;
    private List<String> placedWords;
    private Map<String, int[]> wordPositions;
    private Random rand;

    public crossword_generator(int size, long seed) {
        this.size = size;
        this.grid = new char[size][size];
        this.solution = new char[size][size];
        this.words = Arrays.asList("АРБУЗ", "БУКЕТ", "ГАМА", "ДОМ", "ЕЛЬ", "ЁЖИК", "ЖУК", "ЗОНТ", "ИГЛА");
        this.placedWords = new ArrayList<>();
        this.wordPositions = new HashMap<>();
        this.rand = seed != 0 ? new Random(seed) : new Random();
        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                grid[i][j] = ' ';
                solution[i][j] = ' ';
            }
        }
    }

    public void generate() {
        List<String> sortedWords = new ArrayList<>(words);
        sortedWords.sort((a, b) -> b.length() - a.length());

        for (String word : sortedWords) {
            if (placeWord(word)) {
                placedWords.add(word);
            }
        }

        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                solution[i][j] = grid[i][j];
            }
        }

        String letters = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ";
        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                if (grid[i][j] == ' ') {
                    grid[i][j] = letters.charAt(rand.nextInt(letters.length()));
                }
            }
        }
    }

    private boolean placeWord(String word) {
        int[][] dirs = {{0, 1}, {1, 0}};
        for (int attempt = 0; attempt < 50; attempt++) {
            int row = rand.nextInt(size);
            int col = rand.nextInt(size);
            int[] dir = dirs[rand.nextInt(2)];
            if (canPlace(word, row, col, dir)) {
                placeWordInGrid(word, row, col, dir);
                wordPositions.put(word, new int[]{row, col, dir[0], dir[1]});
                return true;
            }
        }
        return false;
    }

    private boolean canPlace(String word, int row, int col, int[] dir) {
        int dr = dir[0], dc = dir[1];
        if (row + dr * (word.length() - 1) >= size || col + dc * (word.length() - 1) >= size) {
            return false;
        }
        for (int i = 0; i < word.length(); i++) {
            int r = row + dr * i;
            int c = col + dc * i;
            if (grid[r][c] != ' ' && grid[r][c] != word.charAt(i)) {
                return false;
            }
        }
        return true;
    }

    private void placeWordInGrid(String word, int row, int col, int[] dir) {
        int dr = dir[0], dc = dir[1];
        for (int i = 0; i < word.length(); i++) {
            int r = row + dr * i;
            int c = col + dc * i;
            grid[r][c] = word.charAt(i);
        }
    }

    private int getNumber(int row, int col) {
        int num = 1;
        for (Map.Entry<String, int[]> entry : wordPositions.entrySet()) {
            int[] pos = entry.getValue();
            if (pos[0] == row && pos[1] == col) {
                return num;
            }
            num++;
        }
        return 0;
    }

    public void printCrossword(boolean showSolution) {
        char[][] display = showSolution ? solution : grid;
        System.out.println("┌─────┬─────┬─────┬─────┬─────┐");
        for (int i = 0; i < size; i++) {
            System.out.print("│");
            for (int j = 0; j < size; j++) {
                if (display[i][j] == ' ') {
                    System.out.print("   ");
                } else {
                    int num = getNumber(i, j);
                    if (num > 0 && !showSolution) {
                        System.out.printf("%2d ", num);
                    } else {
                        System.out.printf(" %c ", display[i][j]);
                    }
                }
                if (j < size - 1) System.out.print("│");
            }
            System.out.println("│");
            if (i < size - 1) {
                System.out.println("├─────┼─────┼─────┼─────┼─────┤");
            }
        }
        System.out.println("└─────┴─────┴─────┴─────┴─────┘");
    }

    public void printClues() {
        System.out.println("\nСлова по горизонтали:");
        for (String word : placedWords) {
            int[] pos = wordPositions.get(word);
            if (pos[2] == 0 && pos[3] == 1) {
                System.out.println((placedWords.indexOf(word) + 1) + ". " + word);
            }
        }
        System.out.println("\nСлова по вертикали:");
        for (String word : placedWords) {
            int[] pos = wordPositions.get(word);
            if (pos[2] == 1 && pos[3] == 0) {
                System.out.println((placedWords.indexOf(word) + 1) + ". " + word);
            }
        }
    }

    public void saveJSON(String filename) throws IOException {
        Map<String, Object> data = new LinkedHashMap<>();
        data.put("size", size);
        data.put("grid", grid);
        data.put("solution", solution);
        data.put("words", placedWords);
        String json = new com.google.gson.GsonBuilder().setPrettyPrinting().create().toJson(data);
        Files.write(Paths.get(filename), json.getBytes());
        System.out.println("💾 Сохранено: " + filename);
    }

    public void saveHTML(String filename) throws IOException {
        StringBuilder html = new StringBuilder();
        html.append("<!DOCTYPE html>\n<html>\n<head><meta charset=\"UTF-8\"><title>Кроссворд</title>\n");
        html.append("<style>\nbody { font-family: monospace; background: #1a1a2e; color: #fff; padding: 20px; }\n");
        html.append("table { border-collapse: collapse; margin: 20px auto; }\n");
        html.append("td { width: 40px; height: 40px; text-align: center; font-size: 18px; border: 2px solid #444; background: #2a2a4e; }\n");
        html.append("td.filled { background: #3a3a6e; }\n");
        html.append("td.empty { background: #1a1a2e; }\n");
        html.append(".number { font-size: 10px; color: #888; position: relative; top: -15px; left: -15px; }\n");
        html.append("</style>\n</head>\n<body>\n<h1 style=\"text-align:center;\">🧩 Кроссворд</h1>\n<table>\n");
        for (int i = 0; i < size; i++) {
            html.append("<tr>");
            for (int j = 0; j < size; j++) {
                char cell = grid[i][j];
                String cls = cell == ' ' ? "empty" : "filled";
                int num = getNumber(i, j);
                String display = cell == ' ' ? "&nbsp;" : String.valueOf(cell);
                String numHtml = num > 0 ? "<span class=\"number\">" + num + "</span>" : "";
                html.append("<td class=\"").append(cls).append("\">").append(numHtml).append(display).append("</td>");
            }
            html.append("</tr>\n");
        }
        html.append("</table>\n<p style=\"text-align:center;\">💡 Подсказки доступны в консоли</p>\n</body>\n</html>");
        Files.write(Paths.get(filename), html.toString().getBytes());
        System.out.println("💾 Сохранено: " + filename);
    }

    public static void main(String[] args) throws IOException {
        Scanner scanner = new Scanner(System.in);
        System.out.println("🧩 Mini Crossword Generator (Java)");

        System.out.print("Размер кроссворда (5 или 7): ");
        String sizeInput = scanner.nextLine();
        int size = 5;
        if (!sizeInput.isEmpty()) {
            try {
                int s = Integer.parseInt(sizeInput);
                if (s == 5 || s == 7) size = s;
            } catch (NumberFormatException e) {}
        }

        System.out.print("Seed (число, опционально): ");
        String seedInput = scanner.nextLine();
        long seed = 0;
        if (!seedInput.isEmpty()) {
            try {
                seed = Long.parseLong(seedInput);
            } catch (NumberFormatException e) {}
        }

        crossword_generator gen = new crossword_generator(size, seed);
        gen.generate();

        System.out.println("\nКроссворд:");
        gen.printCrossword(false);

        System.out.print("\nПоказать решение? (y/n): ");
        String show = scanner.nextLine();
        if (show.equalsIgnoreCase("y")) {
            System.out.println("\nРешение:");
            gen.printCrossword(true);
        }

        gen.printClues();
        gen.saveJSON("crossword.json");
        gen.saveHTML("crossword.html");
        scanner.close();
    }
}
