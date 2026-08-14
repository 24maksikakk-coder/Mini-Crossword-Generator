<?php
// crossword_generator.php — PHP версия

class MiniCrossword {
    private $size;
    private $grid;
    private $solution;
    private $words;
    private $placedWords;
    private $wordPositions;

    public function __construct($size = 5, $seed = null) {
        $this->size = $size;
        $this->grid = array_fill(0, $size, array_fill(0, $size, ' '));
        $this->solution = array_fill(0, $size, array_fill(0, $size, ' '));
        $this->words = ['АРБУЗ', 'БУКЕТ', 'ГАМА', 'ДОМ', 'ЕЛЬ', 'ЁЖИК', 'ЖУК', 'ЗОНТ', 'ИГЛА'];
        $this->placedWords = [];
        $this->wordPositions = [];
        if ($seed !== null) mt_srand($seed);
    }

    public function generate() {
        $sortedWords = $this->words;
        usort($sortedWords, function($a, $b) {
            return strlen($b) - strlen($a);
        });

        foreach ($sortedWords as $word) {
            if ($this->placeWord($word)) {
                $this->placedWords[] = $word;
            }
        }

        for ($i = 0; $i < $this->size; $i++) {
            for ($j = 0; $j < $this->size; $j++) {
                $this->solution[$i][$j] = $this->grid[$i][$j];
            }
        }

        $letters = str_split('АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ');
        for ($i = 0; $i < $this->size; $i++) {
            for ($j = 0; $j < $this->size; $j++) {
                if ($this->grid[$i][$j] == ' ') {
                    $this->grid[$i][$j] = $letters[array_rand($letters)];
                }
            }
        }
    }

    private function placeWord($word) {
        $dirs = [[0, 1], [1, 0]];
        for ($attempt = 0; $attempt < 50; $attempt++) {
            $row = rand(0, $this->size - 1);
            $col = rand(0, $this->size - 1);
            $dir = $dirs[array_rand($dirs)];
            if ($this->canPlace($word, $row, $col, $dir)) {
                $this->placeWordInGrid($word, $row, $col, $dir);
                $this->wordPositions[$word] = [$row, $col, $dir];
                return true;
            }
        }
        return false;
    }

    private function canPlace($word, $row, $col, $dir) {
        list($dr, $dc) = $dir;
        if ($row + $dr * (strlen($word) - 1) >= $this->size || $col + $dc * (strlen($word) - 1) >= $this->size) {
            return false;
        }
        $chars = str_split($word);
        foreach ($chars as $i => $ch) {
            $r = $row + $dr * $i;
            $c = $col + $dc * $i;
            if ($this->grid[$r][$c] != ' ' && $this->grid[$r][$c] != $ch) {
                return false;
            }
        }
        return true;
    }

    private function placeWordInGrid($word, $row, $col, $dir) {
        list($dr, $dc) = $dir;
        $chars = str_split($word);
        foreach ($chars as $i => $ch) {
            $r = $row + $dr * $i;
            $c = $col + $dc * $i;
            $this->grid[$r][$c] = $ch;
        }
    }

    private function getNumber($row, $col) {
        $num = 1;
        foreach ($this->wordPositions as $pos) {
            list($r, $c, $dir) = $pos;
            if ($r == $row && $c == $col) {
                return $num;
            }
            $num++;
        }
        return 0;
    }

    public function printCrossword($showSolution = false) {
        $grid = $showSolution ? $this->solution : $this->grid;
        echo "┌─────┬─────┬─────┬─────┬─────┐\n";
        for ($i = 0; $i < $this->size; $i++) {
            echo "│";
            for ($j = 0; $j < $this->size; $j++) {
                if ($grid[$i][$j] == ' ') {
                    echo "   ";
                } else {
                    $num = $this->getNumber($i, $j);
                    if ($num > 0 && !$showSolution) {
                        echo str_pad($num, 2, " ", STR_PAD_LEFT) . " ";
                    } else {
                        echo " {$grid[$i][$j]} ";
                    }
                }
                if ($j < $this->size - 1) echo "│";
            }
            echo "│\n";
            if ($i < $this->size - 1) {
                echo "├─────┼─────┼─────┼─────┼─────┤\n";
            }
        }
        echo "└─────┴─────┴─────┴─────┴─────┘\n";
    }

    public function printClues() {
        echo "\nСлова по горизонтали:\n";
        foreach ($this->placedWords as $i => $word) {
            if ($this->wordPositions[$word][2][0] == 0 && $this->wordPositions[$word][2][1] == 1) {
                echo ($i + 1) . ". $word\n";
            }
        }
        echo "\nСлова по вертикали:\n";
        foreach ($this->placedWords as $i => $word) {
            if ($this->wordPositions[$word][2][0] == 1 && $this->wordPositions[$word][2][1] == 0) {
                echo ($i + 1) . ". $word\n";
            }
        }
    }

    public function saveJSON($filename = 'crossword.json') {
        $data = [
            'size' => $this->size,
            'grid' => $this->grid,
            'solution' => $this->solution,
            'words' => $this->placedWords
        ];
        file_put_contents($filename, json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
        echo "💾 Сохранено: $filename\n";
    }

    public function saveHTML($filename = 'crossword.html') {
        $html = '<!DOCTYPE html>
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
';
        for ($i = 0; $i < $this->size; $i++) {
            $html .= "<tr>";
            for ($j = 0; $j < $this->size; $j++) {
                $cell = $this->grid[$i][$j];
                $cls = $cell == ' ' ? 'empty' : 'filled';
                $num = $this->getNumber($i, $j);
                $display = $cell == ' ' ? '&nbsp;' : $cell;
                $numHtml = $num > 0 ? "<span class=\"number\">$num</span>" : '';
                $html .= "<td class=\"$cls\">$numHtml$display</td>";
            }
            $html .= "</tr>\n";
        }
        $html .= '</table>
<p style="text-align:center;">💡 Подсказки доступны в консоли</p>
</body>
</html>';
        file_put_contents($filename, $html);
        echo "💾 Сохранено: $filename\n";
    }
}

function main() {
    echo "🧩 Mini Crossword Generator (PHP)\n";
    echo "Размер кроссворда (5 или 7): ";
    $sizeInput = trim(fgets(STDIN));
    $size = empty($sizeInput) ? 5 : (int)$sizeInput;
    $size = in_array($size, [5, 7]) ? $size : 5;

    echo "Seed (число, опционально): ";
    $seedInput = trim(fgets(STDIN));
    $seed = empty($seedInput) ? null : (int)$seedInput;

    $gen = new MiniCrossword($size, $seed);
    $gen->generate();

    echo "\nКроссворд:\n";
    $gen->printCrossword();

    echo "\nПоказать решение? (y/n): ";
    $show = trim(fgets(STDIN));
    if (strtolower($show) == 'y') {
        echo "\nРешение:\n";
        $gen->printCrossword(true);
    }

    $gen->printClues();
    $gen->saveJSON();
    $gen->saveHTML();
}

main();
?>
