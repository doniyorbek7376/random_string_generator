# Random String Generator

Write a string generator that takes a regular expression as argument and returns strings that match the given regular expression. 

`generate(/[-+]?[0-9]{1,16}[.][0-9]{1,6}/, 10)`

should return results like
```
"-1752643936.096896"
"9519688.31"
"+1.7036"
"+65048.3876"
"-6547028036936294.111"
"07252345.650"
"-27557.78"
"7385289878518.439775"
"13981103761187.90"
"4100273498885.614"
```

The generator should be able to support following features:

`.` Match any character except newline  
`[` Start character class definition  
`]` End character class definition  
`?` 0 or 1 quantifier  
`*` 0 or more quantifiers  
`+` 1 or more quantifier  
`{` Start min/max quantifier  
`}` End min/max quantifier  

Within a character class, the following meta characters should be supported:

`^` Negate the class, but only if the first character  
`-` Indicates character range

Few more regex examples:
```
> /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{8}/
eb4bbfa4-11d4-bce9-0f56e2dc
d896fd6b-9bfe-d0ae-6fb05b52
b36dcaeb-5654-73aa-c9ec7de2

> /.{8,12}/
(<W[+]%i
7QEyw0th
rEF\Ly(C

> /[^aeiouAEIOU0-9]{5}/
#DTdH
B[n<F
rsQgV

> /[a-f-]{5}/
-cfbc
ab---
----f
```


Note: 
1. Please don't use external libraries. Especially, write the tokenizer/parser yourself. If you use tokenizer/parsers shipped with languages (go, and python seem to have regular expression parsers), your solution would be rejected.
2. Generate possibilities randomly. Sequentially listing through the possibilities might not be as elegant/useful.
3. Please strive for simplicity and clarity. Write comments and instructions on how to run, etc. If your solution makes reviewing difficult, reviewer might reject it.

# Bonus
The following features are optional.

`|` Start of alternative branch  
`(` Start subpattern  
`)` End subpattern  
`\1` back reference  

Adding support for these features would make the generator able to handle complex patterns like

`generate(/(1[0-2]|0[1-9])(:[0-5][0-9]){2} (A|P)M/, 10)`

```
"10:43:51 PM"
"10:41:31 PM"
"03:09:55 PM"
"11:19:50 AM"
"11:20:41 PM"
"01:15:54 PM"
"02:10:04 AM"
"03:43:47 PM"
"09:39:03 AM"
"11:23:46 PM"
```
