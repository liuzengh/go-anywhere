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

// 给你一个链表，删除链表的倒数第 n 个结点，并且返回链表的头结点。
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	// head == nil, n  < 0 or len(list) < n
	// var fast, slow *ListNode
	low, fast := head, head
	for i := 0; i < n; i++ {
		fast = fast.Next
	}
	if fast == nil {
		return head.Next
	}
	for fast.Next != nil {
		low = low.Next
		fast = fast.Next
	}
	low.Next = low.Next.Next
	return head

}

// 将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}
	dummyHead := ListNode{Val: -1, Next: nil}
	p := &dummyHead
	for l1 != nil || l2 != nil {
		if l1 == nil {
			p.Next = l2
			l2 = l2.Next
		} else if l2 == nil {
			p.Next = l1
			l1 = l1.Next
		} else {
			if l1.Val < l2.Val {
				p.Next = l1
				l1 = l1.Next
			} else {
				p.Next = l2
				l2 = l2.Next
			}
		}
		p = p.Next
	}
	return dummyHead.Next
}
