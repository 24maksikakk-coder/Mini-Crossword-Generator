# crossword_generator.rb — Ruby версия

class MiniCrossword
  attr_reader :grid, :solution, :placed_words, :word_positions

  def initialize(size = 5, seed = nil)
    @size = size
    @grid = Array.new(size) { Array.new(size, ' ') }
    @solution = Array.new(size) { Array.new(size, ' ') }
    @words = ['АРБУЗ', 'БУКЕТ', 'ГАМА', 'ДОМ', 'ЕЛЬ', 'ЁЖИК', 'ЖУК', 'ЗОНТ', 'ИГЛА']
    @placed_words = []
    @word_positions = {}
    @rng = seed ? Random.new(seed) : Random.new
  end

  def generate
    sorted_words = @words.sort_by { |w| -w.length }

    sorted_words.each do |word|
      if place_word(word)
        @placed_words << word
      end
    end

    @size.times do |i|
      @size.times do |j|
        @solution[i][j] = @grid[i][j]
      end
    end

    letters = 'АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ'.chars
    @size.times do |i|
      @size.times do |j|
        @grid[i][j] = letters.sample(random: @rng) if @grid[i][j] == ' '
      end
    end
  end

  def place_word(word)
    dirs = [[0, 1], [1, 0]]
    50.times do
      row = @rng.rand(@size)
      col = @rng.rand(@size)
      dir = dirs.sample(random: @rng)
      if can_place?(word, row, col, dir)
        place_word_in_grid(word, row, col, dir)
        @word_positions[word] = [row, col, dir]
        return true
      end
    end
    false
  end

  def can_place?(word, row, col, dir)
    dr, dc = dir
    return false if row + dr * (word.length - 1) >= @size || col + dc * (word.length - 1) >= @size
    word.chars.each_with_index do |ch, i|
      r = row + dr * i
      c = col + dc * i
      return false if @grid[r][c] != ' ' && @grid[r][c] != ch
    end
    true
  end

  def place_word_in_grid(word, row, col, dir)
    dr, dc = dir
    word.chars.each_with_index do |ch, i|
      r = row + dr * i
      c = col + dc * i
      @grid[r][c] = ch
    end
  end

  def get_number(row, col)
    num = 1
    @word_positions.each do |_, (r, c, _)|
      return num if r == row && c == col
      num += 1
    end
    0
  end

  def print_crossword(show_solution = false)
    grid = show_solution ? @solution : @grid
    puts "┌─────┬─────┬─────┬─────┬─────┐"
    @size.times do |i|
      print "│"
      @size.times do |j|
        if grid[i][j] == ' '
          print "   "
        else
          num = get_number(i, j)
          if num > 0 && !show_solution
            print "#{num.to_s.rjust(2)} "
          else
            print " #{grid[i][j]} "
          end
        end
        print "│" if j < @size - 1
      end
      puts "│"
      puts "├─────┼─────┼─────┼─────┼─────┤" if i < @size - 1
    end
    puts "└─────┴─────┴─────┴─────┴─────┘"
  end

  def print_clues
    puts "\nСлова по горизонтали:"
    @placed_words.each_with_index do |word, i|
      if @word_positions[word][2] == [0, 1]
        puts "#{i + 1}. #{word}"
      end
    end
    puts "\nСлова по вертикали:"
    @placed_words.each_with_index do |word, i|
      if @word_positions[word][2] == [1, 0]
        puts "#{i + 1}. #{word}"
      end
    end
  end

  def save_json(filename = 'crossword.json')
    data = {
      size: @size,
      grid: @grid,
      solution: @solution,
      words: @placed_words
    }
    File.write(filename, JSON.pretty_generate(data))
    puts "💾 Сохранено: #{filename}"
  end

  def save_html(filename = 'crossword.html')
    html = <<~HTML
      <!DOCTYPE html>
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
    HTML
    @size.times do |i|
      html << "<tr>"
      @size.times do |j|
        cell = @grid[i][j]
        cls = cell == ' ' ? 'empty' : 'filled'
        num = get_number(i, j)
        display = cell == ' ' ? '&nbsp;' : cell
        num_html = num > 0 ? "<span class=\"number\">#{num}</span>" : ''
        html << "<td class=\"#{cls}\">#{num_html}#{display}</td>"
      end
      html << "</tr>\n"
    end
    html << <<~HTML
      </table>
      <p style="text-align:center;">💡 Подсказки доступны в консоли</p>
      </body>
      </html>
    HTML
    File.write(filename, html)
    puts "💾 Сохранено: #{filename}"
  end
end

def main
  puts "🧩 Mini Crossword Generator (Ruby)"
  print "Размер кроссворда (5 или 7): "
  size_input = gets.chomp
  size = size_input.empty? ? 5 : size_input.to_i
  size = 5 unless [5, 7].include?(size)

  print "Seed (число, опционально): "
  seed_input = gets.chomp
  seed = seed_input.empty? ? nil : seed_input.to_i

  gen = MiniCrossword.new(size, seed)
  gen.generate

  puts "\nКроссворд:"
  gen.print_crossword

  print "\nПоказать решение? (y/n): "
  show = gets.chomp
  if show.downcase == 'y'
    puts "\nРешение:"
    gen.print_crossword(true)
  end

  gen.print_clues
  gen.save_json
  gen.save_html
end

main if __FILE__ == $0
