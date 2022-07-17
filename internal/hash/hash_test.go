package hash

import (
	"fmt"

	"github.com/andrei-cloud/go-devops/internal/model"
)

func ExampleCreate() {
	key := []byte("test_key")
	src := "The test string"

	h := Create(src, key)

	fmt.Println(h)

	// Output:
	// 502eaf07c4b145cf91957ccaafa7b2d9a353a1fa67e293755e920e758fb23991

}

func ExampleValidate() {
	key := []byte("test_key")

	v := 0.123

	m := model.Metric{
		ID:    "test",
		MType: "gauge",
		Value: &v,
		Hash:  "",
	}

	// Valid hash
	m.Hash = Create(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), key)

	r, _ := Validate(m, key)
	fmt.Println(r)

	// Invalid hash
	m.Hash = "invalid value"

	r, _ = Validate(m, key)
	fmt.Println(r)

	// Output:
	// true
	// false
}
