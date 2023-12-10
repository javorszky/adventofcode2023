# Day 7

[back to index](https://github.com/javorszky/adventofcode2023/)

[link on aoc website](https://adventofcode.com/2023/day/7)

## Part 1
Your all-expenses-paid trip turns out to be a one-way, five-minute ride in an airship. (At least it's a **cool** airship!) It drops you off at the edge of a vast desert and descends back to Island Island.

"Did you bring the parts?"

You turn around to see an Elf completely covered in white clothing, wearing goggles, and riding a large camel.

"Did you bring the parts?" she asks again, louder this time. You aren't sure what parts she's looking for; you're here to figure out why the sand stopped.

"The parts! For the sand, yes! Come with me; I will show you." She beckons you onto the camel.

After riding a bit across the sands of Desert Island, you can see what look like very large rocks covering half of the horizon. The Elf explains that the rocks are all along the part of Desert Island that is directly above Island Island, making it hard to even get there. Normally, they use big machines to move the rocks and filter the sand, but the machines have broken down because Desert Island recently stopped receiving the **parts** they need to fix the machines.

You've already assumed it'll be your job to figure out why the parts stopped when she asks if you can help. You agree automatically.

Because the journey will take a few days, she offers to teach you the game of **Camel Cards**. Camel Cards is sort of similar to poker except it's designed to be easier to play while riding a camel.

In Camel Cards, you get a list of **hands**, and your goal is to order them based on the **strength** of each hand. A hand consists of five cards labeled one of `A`, `K`, `Q`, `J`, `T`, `9`, `8`, `7`, `6`, `5`, `4`, `3`, or `2`. The relative strength of each card follows this order, where A is the highest and 2 is the lowest.

Every hand is exactly one type. From strongest to weakest, they are:

 * **Five of a kind**, where all five cards have the same label: `AAAAA`
 * **Four of a kind**, where four cards have the same label and one card has a different label: `AA8AA`
 * **Full house**, where three cards have the same label, and the remaining two cards share a different label: `23332`
 * **Three of a kind**, where three cards have the same label, and the remaining two cards are each different from any other card in the hand: `TTT98`
 * **Two pair**, where two cards share one label, two other cards share a second label, and the remaining card has a third label: `23432`
 * **One pair**, where two cards share one label, and the other three cards have a different label from the pair and each other: `A23A4`
 * **High card**, where all cards' labels are distinct: `23456`

Hands are primarily ordered based on type; for example, every full house is stronger than any three of a kind.

If two hands have the same type, a second ordering rule takes effect. Start by comparing the **first card in each hand**. If these cards are different, the hand with the stronger first card is considered stronger. If the first card in each hand have the **same label**, however, then move on to considering the **second card in each hand**. If they differ, the hand with the higher second card wins; otherwise, continue with the third card in each hand, then the fourth, then the fifth.

So, `33332` and `2AAAA` are both **four of a kind** hands, but `33332` is stronger because its first card is stronger. Similarly, `77888` and `77788` are both a **full house**, but `77888` is stronger because its third card is stronger (and both hands have the same first and second card).

To play Camel Cards, you are given a list of hands and their corresponding bid (your puzzle input). For example:

```
32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
```
This example shows five hands; each hand is followed by its **bid** amount. Each hand wins an amount equal to its bid multiplied by its **rank**, where the weakest hand gets rank 1, the second-weakest hand gets rank 2, and so on up to the strongest hand. Because there are five hands in this example, the strongest hand will have rank 5 and its bid will be multiplied by 5.

So, the first step is to put the hands in order of strength:

 * `32T3K` is the only **one pair** and the other hands are all a stronger type, so it gets rank **`1`**.
 * `KK677` and `KTJJT` are both **two pair**. Their first cards both have the same label, but the second card of `KK677` is stronger (`K` vs `T`), so `KTJJT` gets rank **`2`** and `KK677` gets rank **`3`**.
 * `T55J5` and `QQQJA` are both **three of a kind**. `QQQJA` has a stronger first card, so it gets rank **`5`** and `T55J5` gets rank **`4`**.

Now, you can determine the total winnings of this set of hands by adding up the result of multiplying each hand's bid with its rank (`765` * 1 + `220` * 2 + `28` * 3 + `684` * 4 + `483` * 5). So the total winnings in this example are `6440`.

Find the rank of every hand in your set. **What are the total winnings?**

### Solution

Eh, it was ... all right? I had a bug in my sorting algorithm. So anyways.

1. parse all the hands into a hand -> bid map
2. write a function that classifies each hand
   1. to do that, create a map that holds the character -> frequency of each hand, and then
   2. check the length of the map, and if
      1. it's 5, it's a high card
      2. it's 1, it's a 5 of a kind
      3. it's 4, it's a one pair
      4. it's 2, it can be either a
         1. 4 of a kind, in which case the product of the frequencies is 4 ( 4 * 1 )
         2. or a full house, in which case the product of the frequencies is 6 ( 2 * 3 )
      5. it's 3, it can be either a
         1. two pair, product of frequencies is 4 ( 2 * 2 * 1 )
         2. three of a kind, product is 3 ( 3 * 1 * 1 )
3. group all cards with the same classification into their own buckets
4. sort each bucket separately using my super custom algorithm (it's just comparing the values of each card in a hand using `slices.SortFunc` and a custom very simple func) such that the lowest rank card is at the beginning of the sorted slice
5. create a global order of cards and append the groups in order from lowest tier to highest tier
6. and walk through that and sum up the bid * current rank, which is now in order and monotonically increasing

## Part 2

To make things a little more interesting, the Elf introduces one additional rule. Now, `J` cards are jokers - wildcards that can act like whatever card would make the hand the strongest type possible.

To balance this, **J cards are now the weakest** individual cards, weaker even than `2`. The other cards stay in the same order: `A`, `K`, `Q`, `T`, `9`, `8`, `7`, `6`, `5`, `4`, `3`, `2`, `J`.

`J` cards can pretend to be whatever card is best for the purpose of determining hand type; for example, `QJJQ2` is now considered **four of a kind**. However, for the purpose of breaking ties between two hands of the same type, `J` is always treated as `J`, not the card it's pretending to be: `JKKK2` is weaker than `QQQQ2` because `J` is weaker than `Q`.

Now, the above example goes very differently:

```
32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
```
 * `32T3K` is still the only **one pair**; it doesn't contain any jokers, so its strength doesn't increase.
 * `KK677` is now the only **two pair**, making it the second-weakest hand.
 * `T55J5`, `KTJJT`, and `QQQJA` are now all **four of a kind**! `T55J5` gets rank 3, `QQQJA` gets rank 4, and `KTJJT` gets rank 5.

With the new joker rule, the total winnings in this example are `5905`.

Using the new joker rule, find the rank of every hand in your set. What are the new total winnings?

### Solution

Mostly the same, except the value of J card gets changed, and the upgrade rules are implemented, which needs a change in the classify function:
 * first I get the original classification of the hand, then
 * if the hand does not have a J in it, or has 5 Js in it, I return the original classification, otherwise the following rules apply:
   * if it's a four of a kind or a full house originally -> five of a kind
   * if it's a three of a kind -> four of a kind
     * (we know, because if it's a three of a kind with 2 Js, then it would have been a full house, if it's a three of a kind with 3Js, then the other two are different, otherwise it would have been a full house)
   * two pair can be turned into two different things
     * if there's exactly one J in it, it upgrades to a full house. That J takes on the card of one of the pairs making it a pair and a three of a kind, which is a full house, or
     * if there are two Js in it, then it takes on the other pair, making it a four of a kind
   * if it's a one pair -> three of a kind
     * either the pair is non-J, which means there's exactly one J in there, which would make that a three of a kind, or
     * the pair is J, which means the other three are different, so the Js can turn into one of the others, which would make it a three of a kind
   * high card -> one pair
     * there's exactly one J in there, otherwise the others would have matched already, so turn that J to one of the others

No other upgrade paths are possible.

And from this, the usual grouping of the reclassified cards, and the sorting of them happens, and the rollup of the rank * bid.
