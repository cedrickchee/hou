# Hou :monkey:

The Hou programming language based on [Monkey](https://monkeylang.org/), but with a little twist on the syntax and features.

Hou ("Monkey" in Chinese).

## About this Project

My go-to project for practicing a new programming language is:
- Thorsten Ball's writing on interpreters (https://interpreterbook.com) and compilers (https://compilerbook.com)
- Bob Nystrom's [Crafting Interpreters](https://craftinginterpreters.com/) handbook for making programming languages

I reimplemented Monkey tree-walking interpreter in Go as a learning exercise. The code in this repository closely resembles that presented in Thorsten's book. The interpreter is fully working.

Don't miss out the [step-by-step walk-through](#step-by-step-walk-through) in this project, where each commit is a fully working part. Read the books and follow along with the commit history.

## Quick start

Start the **REPL**:

```sh
$ go get github.com/cedrickchee/hou
$ hou
This is the Hou programming language!
Feel free to type in commands
>>
```

Then entering some Hou code:

- Variable bindings

```sh
>> let name = "awesome people"
>> puts("Hello " + name)
Hello awesome people
null
```

- Functions and closures

```sh
>> let newAdder = fn(x) { fn(y) { x + y } };
>> let addTwo = newAdder(2);
>> addTwo(3);
5
```

- Arrays and hash maps

```sh
>> let music = [{"song": "We are the World", "singer": "Michael Jackson", "year": 1985}, {"song": "Help!", "singer": "The Beatles", "year": 1965}]
>> music[0]
{song: We are the World, singer: Michael Jackson, year: 1985}
>> music[1]["song"]
Help!
```

- Errors

```
>> color = "green"
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
Woops! We ran into some monkey business here!
parser errors:
	no prefix parse function for = found
```

## Development

To build, run `make`.

```sh
$ git clone https://github.com/cedrickchee/hou
$ cd hou
$ make
```

To run the tests, run `make test`.

## Step-by-step walk-through

### Writing an Interpreter

- [1.2 Define token](https://github.com/cedrickchee/hou/commit/136d6ff5f7edeff0993dd1adbfe703d8cdab1900)
- [1.3 Lexer (basic)](https://github.com/cedrickchee/hou/commit/02637fcbd622a060bfe1ed4e22cbcd8d1a72190c)
- [1.4 Lexer (extended)](https://github.com/cedrickchee/hou/commit/ccdbaac3b3d61c691a135419beabbd685df80fe9)
- [1.5 REPL (basic)](https://github.com/cedrickchee/hou/commit/02a5f71699fd595fb4a3740f6978c22e65f39743)
- [2.4 Parser (basic)](https://github.com/cedrickchee/hou/commit/b54dc0a3c69bd427ef127feae84f765ae9e04e68)
- [2.4 Parser (error handling)](https://github.com/cedrickchee/hou/commit/58b165c19e0260b912f87d7fa5ccffab0e0fc593)
- [2.5 Parser (return)](https://github.com/cedrickchee/hou/commit/c44317d4476b4dd90e2d6986c666dd85053f67c1)
- [2.6 Pratt Parser (prefix)](https://github.com/cedrickchee/hou/commit/5e371a79310346bb20ec493b8709f093943b11d4)
- [2.6 Pratt Parser (infix)](https://github.com/cedrickchee/hou/commit/0dcf09b9e784ed46d07ac87134f1d99a3eb0dd1e)
- [2.8 Parser tracing](https://github.com/cedrickchee/hou/commit/1f20416f533120beff74e572ac1e98c686d37dd1)
- [2.8 Parser (extended)](https://github.com/cedrickchee/hou/commit/02f2502f66c753e899a2f626c655478f80f22780)
- [2.9 REPL (read-parse-print-loop)](https://github.com/cedrickchee/hou/commit/56a4636be3e13d9aaa2ba713d95f6dbc66d21bb7)
- [3.4 Evaluation (Object System)](https://github.com/cedrickchee/hou/commit/bd039e1c3f60f4363e2226b17b8ec2ad4d822039)
- [3.5 Evaluate Expression (basic)](https://github.com/cedrickchee/hou/commit/b9d8599e88ec6ccc0bf07f1a52a4bb6fd48c93b1)
- [3.5 Complete the REPL](https://github.com/cedrickchee/hou/commit/ced9b41d6c2d7c1298b5e8073c1b75ad9e9652d3)
- [3.5 Evaluation (literals)](https://github.com/cedrickchee/hou/commit/b8518212ef59a56f756c327f58a7ea684be05023)
- [3.5 Evaluation (prefix expressions)](https://github.com/cedrickchee/hou/commit/659c1fa80c763b9c9d0f4a545e5514a6cd202379)
- [3.5 Evaluation (infix expressions)](https://github.com/cedrickchee/hou/commit/c9cb6f062199a08b8221c2cfef4678e1627a2159)
- [3.6 Evaluation (conditionals)](https://github.com/cedrickchee/hou/commit/d9376e13349f044d1feba5dc7aca2bd6d2a9028c)
- [3.7 Evaluation (return statements)](https://github.com/cedrickchee/hou/commit/398a7b2e388e7709ffa7c74d84053bf0928811e1)
- [3.8 Evaluation (error handling)](https://github.com/cedrickchee/hou/commit/fb7f56223c39db90b2f1a361f50383754a90f840)
- [3.9 Evaluation (bindings and environment)](https://github.com/cedrickchee/hou/commit/4a9d4505bd11f85d219d1da859298bad80ffcc65)
- [3.10 Evaluation (functions and call expressions)](https://github.com/cedrickchee/hou/commit/7ae0c379e6732daf9418c87fadc313efd19ecdb0)
- [4.2 Data Types (strings)](https://github.com/cedrickchee/hou/commit/f9e6f144de27f950b2a90e2d7b02efce099c9fed)
- [4.2 Data Types (string concatentation)](https://github.com/cedrickchee/hou/commit/8afb5f238081f52299519f0cb3a73c7b28c618dc)
- [4.3 Builtins (len)](https://github.com/cedrickchee/hou/commit/02b8bcf6d1bc47697c8f50b698b6a1206a1b4f39)
- [4.4 Data Types (arrays) ](https://github.com/cedrickchee/hou/commit/1459f031ee317925956ed54ebdbd9055a66ccd9c)
- [4.4 Arrays (index operator expressions)](https://github.com/cedrickchee/hou/commit/375deda0298b51f1b03de296e0900acd70cd7642)
- [4.4 Arrays (evaluating array literals)](https://github.com/cedrickchee/hou/commit/bda181a5b41d6c81764ba10c4e5b59c8758b3122)
- [4.4 Arrays (indexing)](https://github.com/cedrickchee/hou/commit/5f6d6e55d4947c8d8fb6c89ee2e1ac1ee3b1d464)
- [4.4 Arrays (more built-in functions)](https://github.com/cedrickchee/hou/commit/d3227c27c93368946f0845aa3d87c1912ac84bd4)
- [4.5 Hash](https://github.com/cedrickchee/hou/commit/3386b53a3198d0bb0ba47bc10f2cd534eb0bb3f8)
- [4.6 Hello World](https://github.com/cedrickchee/hou/commit/4512e9711b5d18e9ef49f2393418da44daf64575)