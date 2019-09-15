package main

import (
	"github.com/anmaslov/nec-parser/kpi"
	"github.com/anmaslov/nec-parser/store"
	"github.com/anmaslov/smdr"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

const MaxParsedCircles int = 100

// stListener получение данных со станции
func stListener(phone store.Phones, chCall chan<- store.CallInfo, logger *zap.Logger) {
	for {
		addr := strings.Join([]string{phone.Ip, phone.Port}, ":")
		log := logger.With(zap.String("addr", addr))
		log.With(zap.String("st_description", string(phone.Id)))
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Error("dial error on addr", zap.Error(err))
			time.Sleep(time.Minute * 5)
			continue
		}

		k, i := kpi.NewKpi(), 0
		for {
			r1 := smdr.SetRequest(smdr.DataRequest())
			if wr, err := conn.Write(r1); //Запрос #1
			wr == 0 || err != nil {
				logger.Debug("some error")
				log.Error("error in first query", zap.Error(err))
				break
			}

			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Error("connection timeout", zap.Error(err))
				break
			}

			log.Debug("trying to get response from")
			buff := make([]byte, 1024)
			rd, err := conn.Read(buff)
			if err != nil {
				log.Error("error when get response", zap.Error(err))
			}

			log.Debug("trying to parse data")
			res := smdr.CDR{}
			err = res.Parser(buff[:rd])
			if err != nil {
				log.Error("error when parse data", zap.Error(err))
				k.StepUp()
			} else {
				call := store.FillParam(&res)
				call.Stantion = string(phone.Id)
				chCall <- *call // Отправляем данные в канал

				//Отправляем запрос о том, что все ок
				r4 := smdr.SetRequest(smdr.ClientResponse(res.Sequence))
				if wr, err := conn.Write(r4); //Запрос #4
				wr == 0 || err != nil {
					log.Error("error when four request", zap.Error(err))
				} else {
					k.StepDown() //Уменьшаем интервал
					i++          //Увеличиваем счетчик распарсеных данных
				}
			}

			if i >= MaxParsedCircles {
				log.Info("disconnect from", zap.Int("current_i", i))
				break
			}

			d := time.Duration(k.GetCurrent() * float32(time.Second))
			log.Info("sleep on",
				zap.Int64("seconds", int64(d)))

			time.Sleep(d)
		} //end for

		conn.Close()

		log.Error("error, when connect or receive data, wait 60seconds")
		time.Sleep(time.Minute) //Ждем 1 минуту, прежде чем выполнить повторное подключение
	}
}
