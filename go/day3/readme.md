# Day 3

[back to index](https://github.com/javorszky/adventofcode2023/)

[link on aoc website](https://adventofcode.com/2023/day/3)

## Part 1

You and the Elf eventually reach a gondola lift station; he says the gondola lift will take you up to the water source, but this is as far as he can bring you. You go inside.

It doesn't take long to find the gondolas, but there seems to be a problem: they're not moving.

"Aaah!"

You turn around to see a slightly-greasy Elf with a wrench and a look of surprise. "Sorry, I wasn't expecting anyone! The gondola lift isn't working right now; it'll still be a while before I can fix it." You offer to help.

The engineer explains that an engine part seems to be missing from the engine, but nobody can figure out which one. If you can **add up all the part numbers** in the engine schematic, it should be easy to work out which part is missing.

The engine schematic (your puzzle input) consists of a visual representation of the engine. There are lots of numbers and symbols you don't really understand, but apparently **any number adjacent to a symbol**, even diagonally, is a "part number" and should be included in your sum. (Periods (.) do not count as a symbol.)

Here is an example engine schematic:

```
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..
```

In this schematic, two numbers are **not** part numbers because they are not adjacent to a symbol: `114` (top right) and `58` (middle right). Every other number is adjacent to a symbol and so **is** a part number; their sum is **`4361`**.

Of course, the actual engine schematic is much larger. **What is the sum of all of the part numbers in the engine schematic?**

### Solution

Regex makes an appearance again! I start with a very simple one: `\d+`, as in "one or more number characters", and then use `FindAllStringSubmatchIndex` on each line to get where they are. I turn the substrings in the lines by the start / end indices into integers, then create a map where the keys are numbers, and the values are a slice of 3-element arrays: row, start, end coords. I need slice, in case a given number appears more than once. That way I'd have multiple coordinates for the same number, as is expected.

Once I have that data, I'll loop through each number, and each coordinate for the number (often just the one), and then construct a string based on the line above, overhanging by one coordinate (where applicable, top left corner doesn't get this for example), then I add the same line, adding one character on either side from the schematic, then the line below, overhanging, and then I run a regex looking for symbols: `[^\.\d]` - "not a dot, not a number". If there's a match, the number is a part number, and I add it to the rolling sum.

That rolling sum is the solution.

## Part 2

The engineer finds the missing part and installs it in the engine! As the engine springs to life, you jump in the closest gondola, finally ready to ascend to the water source.

You don't seem to be going very fast, though. Maybe something is still wrong? Fortunately, the gondola has a phone labeled "help", so you pick it up and the engineer answers.

Before you can explain the situation, she suggests that you look out the window. There stands the engineer, holding a phone in one hand and waving with the other. You're going so slowly that you haven't even left the station. You exit the gondola.

The missing part wasn't the only issue - one of the gears in the engine is wrong. A **gear** is any `*` symbol that is adjacent to **exactly two part numbers**. Its **gear ratio** is the result of multiplying those two numbers together.

This time, you need to find the gear ratio of every gear and add them all up so that the engineer can figure out which gear needs to be replaced.

Consider the same engine schematic again:


```
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..
```

In this schematic, there are **two** gears. The first is in the top left; it has part numbers `467` and `35`, so its gear ratio is `16345`. The second gear is in the lower right; its gear ratio is `451490`. (The `*` adjacent to `617` is not a gear because it is only adjacent to one part number.) Adding up all of the gear ratios produces `467835`.

**What is the sum of all of the gear ratios in your engine schematic?**

### Solution

It starts with the same parsing of where the numbers are.

Figuring out how many parts a gear is adjacent though is trickier. I set up a map with a string key and a slice of integers. The string key is the coordinate of the gear: `row-column`.

For each number I'll check the line above with overhangs, the line itself, and the line below (where they exist), and then check for a gear with `strings.Index`. If that returns anything other than -1, I grab the row and column of where the gear is, and store the current number that gear neighbours in the slice value keyed by the coordinate of the gear.

Once I have that, I loop over the gear map, and where the coordinate for a gear has exactly two elements in the slice value, I multiply those together, and sum them up.

That sum is the solution.