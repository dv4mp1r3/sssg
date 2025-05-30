**Генератор статичных сайтов**

Целью было попрактиковаться в изучении go, и опробовать результаты при переносе своего бложика с wordpress на статику. На текущий момент проект находится в стадии концепта и НЕ МОЖЕТ быть рекомендован в качестве использования.

**Требования**

Linux или MACOS, go 1.18 или выше

**Как работает**

Необходимо скопировать файл config.json.example в config.json и заполнить его в соответствии с тем, как организован контент для генерации.

**Запуск в докере**

Самостоятельная сборка:

```
docker build -t sssg:1.0 .
```
Сборка образа с кастомным пользователем (по умолчанию UID/GID=1000):
```
docker build --build-arg UID=1234 --build-arg GID=1234 -t sssg:1.0 .
```

Пример генерации на основе каталога source.example
```
docker run -v "$(pwd)/source.example":/app/data:rw sssg:1.0 -c /app/data/config.json
```

Ниже пример конфига с комментариями по каждому параметру.
```
{
    // title главной страницы
    "label": "Test homepage",

    // меню сайта
    "menu" : [
        {
            // метка меню (название элемента)
            "label" : "Home",

            // куда ведет ()
            "url" : "/"
        },
        {
            "label": "About",
            "url": "/about"
        }
    ],

    // количество постов на страницу
    "postsPerPage": 3,

    // количество символов (с округлением символов до целых слов) из которых формируется превью
    // не используется, если флаг previewByPageBreak выставлен в true
    "previewLength": 15,

    // флаг получения превью поста по специальному разделителю (значение параметра previewPageBreakString)
    // альтернатива использованию параметра previewLength для указания превью индивидуально для каждого поста
    "previewByPageBreak": false,

    // разделитель, по которому будет формироваться превью поста (от 0 символа до этого разделителя)
    // работает вместе с установленным в true флагом previewByPageBreak
    "previewPageBreakString": "<!--pb-->",

    // вариант пагинации с двумя кнопками (назад/вперед), а не всеми возможными страницами
    "maxTwoPaginationButtons": true,

    // путь, к содержимому ИЗ которого нужно генерировать сайт
    "sourcePath": "./source",

    // путь, куда будет сгенерирован сайт
    "resultPath": "./result",

    // путь к статике относительно параметра sourcePath
    "staticPath": "./css",
    "url": "https://domain.com"
}
```

**Правила обработки шаблонов**

В качестве парсинга шаблонов используется стандартная библиотека, обрабатываются все *.html файлы, путь берется из параметра sourcePath конфига. Шаблон для страницы обязан называться page.html, все остальные файлы опциональны. 
Пример можно посмотреть в каталоге source.example.
