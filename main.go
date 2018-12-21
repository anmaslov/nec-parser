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

const MAX_PARSED_CICLES int = 100

//Получение данных со станции
func stantionListener(phone Phones, p DataProducer)  {
	for {
		addr := strings.Join([]string{phone.Ip, phone.Port}, ":")
		stDesc := string(phone.Id)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal("dial error on addr:", addr, err)
			return
		}
		defer conn.Close()
		kpi := fill()
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
			if err != nil{
				log.Println(addr, err)
			}

			log.Println("trying to parse data", stDesc)
			res := smdr.CDR{}
			err = res.Parser(buff[:rd])
			if err != nil {
				log.Println(err, stDesc)
				kpi.stepUp()
			} else {
				call := fillParam(&res)
				call.Stantion = string(phone.Id)
				p.OutChan <- call // Отправляем данные в канал
				//Отправляем запрос о то, что все ок
				r4 := smdr.SetRequest(smdr.ClientResponse(res.Sequence))
				if wr, err := conn.Write([]byte(r4)); //Запрос #4
					wr == 0 || err != nil {
					log.Println(err, stDesc)
				} else {
					kpi.stepDown() //Уменьшаем интервал
					i++ //Увеличиваем счетчик распарсеных данных
				}
			}

			if i >= MAX_PARSED_CICLES {
				log.Println("disconnect from", stDesc, "limit of parsed data", i)
				break
			}

			d := time.Duration(kpi.current * float32(time.Second))
			log.Println("sleep on", d, stDesc)
			time.Sleep(d)
		}//end for

		log.Println("error, when connect or receive data", stDesc, "wait 60seconds")
		time.Sleep(time.Minute) //Ждем 1 минуту, прежде чем выполнить повторное подключение
	}
}

func main() {
	//Загружаем конфигурацию
	cfg.loadConfig()

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

	for data := range p.getOutChan(){ //Ждем данных от канала
		//fmt.Println("get data from chain", data)
		insertCall(&data)
	}

}