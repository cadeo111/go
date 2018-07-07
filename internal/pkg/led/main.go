package led

import (
	"github.com/tarm/serial"
	"log"
	"time"
	"io"
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {

	prtName := "/dev/tty.usbmodem1451"

	c := &serial.Config{Name: prtName, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Error opening serial port %v", err)
	}

	time.Sleep(2 * time.Second) // sleep because when a connection is made with an Arudino it reboots

	down := []int{
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	up := []int{
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	right := []int{
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 1, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 1, 1,
		1, 1, 2, 2, 2, 2, 1, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	left := []int{
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 1, 2, 2, 2, 2, 1, 1,
		1, 1, 2, 2, 2, 2, 2, 1, 1,
		1, 1, 1, 2, 2, 2, 2, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	}

	sendPanel(prepStoSend(sItoI8(up)), s);
	sendPanel(prepStoSend(sItoI8(down)), s);
	sendPanel(prepStoSend(sItoI8(left)), s);
	sendPanel(prepStoSend(sItoI8(right)), s);

	//sendRow(bs1[:40], s)
	//time.Sleep(150 * time.Millisecond)
	//sendRow(bs1, s)
	//time.Sleep(150 * time.Millisecond)
	//sendRow(bs2[:40], s)
	//time.Sleep(150 * time.Millisecond)
	//sendRow(bs2, s)

	//for i := 45 ; i < 81; i++ {
	//	bs = append(bs, int8(1))
	//}
	//sendRow(bs, s);

}

type  LEDpanel struct {
	size int
}


func sendRow(val []int8, s io.ReadWriteCloser) error {
	buf := new(bytes.Buffer)
	for _, i := range val {
		if err := binary.Write(buf, binary.LittleEndian, i); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
	}

	for _, b := range [][]byte{buf.Bytes()} {
		n, err := s.Write(b)
		if err != nil {
			return err
		}

		fmt.Printf("sent %d byte(s)\n", n)
	}
	return nil
}

func sendPanel(val []int8, s io.ReadWriteCloser) error {
	if err := sendRow(val[:40], s); err != nil {
		return fmt.Errorf("Send Panel first 1/2 Failed: %v", err)
	}
	time.Sleep(150 * time.Millisecond)
	sendRow(val[40:81], s)
	time.Sleep(150 * time.Millisecond)
	return nil
}

func prepStoSend(a []int8) []int8 {
	//http://golangcookbook.com/chapters/arrays/reverse/
	reverse := func(numbers []int8) []int8 {
		for i := 0; i < len(numbers)/2; i++ {
			j := len(numbers) - i - 1
			numbers[i], numbers[j] = numbers[j], numbers[i]
		}
		return numbers
	}
	var x []int8

	for i := 1; i < 10; i++ {
		if i%2 == 0 {
			x = append(x, a[9*(i-1):9*i]...)
		} else {
			x = append(x, reverse(a[9*(i-1):9*i])...)
		}
	}
	return x;
}

func sItoI8(a []int) []int8 {
	var ret []int8
	for i := 0; i < len(a); i++ {
		ret = append(ret, int8(a[i]))
	}
	return ret
}

func OpenConnection(portName string) io.ReadWriteCloser{

	//prtName := "/dev/tty.usbmodem1451"
	c := &serial.Config{Name: portName, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Error opening serial port %v", err)
	}
	return s
}

func UpdateBoard(a []int , s io.ReadWriteCloser) error{
	b := sItoI8(a)
	b = prepStoSend(b)
	sendPanel(b, s)
	return nil
}


