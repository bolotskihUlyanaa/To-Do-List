package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "bhefjndlsroif4309rifn"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH" //ключ подписи
	tokenTTL   = 12 * time.Hour
)

// дополнение стандартных claims
// имеет все поля стандартной claims и дополнительно id пользователя
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

// принимаем репозиторий для работы с базой
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// передаем структуру user еще на слой ниже - в репозиторий
func (s *AuthService) CreateUser(user todolist.User) (int, error) {
	//перед записью пользователя в базу, будем хешировать пароль и только потом передавать в слой репозиториев
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// тк пароли не храним открыто в базе
// здесь используется алгоритм хеширования пароля sha-1
func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	//добавим случайные символы к хешу - соль
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	//для генерации токена надо получить пользователя из базы,
	//если такого пользователя нет, то вернуть ошибку
	//иначе генерируем токен в который записываем id пользователя

	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	//в качестве аргументов стандарный метод для подписи
	//Claims - json объект с набором различных полей
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), //токен будет валидным 12 часов
			IssuedAt:  time.Now().Unix(),               //время когда токен был сгенерирован
		},
		user.Id,
	})
	//возвращает подписанный токен
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	//ParseWithClaims возвращает объект токена в котором есть поле Claims типа интерфейс
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) { //внутренняя функция возвращает ключ пользователя или ошибку
		//проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	//в token есть поле claims типа интерфейс, приводим к нашему объекту и проверяем удачно ли
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	//если удачно распарсили
	return claims.UserId, nil
}
