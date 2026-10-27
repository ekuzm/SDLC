// Концептуальный контекст; выбор СУБД и веб-фреймворка не фиксируется.
package Context {
    item def 'Действия';
    item def 'Страницы и прогнозы';
    item def 'Запросы';
    item def 'Результаты и ошибки';
    item def 'Данные бюджета';
    item def 'Записи и планы';
    part <ctx> 'Контекст системы Баланс' {
        part <user> 'Пользователь' {
            out item actions : 'Действия';
            in item pages : 'Страницы и прогнозы';
        }
        part <system> 'Баланс — веб-приложение' {
            part <web> 'Интерфейс в браузере' {
                in item actions : 'Действия';
                out item pages : 'Страницы и прогнозы';
                out item requests : 'Запросы';
                in item results : 'Результаты и ошибки';
            }
            part <server> 'Сервер Go' {
                in item requests : 'Запросы';
                out item results : 'Результаты и ошибки';
                out item changes : 'Данные бюджета';
                in item records : 'Записи и планы';
            }
            part <storage> 'Постоянное хранилище' {
                in item changes : 'Данные бюджета';
                out item records : 'Записи и планы';
            }
            flow of 'Запросы' from web.requests to server.requests;
            flow of 'Результаты и ошибки' from server.results to web.results;
            flow of 'Данные бюджета' from server.changes to storage.changes;
            flow of 'Записи и планы' from storage.records to server.records;
        }
        flow of 'Действия' from user.actions to system.web.actions;
        flow of 'Страницы и прогнозы' from system.web.pages to user.pages;
    }
}
