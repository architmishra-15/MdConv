# Sample Markdown Document

This is a **sample** markdown document with syntax highlighting.

## Code Examples

Here's some Go code:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    for i := 0; i < 5; i++ {
        fmt.Printf("Count: %d\n", i)
    }
}
```

And here's some Python:

```python
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

# Generate first 10 Fibonacci numbers
for i in range(10):
    print(f"F({i}) = {fibonacci(i)}")
```

## Features

- Syntax highlighting with Chroma
- Multiple programming languages supported
- Clean HTML output
- PDF conversion ready

## Regular paragraph

This is just a regular paragraph to show normal text rendering alongside code blocks.
