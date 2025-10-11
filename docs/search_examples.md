
# Search examples

## Word search

* Search for a precise word
```
   !a1!
   "ephemeral or persistent"
```

* Search for any combination of a set of words
```
  word1 word2 word3
  recipe fish soup
```

* Search for any combination of a set of words in a chapter
```
  word1 word2 word3 \chapter dictionary
  recipe fish soup  \chapter "my recipes"
```

* Search for any combination of a set of words in named context, any chapter
```
  word1 word2 word3 \context "weird words"
  recipe fish soup  \chapter food
```

* General word search

```
  word1 word2 word3 \chapter "my chapter" yourchapter \context "weird words"
  recipe fish soup  \chapter food \context food, recipes, dishes
```

## Table of contents
```
\chapters
\chapter mychapter
```

## Notes

* Print original notes from a chapter
```
\notes mychapter

```

## Stories and sequences
```
\story (wuya)
\story Mary
\sequence "create a pod"
\seq any \in \chapter kubernetes
\story any \chapter moon
```

## Path solutions

```
\paths \from start \to target
\from !a1! \to b6
```

## Look for an arrow

\arrow succeed
\arrow 1
\arrow 232