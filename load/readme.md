# Blooming Password - load

The load program creates a new bloom filter. It takes two arguments. 

  1. Path to the text file containing partial SHA1 hashes (one hash per line). The partial SHA1 hashes must be **UPPERCASE**. 
  2. Path to where you'd like to save the bloom filter.

## What the partial SHA1 hash file should look like

```bash
head /path/to/16.txt 
7C4A8D09CA3762AF
F7C3BC1D808E0473
B1B3773A05C0ED01
5BAA61E4C9B93F3F
3D4F2BF07DC1BE38
```

## How to run load

```bash
   ./load /path/to/16.txt /path/to/output.filter
```

## Notes

  * The Blooming Password **check.go** program reads the bloom filter produced by **load.go**.
  * Filters can be read in from a local file or from a URL.

