[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# Проверяет доступность http-ссылок. Ссылки берет из текстового файла. Ошибки отправляет в телеграм (Russian)

Имя файла со ссылками можно передать параметром
```shell
checker --file <filename.txt>
```
По имолчанию `links.txt`

В .env файле нужно указать настройки для Telegram (токен и Chat ID)

Также можно изменить кол-во ворекров для параллельной работы (NUM_WORKERS)
