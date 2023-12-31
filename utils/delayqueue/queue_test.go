package delayqueue

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestQueue_Push(t *testing.T) {
	q := NewQueue("ns-queue", "q1", R)
	job := NewJob("ns-queue", "q1", []byte("hello msg 1"), 10, 0, 1)
	if err := q.Push(job); err != nil {
		t.Fatalf("Failed to push job into queue: %s", err)
	}

	job2 := NewJob("ns-queue", "q2", []byte("hello msg 1"), 10, 0, 1)
	if err := q.Push(job2); err != ErrWrongQueue {
		t.Fatalf("Expected to get wrong queue error, but got: %s", err)
	}
}

func TestQueue_Poll(t *testing.T) {
	q := NewQueue("ns-queue", "q2", R)
	p := NewPool(R)
	job := NewJob("ns-queue", "q2", []byte("hello msg 2"), 10, 0, 1)
	go func() {
		time.Sleep(time.Second)
		p.Add(job)
		q.Push(job)
	}()
	jobID, err := q.Poll(2, 2)
	if err != nil || jobID == "" {
		t.Fatalf("Failed to poll job from queue: %s", err)
	}
	if job.ID() != jobID {
		t.Fatal("Mismatched job")
	}
}

func TestQueue_Peek(t *testing.T) {
	q := NewQueue("ns-queue", "q3", R)
	job := NewJob("ns-queue", "q3", []byte("hello msg 3"), 10, 0, 1)
	q.Push(job)
	jobID, err := q.Peek()
	if err != nil || jobID == "" {
		t.Fatalf("Failed to peek job from queue: %s", err)
	}
	if job.ID() != jobID {
		t.Fatal("Mismatched job")
	}
}

func TestQueue_Destroy(t *testing.T) {
	q := NewQueue("ns-queue", "q4", R)
	job := NewJob("ns-queue", "q4", []byte("hello msg 4"), 10, 0, 1)
	q.Push(job)
	count, err := q.Destroy()
	if err != nil {
		t.Fatalf("Failed to destroy queue: %s", err)
	}
	if count != 1 {
		t.Fatalf("Mismatched deleted jobs count")
	}
	size, _ := q.Size()
	if size != 0 {
		t.Fatalf("Destroyed queue should be of size 0")
	}
}

func TestPopMultiQueues(t *testing.T) {
	namespace := "ns-queueName"
	queues := make([]queue, 3)
	queueNames := make([]string, 3)
	for i, queueName := range []string{"q6", "q7", "q8"} {
		queues[i] = queue{namespace: namespace, queue: queueName}
		queueNames[i] = queues[i].Encode()
	}
	gotQueueName, gotVal, err := popMultiQueues(R, queueNames, 2)
	if err != redis.Nil {
		t.Fatalf("redis nil err was expected, but got %s", err.Error())
	}
	if gotQueueName != "" || gotVal != "" || err != redis.Nil {
		t.Fatal("queueName name and value should be empty")
	}

	queueName := "q7"
	q := NewQueue(namespace, queueName, R)
	p := NewPool(R)
	msg := "hello msg 7"
	job := NewJob(namespace, queueName, []byte(msg), 30, 0, 2)
	p.Add(job)
	q.Push(job)

	gotQueueName, gotVal, err = popMultiQueues(R, queueNames, 2)
	if err != nil {
		t.Fatalf("nil err was expected, but got %s", err.Error())
	}
	if gotQueueName != q.queue.Encode() {
		t.Fatalf("invalid queueName name, %s was expected but got %s", q.queue.Encode(), gotQueueName)
	}

	// single queue condition
	queueName = "q8"
	job = NewJob(namespace, queueName, []byte(msg), 30, 0, 2)
	p = NewPool(R)
	q = NewQueue(namespace, queueName, R)
	p.Add(job)
	q.Push(job)

	gotQueueName, gotVal, err = popMultiQueues(R, []string{queueNames[2]}, 2)
	if err != nil {
		t.Fatalf("redis nil err was expected, but got %s", err.Error())
	}
	if gotQueueName != q.queue.Encode() {
		t.Fatalf("invalid queueName name, %s was expected but got %s", q.queue.Encode(), gotQueueName)
	}
}

func TestBpopMultiQueues(t *testing.T) {
	namespace := "ns-queueName"
	queues := make([]queue, 3)
	queueNames := make([]string, 3)
	for i, queueName := range []string{"q9", "q10", "q11"} {
		queues[i] = queue{namespace: namespace, queue: queueName}
		queueNames[i] = queues[i].Encode()
	}
	gotQueueName, gotVal, err := bpopMultiQueues(R, queueNames, 2, 2)
	if err != redis.Nil {
		t.Fatalf("redis nil err was expected, but got %v", err)
	}
	if gotQueueName != "" || gotVal != "" || err != redis.Nil {
		t.Fatal("queueName name and value should be empty")
	}

	// block consume return immediately
	queueName := "q9"
	q := NewQueue(namespace, queueName, R)
	p := NewPool(R)
	msg := "hello msg 9"
	job := NewJob(namespace, queueName, []byte(msg), 30, 0, 2)
	p.Add(job)
	q.Push(job)
	gotQueueName, gotVal, err = bpopMultiQueues(R, queueNames, 2, 3)
	if err != nil {
		t.Fatalf("nil err was expected, but got %s", err.Error())
	}
	if gotQueueName != q.queue.Encode() {
		t.Fatalf("invalid queueName name, %s was expected but got %s", q.queue.Encode(), gotQueueName)
	}

	queueName = "q10"
	q = NewQueue(namespace, queueName, R)
	p = NewPool(R)
	msg = "hello msg 10"
	job = NewJob(namespace, queueName, []byte(msg), 30, 0, 2)
	go func() {
		time.Sleep(time.Second)
		p.Add(job)
		q.Push(job)
	}()

	gotQueueName, gotVal, err = bpopMultiQueues(R, queueNames, 2, 3)
	if err != nil {
		t.Fatalf("nil err was expected, but got %s", err.Error())
	}
	if gotQueueName != q.queue.Encode() {
		t.Fatalf("invalid queueName name, %s was expected but got %s", q.queue.Encode(), gotQueueName)
	}

	// single queue condition
	queueName = "q11"
	job = NewJob(namespace, queueName, []byte(msg), 30, 0, 2)
	p = NewPool(R)
	q = NewQueue(namespace, queueName, R)
	go func() {
		time.Sleep(time.Second)
		p.Add(job)
		q.Push(job)
	}()

	gotQueueName, gotVal, err = bpopMultiQueues(R, []string{queueNames[2]}, 2, 3)
	if err != nil {
		t.Fatalf("redis nil err was expected, but got %s", err.Error())
	}
	if gotQueueName != q.queue.Encode() {
		t.Fatalf("invalid queueName name, %s was expected but got %s", q.queue.Encode(), gotQueueName)
	}
}
