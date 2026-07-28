# trustmefetch

`trustmefetch` превращает вопрос `are you a linux?` в цветную системную сводку в духе fastfetch. Программа получает настоящие данные Mac и оформляет их как выбранный дистрибутив Linux.

## Возможности

- 32 встроенные темы
- 10 шуточных вариантов
- Переливающаяся тема `100% LINUX!!!!!!`
- 22 оформления на основе популярных дистрибутивов
- Полноэкранная настройка с живым предпросмотром
- Живой экран с обновлением загрузки процессора, памяти, диска, батареи и времени работы
- Сбор модели Mac, версии системы, ядра, времени работы, оболочки, терминала, процессора, графики, памяти, диска, экрана и батареи
- Поддержка Apple Silicon и Intel
- Безопасное повторное выполнение установщика
- Удаление интеграции из `.zshrc`

## Установка

После публикации первого выпуска:

```sh
curl -fsSL https://raw.githubusercontent.com/necrasov-ilya/trustmefetch/main/install.sh | sh
```

Откройте новое окно терминала и выполните:

```zsh
are you a linux?
```

Интерфейс настройки запускается командой:

```sh
trustmefetch config
```

Также доступны команды:

```sh
trustmefetch themes
trustmefetch theme arch-btw
trustmefetch preview rgb-linux
trustmefetch live
trustmefetch mode live
trustmefetch random
trustmefetch doctor
```

Обычная команда `trustmefetch` работает как fastfetch: печатает снимок и возвращает приглашение оболочки. Команда `trustmefetch live` остаётся открытой и обновляет показатели до нажатия `q`. Поведение вопроса выбирается через `trustmefetch mode snapshot` или `trustmefetch mode live`. В интерфейсе настройки режим переключается клавишей `m`.

## Локальная разработка

```sh
git clone https://github.com/necrasov-ilya/trustmefetch.git
cd trustmefetch
make test
make install
```

Проект распространяется по лицензии MIT. ASCII-логотипы дистрибутивов взяты из официального каталога fastfetch на условиях MIT с сохранением уведомления об авторских правах. Названия и знаки дистрибутивов принадлежат их правообладателям. Проект не связан с Apple, KDE, fastfetch и разработчиками дистрибутивов Linux.
