# bitreader

bitreader는 비트를 빅엔디안의 MSB부터 순차적으로 읽는 라이브러리입니다. 

## Usage

```go
package main

import (
	"github.com/snowmerak/bitreader"
	"bytes"
	"fmt"
)

func main() {
	data := []byte{254, 255, 232, 134, 86}

	bytesBuffer := bytes.NewReader(data)
	reader, err := bitreader.New(bytesBuffer)
	if err != nil {
		panic(err)
	}

	for _, v := range data {
		fmt.Printf("%08b ", v)
	}
	fmt.Println()

	for i := 0; i < 8; i++ {
		values, err := reader.Peek(10)
		if err != nil {
			panic(err)
		}
		for i, v := range values {
			fmt.Printf("%d: %b == %d\n", i, v, v)
		}
	}

	for i := 0; i < 4; i++ {
		values, err := reader.Read(10)
		if err != nil {
			panic(err)
		}
		for i, v := range values {
			fmt.Printf("%d: %b == %d\n", i, v, v)
		}
	}
}
```

### Read

`Read` 메서드는 입력한 bit size만큼 읽어서 uint8 슬라이스로 반환합니다.

이때 uint8에 저장되는 비트의 순서 또한 빅엔디안이고, 만약 비트 수가 8의 배수가 아니라면 MSB부터 채워집니다.

### Peek

`Peek` 메서드는 `Read` 메서드와 동일한 동작을 보여주지만, `Read` 메서드는 결과적으로 커서가 이동해서 다음 비트를 읽게 되는데 반해, `Peek` 메서드는 커서가 이동하지 않습니다.

### MoveTo

`MoveTo` 메서드는 커서를 입력한 bit 위치에 이동시킵니다.

