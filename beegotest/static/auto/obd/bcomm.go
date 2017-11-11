// 实现通信层协议
package obd

import (
	"fmt"
	"time"
)

type BComm struct {
	rAck     chan int
	recvdata []byte
	mf       []byte
}

func NewBComm() *BComm {
	return &BComm{
		rAck: make(chan int),
	}
}

// 从原始数据生成AA帧或AB、AC帧
func (bcomm *BComm) Encode(data []byte) ([][]byte, error) {
	if (len(data) == 0) || (len(data) > 0x800000) {
		return nil, fmt.Errorf("Data size is not legal")
	}

	frames := make([][]byte, 0, 1)
	if len(data) <= 1020 {
		frame := make([]byte, 3, len(data)+4)
		frame[0] = 0xAA
		frame[1] = byte(len(data) >> 8)
		frame[2] = byte(len(data) & 0xFF)
		frame = append(frame, data...)
		frame = append(frame, calculateCRC(frame))
		frames = append(frames, frame)
	} else {
		j := len(data) / 1016
		k := len(data) % 1016

		abFrame := make([]byte, 8)
		abFrame[0] = 0xAB
		abFrame[1] = 0x00
		abFrame[2] = 0x04
		abFrame[3] = byte(len(data) >> 24)
		abFrame[4] = byte(len(data) >> 16)
		abFrame[5] = byte(len(data) >> 8)
		abFrame[6] = byte(len(data))
		abFrame[7] = calculateCRC(abFrame[:7])
		frames = append(frames, abFrame)

		for i := 0; i < j; i++ {
			acFrame := make([]byte, 7, 1024)
			acFrame[0] = 0xAC
			acFrame[1] = 0x03
			acFrame[2] = 0xFC

			sent := i * 1016
			acFrame[3] = byte(sent >> 24)
			acFrame[4] = byte(sent >> 16)
			acFrame[5] = byte(sent >> 8)
			acFrame[6] = byte(sent)

			acFrame = append(acFrame, data[i*1016:(i+1)*1016]...)
			acFrame = append(acFrame, calculateCRC(acFrame))

			frames = append(frames, acFrame)
		}

		if k != 0 {
			acFrame := make([]byte, 7, 8+k)
			acFrame[0] = 0xAC
			acFrame[1] = byte((4 + k) >> 8)
			acFrame[2] = byte(4 + k)

			sent := j * 1016
			acFrame[3] = byte(sent >> 24)
			acFrame[4] = byte(sent >> 16)
			acFrame[5] = byte(sent >> 8)
			acFrame[6] = byte(sent)

			acFrame = append(acFrame, data[j*1016:]...)
			acFrame = append(acFrame, calculateCRC(acFrame))

			frames = append(frames, acFrame)
		}
	}

	return frames, nil
}

// 从AA帧或AB、AC帧生成原始数据
func (bcomm *BComm) Decode(data []byte) ([][]byte, error) {
	if (len(data) == 0) || (len(data) > 0x800000) {
		return nil, fmt.Errorf("data size is not legal")
	}
	bcomm.recvdata = append(bcomm.recvdata, data...)

	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	ds := make([][]byte, 0, 1)
	for {
		if len(bcomm.recvdata) <= 4 {
			break
		}

		if (bcomm.recvdata[0] != 0xAA) && (bcomm.recvdata[0] != 0xAB) && (bcomm.recvdata[0] != 0xAC) && (bcomm.recvdata[0] != 0xAD) {
			fmt.Printf("data error\n")
			bcomm.recvdata = bcomm.recvdata[1:]
			continue
		}

		length := (int(bcomm.recvdata[1]) & 0x000000FF)
		length = length << 8
		length |= (int(bcomm.recvdata[2]) & 0x000000FF)
		if len(bcomm.recvdata) < length+4 {
			break
		}

		if calculateCRC(bcomm.recvdata[:length+3]) != bcomm.recvdata[length+3] {
			fmt.Printf("CRC error\n")
			bcomm.recvdata = bcomm.recvdata[1:]
			continue
		}

		frame := bcomm.recvdata[:length+4]
		bcomm.recvdata = bcomm.recvdata[length+4:]

		if frame[0] == 0xAA {
			ds = append(ds, frame[3:length+3])

			bcomm.sendACK(length)
		} else if frame[0] == 0xAB {
			bcomm.mf = nil
			if frame[3] != 0 {
				fmt.Printf("mutiple frame size is too large\n")
				continue
			}

			mfLen := (int(frame[4]) & 0x000000FF)
			mfLen = mfLen << 8
			mfLen |= (int(frame[5]) & 0x000000FF)
			mfLen = mfLen << 8
			mfLen |= (int(frame[6]) & 0x000000FF)
			bcomm.mf = make([]byte, 0, mfLen)

			bcomm.sendACK(0)
		} else if frame[0] == 0xAC {
			if bcomm.mf == nil {
				fmt.Printf("no header found\n")
				continue
			}

			sent := (int(frame[4]) & 0x000000FF)
			sent = sent << 8
			sent |= (int(frame[5]) & 0x000000FF)
			sent = sent << 8
			sent |= (int(frame[6]) & 0x000000FF)

			if sent != len(bcomm.mf) {
				fmt.Printf("there is an error in data transfer\n")
				bcomm.mf = nil
				continue
			}

			if len(bcomm.mf)+length-4 > cap(bcomm.mf) {
				fmt.Printf("there is an error in data transfer\n")
				bcomm.mf = nil
				continue
			}

			bcomm.mf = append(bcomm.mf, frame[7:length+3]...)
			if len(bcomm.mf) == cap(bcomm.mf) {
				ds = append(ds, bcomm.mf)
				bcomm.mf = nil
			}

			bcomm.sendACK(sent + length - 4)
		} else if frame[0] == 0xAD {
			if frame[3] != 0 {
				fmt.Printf("ack frame size is too large\n")
				continue
			}

			fLen := (int(frame[4]) & 0x000000FF)
			fLen = fLen << 8
			fLen |= (int(frame[5]) & 0x000000FF)
			fLen = fLen << 8
			fLen |= (int(frame[6]) & 0x000000FF)

			timer.Reset(100 * time.Millisecond)
			select {
			case <-timer.C:
				fmt.Printf("discard ack\n")
			case bcomm.rAck <- fLen:
				//			default:
				//				fmt.Printf("discard ack\n")
			}
		}
	}
	return ds, nil
}

func (bcomm *BComm) sendACK(length int) []byte {
	frame := make([]byte, 8)
	frame[0] = 0xAD
	frame[1] = 0x00
	frame[2] = 0x04
	frame[3] = byte(length >> 24)
	frame[4] = byte(length >> 16)
	frame[5] = byte(length >> 8)
	frame[6] = byte(length)
	frame[7] = calculateCRC(frame[:7])

	// 发送ACK

	return frame
}

func (bcomm *BComm) recvACK(length int) error {
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout")
		case l, _ := <-bcomm.rAck:
			//			fmt.Printf("l: %d\n", l)
			if l == length {
				return nil
			}
		}
	}
}

func calculateCRC(data []byte) byte {
	var crc byte = 0x00
	for i := 0; i < len(data); i++ {
		crc += data[i]
	}
	return crc
}

func bsTostr(data []byte) string {
	str := ""
	for i := 0; i < len(data); i++ {
		str += fmt.Sprintf("%02X", data[i])
	}
	return str
}
