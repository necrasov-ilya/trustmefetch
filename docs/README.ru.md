# trustmefetch

`trustmefetch` превращает вопрос `you are a linux?` в цветную системную сводку в духе fastfetch. Программа получает настоящие данные Mac и оформляет их как выбранный дистрибутив Linux.

## Возможности

- 32 встроенные темы
- 10 шуточных вариантов
- Переливающаяся тема `100% LINUX!!!!!!`
- 22 оформления на основе популярных дистрибутивов
- Полноэкранная настройка с живым предпросмотром
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
you are a linux?
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
trustmefetch random
trustmefetch doctor
```

## Локальная разработка

```sh
git clone https://github.com/necrasov-ilya/trustmefetch.git
cd trustmefetch
make test
make install
```

Проект распространяется по лицензии MIT. Названия и знаки дистрибутивов принадлежат их правообладателям. Проект не связан с Apple, KDE, fastfetch и разработчиками дистрибутивов Linux.

