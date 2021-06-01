package sort

import "sort"

type Seq struct {
	Data      []rune
	frequency map[rune]int
}

func (s Seq) Len() int {
	return len(s.Data)
}
func (s Seq) Less(i, j int) bool {
	return (s.frequency[s.Data[j]] < s.frequency[s.Data[i]]) || (s.frequency[s.Data[j]] == s.frequency[s.Data[i]] && s.Data[j] < s.Data[i])
}
func (s Seq) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func frequencySort(s string) string {
	fre := make(map[rune]int)
	for _, ch := range s {
		fre[ch] += 1
	}
	// fmt.Println(fre)
	seq := Seq{Data: []rune(s), frequency: fre}
	sort.Sort(seq)
	return string(seq.Data)
}
