package main

import (
	"fmt"
	"github.com/anmaslov/nec-parser/config"
	"github.com/anmaslov/nec-parser/kpi"
	"github.com/anmaslov/smdr"
	"github.com/jessevdk/go-flags"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var cfg config.Config

type DataProducer struct {
	OutChan chan CallInfo
}

func (p *DataProducer) getOutChan() <-chan CallInfo {
	return p.OutChan
}

const MAX_PARSED_CICLES int = 100

//Получение данных со станции
func stantionListener(phone Phones, p DataProducer) {
	for {
		addr := strings.Join([]string{phone.Ip, phone.Port}, ":")
		stDesc := string(phone.Id)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println("dial error on addr:", addr, err)
			time.Sleep(time.Minute * 5)
			continue
		}
		k := kpi.NewKpi()
		//Основной цикл для получения данных
		i := 0 //Счетчик распарсенных данных
		for {
			r1 := smdr.SetRequest(smdr.DataRequest())
			if wr, err := conn.Write([]byte(r1)); //Запрос #1
			wr == 0 || err != nil {
				log.Println(addr, err)
				break
			}

			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Println(addr, err)
				break
			}

			log.Println("trying to get response from", stDesc)
			buff := make([]byte, 1024)
			rd, err := conn.Read(buff)
			if err != nil {
				log.Println(addr, err)
			}

			log.Println("trying to parse data", stDesc)
			res := smdr.CDR{}
			err = res.Parser(buff[:rd])
			if err != nil {
				log.Println(err, stDesc)
				k.StepUp()
			} else {
				call := fillParam(&res)
				call.Stantion = string(phone.Id)
				p.OutChan <- call // Отправляем данные в канал
				//Отправляем запрос о том, что все ок
				r4 := smdr.SetRequest(smdr.ClientResponse(res.Sequence))
				if wr, err := conn.Write([]byte(r4)); //Запрос #4
				wr == 0 || err != nil {
					log.Println(err, stDesc)
				} else {
					k.StepDown() //Уменьшаем интервал
					i++          //Увеличиваем счетчик распарсеных данных
				}
			}

			if i >= MAX_PARSED_CICLES {
				log.Println("disconnect from", stDesc, "limit of parsed data", i)
				break
			}

			d := time.Duration(k.GetCurrent() * float32(time.Second))
			log.Println("sleep on", d, stDesc)

			time.Sleep(d)
		} //end for

		conn.Close()

		log.Println("error, when connect or receive data", stDesc, "wait 60seconds")
		time.Sleep(time.Minute) //Ждем 1 минуту, прежде чем выполнить повторное подключение
	}
}

func init() {
	parser := flags.NewParser(&cfg, flags.Default)
	parser.SubcommandsOptional = true
	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error init: %s.\nFor help use -h\n", err)
		os.Exit(1)
	}
}

func main() {
	//todo переделать логгер
	f, err := os.OpenFile("/var/log/phone.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("file no exist", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	session := initialiseMongo()
	mongoStore.session = session
	defer session.Close()

	//Создаем канал
	p := DataProducer{
		OutChan: make(chan CallInfo),
	}

	phones, err := getPhones()
	if err != nil {
		log.Fatal(err)
	}
	for _, phone := range phones {
		go stantionListener(phone, p) //Запускаем столько гоурутин, сколько телефонов
	}

	for data := range p.getOutChan() { //Ждем данных от канала
		//fmt.Println("get data from chain", data)
		err := insertCall(&data)
		if err != nil {
			log.Fatal(err) //Падаем, т.к. запись в базу - критична
		} else {
			log.Println("write to DB success, date end call:", data.Cvt.DateEnd.String())
		}
	}

}
