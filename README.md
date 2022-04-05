# glox

A Lox interpreter implemented with golang

From the book: http://www.craftinginterpreters.com/

Try it Online on [Playground](http://122.112.198.192/)

## Usage

build your glox
```shell
go build glox
```

run some lox code
```shell
./glox ${InputFile}
```

here are some test cases
```shell
./glox test_case/01-scope.glox
./glox test_case/02-scope.glox
./glox test_case/03-scope.glox
```
## About lox language
```shell
print "Hello, world!";

// Boolean
print true;             // Not false.
print false;            // Not not false.

// Number
print 1234;             // An integer.
print 12.34;            // A decimal number.

// String
print "I am a string";
print "";               // The empty string.
print "123";            // This is a string, not a number.

// Nil
print nil;

// Comparison and equality
print 1 == 2;           // false.
print "cat" != "dog";   // true.

print 314 == "pi";      // false.
print 123 == "123";     // false.

print !true;            // false.
print !false;           // true.

// Logical operators
print true and false;   // false.
print true and true;    // true.
print false or false;   // false.
print true or false;    // true.

// Precedence and grouping
print (1 + 2) / 2;      // 1.5

// Variables
var breakfast = "bagels";
print breakfast;        // "bagels".
breakfast = "beignets";
print breakfast;        // "beignets".

// Control Flow
var condition = true;
if (condition) {
  print "yes";          // "yes"
} else {
  print "no";
}

var a = 1;
while (a < 10) {
  print a;
  a = a + 1;
}

for (var a = 1; a < 10; a = a + 1) {
  print a;
}

// Functions
fun sum(a, b) {
  return a + b;
}

print sum(1, 2);
```

## As plugin
```shell
// 编译成动态链接库作为插件
go build -o glox.so -buildmode=plugin *.go
go run ${宿主程序}.go glox.so
```
参考下面这个例程可以加载插件到你的程序中。

其中，play函数作为调用glox解释器的入口，buf保存每次调用glox解释器执行的输出结果。
```go
//
// load the application "glox"  from a plugin file "glox.so"
//
func loadPlugin(filename string) (func(string), *bytes.Buffer) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}
	xplay, err := p.Lookup("Play")
	if err != nil {
		log.Fatalf("cannot find Play in %v", filename)
	}
	play := xplay.(func(string))
	xbuf, err := p.Lookup("Buf")
	if err != nil {
		log.Fatalf("cannot find Buf in %v", filename)
	}
	buf := xbuf.(*bytes.Buffer)

	return play, buf
}
```

