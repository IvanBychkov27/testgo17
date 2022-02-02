package main

import (
	"fmt"
)

func main() {

	a := []int{-9, -7, -6, -1, -1, 3, 4, 6, 8, 9, 12, 21}
	n := len(a)
	k := -2

	res := make([]int, 0, 2)
	count := 0
	i := 1
	j := 0
	for {
		s := a[j]
		count++

		if s+a[n-i] == k {
			res = append(res, s, a[n-i])
			break
		}

		if s+a[n-i] < k {
			j++
			continue
		}

		if s+a[n-i] > k {
			i++
			continue
		}

		if i == j {
			break
		}

	}

	fmt.Println("count =", count)

	if len(res) != 0 {
		for _, el := range res {
			fmt.Printf("%d ", el)
		}
	} else {
		fmt.Println("None")
	}
}

//
//func main() {
//	//reader := bufio.NewReader(os.Stdin)
//	writer := bufio.NewWriter(os.Stdout)
//
//	//line, _ := reader.ReadString('\n')
//	//n, _ := strconv.Atoi(strings.TrimSpace(line))
//	//
//	//line, _ = reader.ReadString('\n')
//	//splitted := strings.Split(line, " ")
//	//
//	//a := make([]int, 0, len(splitted))
//	//for _, elem := range splitted {
//	//	x, _ := strconv.Atoi(strings.TrimSpace(elem))
//	//	a = append(a, x)
//	//}
//	//
//	//line, _ = reader.ReadString('\n')
//	//k, _ := strconv.Atoi(strings.TrimSpace(line))
//	//
//	//if len(a) != n {
//	//	fmt.Println("error input array")
//	//	return
//	//}
//
//	a := []int{-9, -7, -6, -1, -1, 3, 4, 6, 8, 9, 12, 21}
//	n := len(a)
//	k := -2
//
//	res := make([]int, 0, 2)
//	ok := false
//
//	count := 0
//	for j := 0; j < n-1; j++ {
//		s := a[j]
//		count++
//		if s + a[n-1] < k {
//			continue
//		}
//
//		if ok || s+a[j+1] > k {
//			break
//		}
//
//		ii := (n + j) / 2
//		nn := n
//		if s+a[ii] > k {
//			n = ii
//			ii = j
//		}
//
//		for i := ii; i < nn; i++ {
//			count++
//
//			if s+a[i] == k {
//				res = append(res, s, a[i])
//				ok = true
//				break
//			}
//
//			if s+a[i] > k {
//				break
//			}
//		}
//	}
//
//	fmt.Println("count =", count)
//
//	if len(res) != 0 {
//		for _, el := range res {
//			fmt.Fprintf(writer, "%d ", el)
//		}
//	} else {
//		fmt.Fprintf(writer, "None")
//	}
//
//	writer.Flush()
//}

//Рита и Гоша играют в игру. У Риты есть n фишек, на каждой из которых написано количество очков. Сначала Гоша называет число k, затем Рита должна выбрать две фишки, сумма очков на которых равна заданному числу.
//Рите надоело искать фишки самой, и она решила применить свои навыки программирования для решения этой задачи. Помогите ей написать программу для поиска нужных фишек.

//Формат ввода
//В первой строке записано количество фишек n, 2 ≤ n ≤ 104.
//Во второй строке записано n целых чисел —– очки на фишках Риты в диапазоне от -105 до 105.
//В третьей строке —– загаданное Гошей целое число k, -105 ≤ k ≤ 105.

//Формат вывода
//Нужно вывести два числа —– очки на двух фишках, в сумме дающие k.
//Если таких пар несколько, то можно вывести любую из них.
//Если таких пар не существует, то вывести «None».
//func main() {
//	reader := bufio.NewReader(os.Stdin)
//	reader = bufio.NewReaderSize(reader, 614400)
//	writer := bufio.NewWriter(os.Stdout)
//
//	line, _ := reader.ReadString('\n')
//	n, _ := strconv.Atoi(strings.TrimSpace(line))
//
//	line, _ = reader.ReadString('\n')
//	splitted := strings.Split(line, " ")
//
//	a := make([]int, 0, len(splitted))
//	for _, elem := range splitted {
//		x, _ := strconv.Atoi(strings.TrimSpace(elem))
//		a = append(a, x)
//	}
//
//	line, _ = reader.ReadString('\n')
//	k, _ := strconv.Atoi(strings.TrimSpace(line))
//
//	if len(a) != n {
//		fmt.Println("error input array")
//		return
//	}
//
//	res := make([]int, 0, 2)
//	ok := false
//
//	for j := 0; j < n-1; j++ {
//		s := a[j]
//		for i := j + 1; i < n; i++ {
//			if s+a[i] == k {
//				res = append(res, s, a[i])
//				ok = true
//				break
//			}
//		}
//		if ok {
//			break
//		}
//	}
//
//	if len(res) != 0 {
//		for _, el := range res {
//			fmt.Fprintf(writer, "%d ", el)
//		}
//	} else {
//		fmt.Fprintf(writer, "None")
//	}
//
//	writer.Flush()
//}

// Скользящее среднее
//func main() {
//	reader := bufio.NewReader(os.Stdin)
//	reader = bufio.NewReaderSize(reader, 614400)
//	writer := bufio.NewWriter(os.Stdout)
//
//	line, _ := reader.ReadString('\n')
//	n, _ := strconv.Atoi(strings.TrimSpace(line))
//
//	line, _ = reader.ReadString('\n')
//	splitted := strings.Split(line, " ")
//
//	a := make([]int, 0, len(splitted))
//	for _, elem := range splitted {
//		x, _ := strconv.Atoi(strings.TrimSpace(elem))
//		a = append(a, x)
//	}
//
//	line, _ = reader.ReadString('\n')
//	k, _ := strconv.Atoi(strings.TrimSpace(line))
//
//	if len(a) != n {
//		fmt.Println("error input array")
//		return
//	}
//
//	res := make([]float64, 0, n-k+1)
//	c := 0
//	for i := 0; i < (n - k + 1); i++ {
//		if i == 0 {
//			b := a[:k]
//			for _, v := range b {
//				c += v
//			}
//		} else {
//			c = c - a[i-1] + a[i+k-1]
//		}
//		res = append(res, float64(c)/float64(k))
//	}
//
//	for _, el := range res {
//		fmt.Fprintf(writer, "%f ", el)
//	}
//
//	writer.Flush()
//}

//
//func main() {
//	reader := bufio.NewReader(os.Stdin)
//	writer := bufio.NewWriter(os.Stdout)
//
//	line, _ := reader.ReadString('\n')
//	n, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
//
//	line, _ = reader.ReadString('\n')
//	splitted := strings.Split(line, " ")
//
//	a := make([]int64, 0, len(splitted))
//	for _, elem := range splitted {
//		x, _ := strconv.ParseInt(strings.TrimSpace(elem), 10, 64)
//		a = append(a, x)
//	}
//
//	line, _ = reader.ReadString('\n')
//	splitted = strings.Split(line, " ")
//	b := make([]int64, 0, len(splitted))
//	for _, elem := range splitted {
//		x, _ := strconv.ParseInt(strings.TrimSpace(elem), 10, 64)
//		b = append(b, x)
//	}
//
//	c := make([]int64, 0, len(splitted))
//	j := 0
//	for i := 0; i < int(2*n); i++ {
//		if i%2 == 0 {
//			c = append(c, a[j])
//		} else {
//			c = append(c, b[j])
//			j++
//		}
//	}
//
//	for _, el := range c {
//		fmt.Fprintf(writer, "%d ", el)
//	}
//
//	writer.Flush()
//}
