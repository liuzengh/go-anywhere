package twopointers

// Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

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
