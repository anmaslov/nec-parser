package main

import (
	"github.com/anmaslov/smdr"
	"log"
	"net"
	"strings"
	"time"
)

var cfg = Configuration{}

type DataProducer struct {
	OutChan chan CallInfo
}

func (p *DataProducer) getOutChan() <- chan CallInfo{
	return p.OutChan
}

//Получение данных со станции
func stantionListener(phone Stantion, p DataProducer)  {
	for {
		addr := strings.Join([]string{phone.Ip, phone.Port}, ":")
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal("dial error on addr:", addr, err)
			return
		}
		defer conn.Close()
		kpi := fill()
		//Основной цикл для получения данных
		for {
			r1 := smdr.SetRequest(smdr.DataRequest())
			if wr, err := conn.Write([]byte(r1)); //Запрос #1
				wr == 0 || err != nil {
				log.Println(phone.Name, err)
				break
			}

			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Println(phone.Name, err)
				break
			}

			log.Println("trying to get response from phone stantion", phone.Name)
			buff := make([]byte, 1024)
			rd, err := conn.Read(buff)
			if err != nil{
				log.Println(phone.Name, ": ", err)
			}

			log.Println("trying to parse data", addr)
			res := smdr.CDR{}
			err = res.Parser(buff[:rd])
			if err != nil {
				log.Println(err)
				log.Println(phone.Name, ":", err)
				kpi.stepUp()
			} else {
				call := fillParam(&res)
				call.Stantion = phone.Name
				p.OutChan <- call // Отправляем данные в канал
				kpi.stepDown() //Уменьшаем интервал

				//Отправляем запрос о то, что все ок
				r4 := smdr.SetRequest(smdr.ClientResponse(res.Sequence))
				if wr, err := conn.Write([]byte(r4)); //Запрос #4
					wr == 0 || err != nil {
					log.Println("error on", phone.Name, "error text", err)
				}
			}

			d := time.Duration(kpi.current * float32(time.Second))
			log.Println(phone.Name, "is sleep on", d)
			time.Sleep(d)
		}//end for

		log.Println("error, when connect or receive data from stantion", addr)
		time.Sleep(time.Minute) //Ждем 1 минуту, прежде чем выполнить повторное подключение
	}
}

func main() {
	//Загружаем конфигурацию
	cfg.loadConfig()

	session := initialiseMongo()
	mongoStore.session = session
	defer session.Close()

	go func(){
		startServer()
	}()

	//Создаем канал
	p := DataProducer{
		OutChan: make(chan CallInfo),
	}

	for _, phone := range cfg.Phones {
		go stantionListener(phone, p) //Запускаем столько гоурутин, сколько телефонов
	}

	for data := range p.getOutChan(){ //Ждем данных от канала
		//fmt.Println("get data from chain", data)
		insertCall(&data)
	}

}