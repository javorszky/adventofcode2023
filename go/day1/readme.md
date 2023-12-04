# Day 1

[back to index](https://github.com/javorszky/adventofcode2023/)

[link on aoc website](https://adventofcode.com/2023/day/1)

## Part 1

Something is wrong with global snow production, and you've been selected to take a look. The Elves have even given you a map; on it, they've used stars to mark the top fifty locations that are likely to be having problems.

You've been doing this long enough to know that to restore snow operations, you need to check all fifty stars by December 25th.

Collect stars by solving puzzles. Two puzzles will be made available on each day in the Advent calendar; the second puzzle is unlocked when you complete the first. Each puzzle grants one star. Good luck!

You try to ask why they can't just use a weather machine ("not powerful enough") and where they're even sending you ("the sky") and why your map looks mostly blank ("you sure ask a lot of questions") and hang on did you just say the sky ("of course, where do you think snow comes from") when you realize that the Elves are already loading you into a trebuchet ("please hold still, we need to strap you in").

As they're making the final adjustments, they discover that their calibration document (your puzzle input) has been amended by a very young Elf who was apparently just excited to show off her art skills. Consequently, the Elves are having trouble reading the values on the document.

The newly-improved calibration document consists of lines of text; each line originally contained a specific calibration value that the Elves now need to recover. On each line, the calibration value can be found by combining the first digit and the last digit (in that order) to form a single two-digit number.

For example:

```
1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
```

In this example, the calibration values of these four lines are `12`, `38`, `15`, and `77`. Adding these together produces `142`.

Consider your entire calibration document. What is the sum of all of the calibration values?

### Solution

Each line needs to be filtered to only contain numbers. I therefore initialise a new `strings.Builder`, then
for each line I step through the characters. If it's a number (ie char point 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57), I add that to the builder.

Then I take the first and last characters of that string. If the string is just a single digit, like `2`, then the first and last are the same, and the two character string will be `22`.

From then, I need to convert that to an `int`, then sum them up.

## Part 2

Your calculation isn't quite right. It looks like some of the digits are actually spelled out with letters: `one`, `two`, `three`, `four`, `five`, `six`, `seven`, `eight`, and `nine` also count as valid "digits".

Equipped with this new information, you now need to find the real first and last digit on each line. For example:

```
two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
```

In this example, the calibration values are `29`, `83`, `13`, `24`, `42`, `14`, and `76`. Adding these together produces `281`.

What is the sum of all of the calibration values?

### Solution

This was annoying, because the example and description did not mention that cases like `eighthree` would be `83`, ie it's not a replacement, but an existence check.

Given that I reworked the data model. For each line I created an array filled with the number 0 for the same number of elements as the string was. For example the string `two1nine` would have `[]int{0 0 0 0 0 0 0 0}` as its starting position.

Then I match regexes for the words or the number. So for "1" the regex is `one|1`. With those regexes I do a `FindAllStringSubmatchIndex`, which gives me the start and end indices of each occurrence of a match in the regex. Because it doesn't matter whether it's a word or a number that matches, I can use that as is, and because no two different numbers are going to start on the same index, I can use the starting indices in the array with the zeroes.

Taking each line in the input, I loop through the regexes for the numbers, then replace the zeroes at the index positions with whatever number I'm currently checking for, if any.

Next step is filtering out the remaining zeroes, and I'm left with a bunch of numbers.

Take the first, then last, convert to int, sum up, and done.
