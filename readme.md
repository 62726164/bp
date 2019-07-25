# Blooming Password - bp

A program that implements the [NIST 800-63-3b Banned Password Check](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-63b.pdf) using a [bloom filter](https://dl.acm.org/citation.cfm?doid=362686.362692) built from the [Have I been pwned](https://haveibeenpwned.com/Passwords) SHA1 password hash list. The Have I Been Pwned 3.0 list contains more than 517 million password hashes and is 22GB uncompressed (as of Aug 2018). A bloom filter of this list is only 887MB and will fit entirely into memory on a virtual machine or Docker container with just 2GB of RAM.

## Why a Bloom Filter?

It's one of the simplest, smallest and fastest data structures for this task. Bloom filters have constant time O(1) performance (where K is the constant) for insertion and lookup. K is the number of times a password is hashed. Bloom filters can easily handle billions of banned password hashes with very modest resources. When a test for membership returns 404 (Not found) then it's safe to use that password.

## Partial SHA1 Hashes

SHA1 hashes are 20 bytes of raw binary data and thus typically hex encoded for a total of 40 characters. Blooming Password uses just the first 16 hex encoded characters of the hashes to build the bloom filter and to test the filter for membership. The program rejects complete hashes if they are sent. False positive rates in the bloom filter are not impacted by the shortening of the SHA1 password hashes. The cardinality of the set is unchanged. The FP rate is .001 (1 in 1,000). You may verify the cardinality is unchanged after truncating the hashes.

```
  $ wc -l pwned-passwords-ordered-by-count.txt 
  517238891 pwned-passwords-ordered-by-count.txt

  $ sort -T /tmp/ -u 16.txt | wc -l
  517238891
```

## How to Construct the Partial SHA1 Hash List

```
  $ 7z e pwned-passwords-ordered-by-count.7z

  $ cut -c 1-16 pwned-passwords-ordered-by-count.txt > 16.txt

  $ head 16.txt 
  7C4A8D09CA3762AF
  F7C3BC1D808E0473
  B1B3773A05C0ED01
  ...
```

## How to Create the Bloom Filter

```
  $ load /path/to/16.txt /path/to/output.filter
```

## Test the Bloom Filter for Membership

Send the first 16 characters of the hex encoded SHA1 hash to the Blooming Password program. Some examples using curl:

  * curl -4 https://check.aws.cloud.iso.vt.edu/hashes/sha1/0123456789ABCDEF
  * curl -6 https://check.aws.cloud.iso.vt.edu/hashes/sha1/F7C3BC1D808E0473
  * curl -4 https://check.aws.cloud.iso.vt.edu/hashes/sha1/$(echo -n "secret123" | shasum | cut -c 1-16)

## Return Codes

  * [200](https://check.aws.cloud.iso.vt.edu/hashes/sha1/F7C3BC1D808E0473) - OK. The hash is probably in the bloom filter.
  * [400](https://check.aws.cloud.iso.vt.edu/hashes/sha1/PASSWORD) - Bad request. The client sent a bad request.
  * [404](https://check.aws.cloud.iso.vt.edu/hashes/sha1/0123456789ABCDEF) - Not found. The hash is definitely not in the bloom filter.

## Notes

  * Blooming Password is written in [Go](https://golang.org).
  * It uses [willf's excellent bloom filter](https://github.com/willf/bloom) implementation.
  * The Examples above are hosted in AWS ECS (Docker) with 2 GB of memory. A VPS works fine too.
  * [OPUS](https://dl.acm.org/citation.cfm?id=134593) is an example of earlier work using a much smaller filter. (Eugene Spafford, 1992).
