package model

import (
	"errors"
	"fmt"
)

//Node n
type qNode struct {
	date      int
	next      *qNode
	priority  int
	bookingid int
}

//Queue :
type Queue struct {
	front *qNode
	back  *qNode
	size  int
}

func (p *Queue) enqueue(date int, bookingid int, pr int) error {
	newNode := &qNode{
		date:      date,
		next:      nil,
		priority:  pr,
		bookingid: bookingid,
	}
	if p.front == nil {
		p.front = newNode

	} else {
		currentN := p.front
		for currentN.next != nil {
			if currentN.next.priority > pr {
				break
			}
			currentN = currentN.next
		}
		tempNode := currentN.next
		currentN.next = newNode
		newNode.next = tempNode

		//p.back.next = newNode
	}
	p.back = newNode
	p.size++
	return nil
}

func (p *Queue) dequeue() (int, int, error) {
	var date int
	var bookingid int
	if p.front == nil {
		return 0, 0, errors.New("Empty queue")
	}
	bookingid = p.front.bookingid
	date = p.front.date
	if p.size == 1 {
		p.front = nil
		p.back = nil
	} else {
		p.front = p.front.next
	}
	p.size--
	return date, bookingid, nil
}

func (p *Queue) printAllNodes() error {
	currentNode := p.front
	if currentNode == nil {
		fmt.Println("Queue is empty.")
		return nil
	}
	fmt.Printf("%+v\n", currentNode.date)
	for currentNode.next != nil {
		currentNode = currentNode.next
		fmt.Printf("%+v\n", currentNode.date)
	}

	return nil
}

func (p *Queue) isEmpty() bool {
	return p.size == 0
}
