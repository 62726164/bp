# Blooming Password - bp

A program that implements the [NIST 800-63-3b Banned Password Check](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-63b.pdf) using a [bloom filter](https://dl.acm.org/citation.cfm?doid=362686.362692) built from the [Have I been pwned 3.0](https://haveibeenpwned.com/Passwords) SHA1 password hash list. The Have I Been Pwned 3.0 SHA1 password hash list contains more than 517 million hashes and is 22GB uncompressed (as of July 2018). The bloom filter of these SHA1 password hashes is only 860MB and will fit entirely into memory on a virtual machine or Docker container with 2GB of RAM.

## Partial SHA1 Hashes

SHA1 hashes are 20 bytes of raw binary data and thus typically hex encoded for a total of 40 characters. Blooming Password uses just the first 16 hex encoded characters of the hashes to build the bloom filter and to test the filter for membership. The program rejects complete hashes if they are sent.

False positive rates in the bloom filter are not impacted by the shortening of the SHA1 password hashes. The cardinality of the set is unchanged. The FP rate is .001 (1 in 1,000).

## Why a Bloom Filter?

It's the simplest, smallest and fastest way to accomplish this task. Bloom filters have constant time performance (where K is the constant) for insertion and lookup. They can easily handle billions of banned password hashes with very modest resources. When a test for membership returns [404](https://www.bloomingpassword.fun/hashes/sha1/0123456789ABCDEF) then it's safe to use that password.

## How to Construct the Partial SHA1 Hash List

```
  $ cut -c 1-16 pwned-passwords-ordered-by-count.txt > 16.txt

  $ wc -l pwned-passwords-ordered-by-count.txt 
  517238891 pwned-passwords-ordered-by-count.txt

  $ sort -T /tmp/ -u 16.txt | wc -l
  517238891

  $ head 16.txt 
  7C4A8D09CA3762AF
  F7C3BC1D808E0473
  B1B3773A05C0ED01
  ...
```

## How to Create the Bloom Filter

```
  load /path/to/16.txt /path/to/output.filter
```

## Test the Bloom Filter for Membership

Send the first 16 characters of the hex encoded SHA1 hash to the Blooming Password program. Some examples using curl:

  * curl -4 https://www.bloomingpassword.fun/hashes/sha1/0123456789ABCDEF
  * curl -6 https://www.bloomingpassword.fun/hashes/sha1/F7C3BC1D808E0473

## Return Codes

  * [200](https://www.bloomingpassword.fun/hashes/sha1/F7C3BC1D808E0473) - The hash is probably in the bloom filter.
  * [400](https://www.bloomingpassword.fun/hashes/sha1/PASSWORD) - The client sent a bad request.
  * [404](https://www.bloomingpassword.fun/hashes/sha1/0123456789ABCDEF) - The hash is definitely not in the bloom filter.

## Notes

  * Blooming Password is written in [Go](https://golang.org).
  * It uses [willf's excellent bloom filter](https://github.com/willf/bloom) implementation.
  * The Examples above are hosted on a [Linode](http://linode.com/) VPS with 2 GB of memory.
  * [OPUS](https://dl.acm.org/citation.cfm?id=134593) is an example of earlier work using a much smaller filter. (Eugene Spafford, 1992).
