package loggingsystemfromicroservices

import (
	"bufio"
	"io"
	"os"
	"time"
)

type FileHandler struct{
	Handler
	writer *bufio.Writer
}

//NewFileHandler return FileHandler instance 
func NewFileHandler(fp *os.File, flushIntervalMs int) Handler{
	h := FileHandler{
		writer: bufio.NewWriter(fp),
	}

	timer := time.NewTimer(time.Duration(flushIntervalMs) * time.Millisecond)
	go func ()  {
		<-timer.C
		_= h.writer
	}()
	return h
}

func (handler FileHandler) append(msg string) error{
	data := []byte(msg+"\n")
	n, err := handler.writer.Write(data)
	if err == nil && n < len(data){
		err = io.ErrShortWrite
	}
	return err
}

func (handler FileHandler) debug(output string) error{
	return handler.append(output)
}

func (handler FileHandler) info(output string) error{
	return handler.append(output)
}

func (handler FileHandler) warn(output string) error{
	return handler.append(output)
}

func (handler FileHandler) error(output string) error{
	return handler.append(output)
}










