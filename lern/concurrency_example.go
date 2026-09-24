package main

import (
	"fmt"
	"sync"
)

func Worker(id int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()                                        // สั่งบอก WaitGroup ว่างานชิ้นนี้ทำเสร็จแล้ว
	ch <- fmt.Sprintf("Workder number %d work succes", id) // คนนี้ทำงานสำเร็จแล้ว
}

func main() {
	var wg sync.WaitGroup
	resultChan := make(chan string, 5) // Channel สำหรับส่งข้อควมหากันข้าม Goroutine

	for i := 1; i <= 5; i++ {
		wg.Add(1) //แจ้งจะเพิ่มงาน 1ชิ้น
		go Worker(i, resultChan, &wg)
	}
	wg.Wait()         //รอจนกว่างานทุกคนจะเรียก wg.Done()
	close(resultChan) //ปิดchannel เมื่อส่งข้อมูลครบ

	for msg := range resultChan {
		fmt.Println(msg)
	}
}
