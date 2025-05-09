![Main page](https://github.com/mayye4ka/notpastebin/raw/master/imgs/admin.png)

### NotPastebin

Notpastebin - cервис для создания заметок, который позволяет делиться заметкой только для чтения или с возможностью редактирования/удаления.

Backend для приложения notpastebin. 

#### Запуск

Бэкенд приложения notpastebin. Для полноценной работы нужен [frontend](https://github.com/mayye4ka/notpastebin-frontend).

* Заполните  `.env.local` и выполните `make run` для запуска локально(нужен docker compose)

* Выполните `make test` / `make cover` для прогона тестов / оценки покрытия

* Выполните `go build .` для сборки докер-образа, `docker run` для запуска контейнера с приложением. Настроить приложение можно через переменные окружения, посмотреть их можно в `.env.local`

#### ToDo

- [ ] Провести ручное тестирование
- [ ] Экспирация заметок по времени с момента создания
- [ ] Экспирация заметок по времени последнего просмотра
- [ ] Кеширование данных в redis
- [ ] Возможность создания бессрочных заметок
- [ ] Возможность создания заметок с очень ссылками / пользовательским текстом в ссылке
- [ ] Регистрация и привязка заметок к профилю автора
