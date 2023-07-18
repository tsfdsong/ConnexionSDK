package tools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//include quo
func TestGenCode(t *testing.T) {
	s := GenCode(8)
	assert.Equal(t, len(s), 8)
}

func TestAsDashboardDisplayAmount(t *testing.T) {
	type args struct {
		s       string
		n       int
		want    string
		wantErr bool
	}

	testList := []args{
		{s: "0.8888.444", n: 8, want: "", wantErr: true},
		{s: "0.8888444", n: 0, want: "", wantErr: true},
		{s: "0", n: 8, want: "0", wantErr: false},
		{s: "5", n: 8, want: "5", wantErr: false},
		{s: "0.0000000000000000000000", n: 8, want: "0", wantErr: false},
		{s: "5.0000", n: 8, want: "5", wantErr: false},
		{s: "5.0000000000000000000000", n: 8, want: "5", wantErr: false},
		{s: "0.8888444", n: 8, want: "0.8888444", wantErr: false},
		{s: "0.888844444", n: 8, want: "0.88884444", wantErr: false},
		{s: "0.88884444", n: 8, want: "0.88884444", wantErr: false},
		{s: "0.888840000000000", n: 8, want: "0.88884", wantErr: false},
		{s: "0.000000000000001", n: 8, want: "0", wantErr: false},
		{s: "0.01", n: 2, want: "0.01", wantErr: false},
		{s: "0.000100000000001", n: 8, want: "0.0001", wantErr: false},
	}

	for _, e := range testList {
		r, err := AsDashboardDisplayAmount(e.s, e.n)
		assert.Equal(t, e.wantErr, err != nil)
		assert.Equal(t, r, e.want)
	}
}

//don't delete!
// func TestQuo(t *testing.T) {
// 	type args struct {
// 		s    string
// 		want string
// 	}

// 	testList := []args{
// 		{s: "50000000000000000000000000000", want: "50000000000"},
// 		{s: "5000000000000000000000", want: "5000"},
// 		{s: "500000000", want: "0.0000000005"},
// 	}
// 	m, _ := big.NewFloat(0).SetString("1000000000000000000") //18
// 	for _, e := range testList {
// 		n, _ := big.NewFloat(0).SetString(e.s) //18
// 		f := big.NewFloat(0).Quo(n, m)
// 		fmt.Println(f)
// 		nf, _ := big.NewFloat(0).SetString(e.want)
// 		assert.Same(t, f, nf)
// 	}

// 	s, _ := big.NewInt(0).SetString("5000000000000000000000000000000000000000000000000", 0)
// 	assert.EqualValues(t, s.String(), "5000000000000000000000000000000000000000000000000")
// }

func TestGetTokenAmount2(t *testing.T) {
	type args struct {
		s       string
		p       int32
		n       int
		want    string
		wantErr bool
	}

	testList := []args{
		{s: "8888444", p: 8, n: 0, want: "", wantErr: true},
		{s: "8888444", p: 0, n: 0, want: "", wantErr: true},
		{s: "0888444", p: 0, n: 0, want: "", wantErr: true},
		{s: "0", p: 7, n: 4, want: "0", wantErr: false},
		//len(s) >= p

		//mod == 0
		{s: "80000000", p: 7, n: 4, want: "8", wantErr: false},
		//mod !=0

		//p > n
		//repeat0
		{s: "80000002", p: 7, n: 2, want: "8", wantErr: false},
		//not repeat 0
		//have zero
		{s: "80000202", p: 7, n: 6, want: "8.00002", wantErr: false},
		//have not zero
		{s: "80000221", p: 7, n: 6, want: "8.000022", wantErr: false},

		//p <= n
		//have not zero
		{s: "80000022", p: 7, n: 8, want: "8.0000022", wantErr: false},
		//have zero
		{s: "800000020", p: 7, n: 8, want: "80.000002", wantErr: false},

		//lena(s) < p
		//p > n
		//has zero
		{s: "8000", p: 8, n: 4, want: "0", wantErr: false},
		{s: "8000", p: 8, n: 6, want: "0.00008", wantErr: false},
		//have not zero
		{s: "8000", p: 8, n: 5, want: "0.00008", wantErr: false},
		{s: "8200", p: 8, n: 6, want: "0.000082", wantErr: false},

		//p<=n
		//has zero
		{s: "8000", p: 8, n: 9, want: "0.00008", wantErr: false},
		//have not zero
		{s: "8888", p: 8, n: 9, want: "0.00008888", wantErr: false},
	}

	for _, e := range testList {
		r, err := GetTokenAmount2(e.s, e.p, e.n)
		fmt.Printf("r:%+v,e:%+v\n", r, e)
		assert.Equal(t, e.wantErr, err != nil)
		assert.Equal(t, r, e.want)
	}
}

func TestGetTokenExactAmount(t *testing.T) {
	type args struct {
		s       string
		p       int32
		want    string
		wantErr bool
	}

	testList := []args{
		{s: "8888444", p: 0, want: "", wantErr: true},
		{s: "0888444", p: 3, want: "", wantErr: true},
		{s: "0", p: 8, want: "0", wantErr: false},
		{s: "876543210", p: 3, want: "876543.21", wantErr: false},
		{s: "87654321", p: 3, want: "87654.321", wantErr: false},
		{s: "80000000", p: 8, want: "0.8", wantErr: false},
		{s: "8765432101", p: 10, want: "0.8765432101", wantErr: false},
		{s: "8000", p: 8, want: "0.00008", wantErr: false},
	}

	for _, e := range testList {
		r, err := GetTokenExactAmount(e.s, e.p)
		fmt.Printf("r:%+v,e:%+v\n", r, e)
		assert.Equal(t, e.wantErr, err != nil)
		assert.Equal(t, r, e.want)
	}
}
