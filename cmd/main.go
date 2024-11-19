package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/handler"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/repository"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	//загрузить переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable: %s", err.Error())
	}

	//инициализируем базу
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("error failed to initialize db: %s", err.Error())
	}

	//объявляем все зависимости в нужном порядке (внутренний круг ничего не знает о внешнем)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)    //сервис зависит от репозитория
	handlers := handler.NewHandler(services) //обработчик зависит от сервиса

	srv := new(todolist.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil { //HTTP сервер по умолчанию запускается на 80 порту
			logrus.Fatalf("error occured while running http server: %s", err.Error()) //вывести ошибку на экран и выйти из приложения
		}
	}()
	logrus.Print("TodoApp started")
	quit := make(chan os.Signal, 1) //канал типа os.Signal
	//запись в канал будет когда процесс в котором выполняется приложение получит от системы сигнал типа syscall.SIGTERM или syscall.SIGINT
	signal.Notify(quit, syscall.SIGTERM)
	//чтение из канала, блокирует выполнение главной горутины main
	<-quit
	//вывод информации о том, что приложение заканчивает свое выполнение
	logrus.Print("TodoApp shutting down")
	//остановка сервера
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shuting down: %s", err.Error())
	}
	//закрытие всех соединений с бд
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

// для инициализации конфигурационного файла
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig() //записывает во внутренний конфиг viper
}
