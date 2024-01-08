package main

import (
	"fmt"
)

type Golay struct {
	G  [12][24]int
	HT [24][12]int
	P  [12][12]int
	u  [24]int
	r  [24]int
	s  [12]int
}

func NewGolay() *Golay {
	g := &Golay{}
	g.getMatrices()

	g.P = [12][12]int{
		{1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 1},
		{0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1},
		{0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1},
		{1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1},
		{1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1},
		{1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1},
		{0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 1},
		{1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1},
		{0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	}

	return g
}

func (g *Golay) getMatrices() {
	var I [12][12]int
	for i := range I {
		for j := range I[0] {
			if i == j {
				I[i][j] = 1
			} else {
				I[i][j] = 0
			}
		}
	}

	for i := 0; i < 12; i++ {
		for j := 0; j < 24; j++ {
			if j < 12 {
				g.HT[j][i] = I[i][j]
				g.G[i][j] = g.P[i][j]
			} else {
				g.G[i][j] = I[i][j-12]
				g.HT[j][i] = g.P[i][j-12]
			}
		}
	}
}

func binAdd(a, b int) int {
	return (a + b) % 2
}

func binMult(a, b int) int {
	return (a * b) % 2
}

func weight(v [12]int) int {
	w := 0
	for _, value := range v {
		w += value
	}
	return w
}

func (g *Golay) encode(m [12]int) {
	fmt.Print("Encoded codeword: ")
	for i := 0; i < 24; i++ {
		g.u[i] = 0
		for j := 0; j < 12; j++ {
			g.u[i] = binAdd(g.u[i], binMult(m[j], g.G[j][i]))
			g.r[i] = g.u[i]
		}
		fmt.Print(g.u[i])
	}
	fmt.Println()
}

func (g *Golay) addErrors() {
	fmt.Println("Pridėti klaidų kanale (iki 5 аrba bus klaida)")

	var num int
	fmt.Scanf("%d", &num)
	positions := make([]int, num)
	for i := 0; i < num; i++ {
		fmt.Scanf("%d", &positions[i])
	}
	fmt.Println()
	fmt.Println()

	for i := 0; i < num; i++ {
		g.r[positions[i]] = binAdd(g.r[positions[i]], 1)
		fmt.Println(positions[i])
	}

}

func (g *Golay) getSyndrome() {
	for i := 0; i < 24; i++ {
		for j := 0; j < 12; j++ {
			g.s[j] = binAdd(g.s[j], binMult(g.r[i], g.HT[i][j]))
		}
	}

	fmt.Print("Syndrome: ")
	for _, value := range g.s {
		fmt.Print(value)
	}
	fmt.Println()
}

func (g *Golay) decode() {
	g.getSyndrome()

	var e [24]int
	decodable := true
	var sp [12]int

	if weight(g.s) <= 3 {
		for i := 0; i < 24; i++ {
			if i < 12 {
				e[i] = g.s[i]
			} else {
				e[i] = 0
			}
		}
		fmt.Println("w(s) = ", weight(g.s), " <= 3")
		goto result
	}

	for i := 0; i < 12; i++ {
		var spi [12]int
		for j := 0; j < 12; j++ {
			spi[j] = binAdd(g.s[j], g.P[i][j])
		}
		if weight(spi) <= 2 {
			for k := 0; k < 24; k++ {
				if k < 12 {
					e[k] = spi[k]
				} else {
					e[k] = 0
					if i == k-12 {
						e[k] = 1
					}
				}
			}
			fmt.Println("w(s + p", i, ") = ", weight(spi), " <= 2")
			goto result
		}
	}

	for i := 0; i < 12; i++ {
		for j := 0; j < 12; j++ {
			sp[j] = binAdd(sp[j], binMult(g.s[i], g.P[i][j]))
		}
	}

	if weight(sp) == 2 || weight(sp) == 3 {
		for i := 0; i < 24; i++ {
			if i < 12 {
				e[i] = 0
			} else {
				e[i] = sp[i-12]
			}
		}
		fmt.Println("w(s*P) = ", weight(sp))
		goto result
	}

	for i := 0; i < 12; i++ {
		var sppi [12]int
		for j := 0; j < 12; j++ {
			sppi[j] = binAdd(sp[j], g.P[i][j])
		}
		if weight(sppi) == 2 {
			for k := 0; k < 24; k++ {
				if k < 12 {
					e[k] = 0
					if i == k {
						e[k] = 1
					}
				} else {
					e[k] = sppi[k-12]
				}
			}
			fmt.Println("w(s*P + p", i, ") = ", weight(sppi))
			goto result
		}
	}

	decodable = false

result:
	if decodable {
		fmt.Print("Error pattern: ")
		for _, v := range e {
			fmt.Print(v)
		}
		fmt.Println()
		fmt.Print("Decoded codeword: ")
		for i := 0; i < 24; i++ {
			fmt.Print(binAdd(g.r[i], e[i]))
		}
		fmt.Println()
	} else {
		fmt.Println("Klaida: Žinutė nedekoduojama.")
	}
}

func main() {
	message := [12]int{2, 3, 1, 5, 3, 4, 2, 4, 8, 9, 2, 1}
	golay := NewGolay()
	fmt.Println("Message (separate digits by spaces or new lines): ")
	for i := 0; i < 12; i++ {
		fmt.Scanf("%d", &message[i])
	}
	golay.encode(message)
	golay.addErrors()
	golay.decode()
}
