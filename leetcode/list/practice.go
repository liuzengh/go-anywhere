package list

// Definition for singly-linked list
type ListNode struct {
	Val  int
	Next *ListNode
}

// 给你两个 非空 的链表，表示两个非负的整数。它们每位数字都是按照 逆序 的方式存储的，并且每个节点只能存储 一位 数字。
// 请你将两个数相加，并以相同形式返回一个表示和的链表。
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	newL := new(ListNode)
	newP := newL
	for p1, p2, carray := l1, l2, 0; p1 != nil || p2 != nil || carray != 0; {
		var n1, n2 int
		if p1 != nil {
			n1 = p1.Val
			p1 = p1.Next
		}
		if p2 != nil {
			n2 = p2.Val
			p2 = p2.Next
		}
		newP.Next = &ListNode{Val: (n1 + n2 + carray) % 10}
		carray = (n1 + n2 + carray) / 10
		newP = newP.Next
	}
	return newL.Next
}
