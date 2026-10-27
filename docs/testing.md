# Проверка приложения «Баланс»

Дата выполнения: 27 октября 2026 года. Проверки выполнены автором; результаты рецензирования заказчиком сюда не приписываются.

## Выполненные проверки

| Проверка | Результат |
|---|---|
| go test -race ./... | Пройдено: финансовая логика, HTTP, хранилище, безопасность, конкурентные запросы |
| go vet ./... | Пройдено |
| go build ./cmd/balance | Пройдено на macOS arm64 и Linux amd64, Go 1.26.4 |
| Docker build / запуск | Пройдено: Linux arm64, пользователь balance без root, healthz=ok; compose config корректен |
| Финансовые примеры F01–F17 | Проверяются тестами; 850,00 и 400,00 BYN воспроизводятся |
| Атомарный отказ сохранения | Пройдено: ошибка записи не изменяет состояние в памяти |
| Параллельное подтверждение платежа | Пройдено: 12 запросов, только одна операция |
| Перезапуск и копии | Пройдено: повторное открытие файла, семь копий, восстановление |
| Сквозная проверка браузером | Chromium 151.0.7922.34 и Firefox 153.0; регистрация, бюджет, операции, лимит, платежи, прогнозы, выход |
| Вёрстка | 360, 768 и 1440 px; без горизонтального переполнения; экраны просмотрены |
| Проверка материалов ЛР №3 | tools/check_artifacts.py: 37 требований, связи моделей, девять рисунков DOCX, расчётные примеры |

Тесты находятся в internal/balance/model_test.go, server_test.go и performance_test.go. Проверка защищённых страниц включает попытки обращения к чужим объектам, CSRF, экранирование HTML, срок сеанса и выход. Пять ошибок входа блокируют правильный пароль до истечения 15 минут. Проверка cookie включает Secure, HttpOnly и SameSite=Strict.

Покрытие пакета приложения на итоговой полной проверке: 82,1% операторов. Запуск main отдельно проверяется сборкой, запуском и браузером; покрытие main unit-тестами не заявляется.

## Локальное измерение производительности

```sh
BALANCE_PERF=1 go test ./internal/balance -run TestLocalPerformance -count=1 -v
```

Фактическое окружение: macOS arm64, Go 1.26.4, доступно 10 логических CPU. 10 пользователей, по 10000 операций и 100 активных правил; 1000 запросов к каждому экрану, по 10 одновременных запросов, после прогрева. Время включает обработчик HTTP, копирование состояния и формирование HTML; сеть и браузер в этот замер не входят.

| Экран | p95 | Максимум |
|---|---|---|
| Операции | 38,91 мс | 47,74 мс |
| Аналитика и оба прогноза | 26,62 мс | 32,79 мс |

Дополнительно выполнен тот же тест в Linux arm64-контейнере с ограничениями SRS: 2 CPU, 2 GiB RAM. Проверены cgroup cpu.max=200000 100000 и memory.max=2147483648. p95: операции 178,92 мс (максимум 329,94 мс), аналитика 114,78 мс (максимум 272,48 мс). Пороги 2 секунды и 1 секунда выдержаны. runtime.NumCPU внутри контейнера показывает 10 CPU хоста; фактическая квота ограничена cgroup двумя CPU.

```sh
docker run --rm --cpus=2 --memory=2g \
  --mount type=bind,src="$PWD",dst=/src,readonly -w /src \
  golang:1.26.4-alpine sh -c \
  'BALANCE_PERF=1 go test ./internal/balance -run TestLocalPerformance -count=1 -v'
```

## Воспроизведение

```sh
go test -race -coverprofile=coverage.out ./...
go vet ./...
go build ./cmd/balance
python3 tools/check_artifacts.py
```

Для необязательной браузерной проверки нужен Node.js с playwright и установленным Chromium. Для Firefox задайте BALANCE_BROWSER=firefox. Основному приложению Node.js не требуется. Запускайте тест на отдельном экземпляре с временным файлом данных: он создаёт синтетическую учётную запись и операции, не удаляя чужие данные.

```sh
BALANCE_ADDR=127.0.0.1:18080 BALANCE_DATA=/tmp/balance-smoke/data.json go run ./cmd/balance
# В другом терминале с доступным playwright:
BALANCE_URL=http://127.0.0.1:18080 node tools/smoke.cjs
```

## Трассировка 37 требований

«Проверено» ниже обозначает перечисленные сценарии, а не полную сертификацию всех свойств.

| ID | Подтверждение реализации | Ограничения / дополнительные проверки |
|---|---|---|
| REQ-FUNC-001 | TestPasswordsAndLockout; браузерная регистрация | Пароли 12–128, уникальный email |
| REQ-FUNC-002 | TestCSRFIsolationAndEscaping; вход/выход в браузере | Сеансы аннулируются после перезапуска |
| REQ-FUNC-003 | TestOperationsAndCategories; бюджет через браузер | В реализации начало не раньше 01.01.2000 |
| REQ-FUNC-004 | TestOperationsAndCategories; TestAllScreens | Категории создаются, переименовываются и архивируются |
| REQ-FUNC-005 | TestMoney, TestOperationsAndCategories, TestDeletedOperationFormCannotRecreateRecord | Валидация и постоянная идемпотентность формы |
| REQ-FUNC-006 | TestOperationsAndCategories; браузерное изменение | Проверка неверной правки и дельты остатка |
| REQ-FUNC-007 | TestForecastExamples; браузерное удаление/отмена | Связанное вхождение снова ожидается |
| REQ-FUNC-008 | TestFiltersLimits | Итоги всей выборки, включительные даты, 20 записей |
| REQ-FUNC-009 | TestForecastExamples; обзор в браузере | План отдельно от фактического остатка |
| REQ-FUNC-010 | TestAllScreens; аналитика в браузере | Нулевая база, суммы и периоды; проценты отображаются округлённо |
| REQ-FUNC-011 | TestFiltersLimits; TestPaymentHTTPReplayAndLimits | При равенстве «Лимит достигнут»; расход не блокируется |
| REQ-FUNC-012 | TestCalendarAndRuleEdits; TestOneToOneLinksAndRuleHistory | Февраль 2027/2028, прошлые и перенесённые вхождения |
| REQ-FUNC-013 | TestStoreTransactionsAndBackups; TestOneToOneLinksAndRuleHistory | Связь один к одному, исключение из истории |
| REQ-FUNC-014 | TestForecastExamples, TestForecastRoundingAndIsolation | F01, F06, F07, F16 |
| REQ-FUNC-015 | TestForecastExamples, TestForecastRoundingAndIsolation | F02, F09, F16, F17 |
| REQ-FUNC-016 | TestForecastExamples | F03, F04, F13, F14 |
| REQ-FUNC-017 | TestAllScreens; аналитика в браузере | История, число дней, формулы и причина недоступности |
| REQ-FUNC-018 | TestForecastExamples | F12, F17; пропуск/перенос без фактического списания |
| REQ-INT-001 | tools/smoke.cjs | Chromium и Firefox 360/768/1440; клавиатурные формы — семантический HTML |
| REQ-INT-002 | TestCSRFIsolationAndEscaping; браузерный неверный ввод | Уведомление отправки, пустые состояния, сохранение полей |
| REQ-INT-003 | Браузеры Chromium/Firefox; TLS-конфигурация сервера и TestHTTPSRegistrationAndSession | Сетевой HTTPS-деплой требует отдельной проверки |
| REQ-PERF-001 | TestLocalPerformance в Linux-контейнере 2 CPU / 2 GiB | p95=178,92 мс, 1000 запросов |
| REQ-PERF-002 | То же окружение, 100 правил на пользователя | p95=114,78 мс, 1000 запросов |
| REQ-SEC-001 | TestPasswordsAndLockout; auth.go | PBKDF2, солёные хеши, постоянные счётчики ошибок |
| REQ-SEC-002 | TestSessionExpiryAndCookie; TestCSRFIsolationAndEscaping; TestHTTPSRegistrationAndSession | Secure в HTTPS-режиме; локальный HTTP — исключение для демонстрации |
| REQ-SEC-003 | TestCSRFIsolationAndEscaping; TestForecastRoundingAndIsolation | Владелец при чтении, изменении, удалении и прогнозе |
| REQ-REL-001 | TestStoreTransactionsAndBackups; TestPaymentHTTPReplayAndLimits | Один процесс, локальная файловая система; без распределённых транзакций |
| REQ-REL-002 | TestStoreTransactionsAndBackups | Семь копий/восстановление; норматив 30 минут в SRS-окружении не измерен |
| REQ-AVAIL-001 | TestStoreTransactionsAndBackups; /healthz | Чтение JSON и пробная запись в каталог |
| REQ-OBS-001 | Инспекция ServeHTTP, health, mutate и main | X-Request-ID, ошибки без финансового содержимого |
| REQ-COMP-001 | Инспекция HTML, футера и прогноза | Нет обещаний точности или банковских услуг |
| REQ-INST-001 | Локальная сборка и запуск; Dockerfile, compose.yaml | Команды в README и implementation.md |
| REQ-BUILD-001 | go.mod, Makefile, Dockerfile, GitHub Actions | Только стандартная библиотека; версия коммита через ldflags |
| REQ-MAINT-001 | Этот протокол, модельные и HTTP-тесты | Непройденные проверки перечислены явно |
| REQ-PORT-001 | Сборка Go, локальная демонстрация, Linux-кросс-сборка | Требуется Unix-подобная ОС из-за flock |
| REQ-COST-001 | Локальные формы/ресурсы | Внешние API отсутствуют; зависимости Go отсутствуют |
| REQ-CM-001 | SRS, этот протокол, customer-review.md, Git | Постоянные 37 ID сохранены |

DOCX ЛР №3 сохранён как отчёт этапа требований, с проектными макетами. Реализация и её фактические проверки описаны здесь отдельно.
