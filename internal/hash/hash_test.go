package hash

import (
	"fmt"
	"testing"

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

func TestValidate(t *testing.T) {
	gaugeValue := 1.234
	counterValue := int64(1234)
	key := []byte("secret")
	type args struct {
		m   model.Metric
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"success gauge",
			args{
				model.Metric{
					ID:    "test",
					MType: "gauge",
					Value: &gaugeValue,
					Hash:  Create(fmt.Sprintf("%s:gauge:%f", "test", gaugeValue), key),
				},
				key,
			},
			true,
			false,
		},
		{
			"success counter",
			args{
				model.Metric{
					ID:    "test",
					MType: "counter",
					Delta: &counterValue,
					Hash:  Create(fmt.Sprintf("%s:counter:%d", "test", counterValue), key),
				},
				key,
			},
			true,
			false,
		},
		{
			"no key error",
			args{
				model.Metric{
					ID:    "test",
					MType: "counter",
					Delta: &counterValue,
					Hash:  Create(fmt.Sprintf("%s:counter:%d", "test", counterValue), key),
				},
				nil,
			},
			true,
			false,
		},
		{
			"not valid",
			args{
				model.Metric{
					ID:    "test",
					MType: "counter",
					Delta: &counterValue,
					Hash:  Create(fmt.Sprintf("%s:counter:%d", "fail", counterValue), key),
				},
				key,
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(tt.args.m, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
