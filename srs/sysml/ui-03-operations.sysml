package UI_03_Operations {
    private import UiKit::*;
    part <operations> 'Операции' : Screen {
        doc /* REQ-FUNC-005, REQ-FUNC-006, REQ-FUNC-007, REQ-FUNC-008 */
        part 'Навигация' : NavBar {
            part 'Обзор' : Link;
            part 'Операции' : Link;
            part 'Бюджет' : Link;
            part 'Платежи' : Link;
            part 'Прогноз' : Link;
            part 'Выйти' : Link;
        }
        part 'Фильтры' : Row {
            part 'С [01.10.2026]' : Input;
            part 'По [26.10.2026]' : Input;
            part 'Тип [все]' : Select;
            part 'Категория [все]' : Select;
            part 'Применить' : Button;
        }
        part 'История операций' : ListView {
            part 'Запись 1' : ListItem {
                part '26.10 · Продукты · −45,00 BYN' : Label;
                part 'Изменить' : Button;
                part 'Удалить' : Button;
            }
            part 'Запись 2' : ListItem {
                part '05.10 · Зарплата · +1500,00 BYN' : Label;
                part 'Изменить' : Button;
                part 'Удалить' : Button;
            }
        }
        part '20 записей на странице; итоги — по всей выборке' : Label;
        part 'Ввод и редактирование' : Form {
            part 'Тип [расход]' : Select;
            part 'Сумма [45,00 BYN]' : Input;
            part 'Дата [26.10.2026]' : Input;
            part 'Категория [Продукты]' : Select;
            part 'Комментарий [Покупки на неделю]' : Input;
            part 'Действия' : Row {
                part 'Сохранить' : Button;
                part 'Отменить' : Button;
            }
        }
        part 'Подтверждение удаления' : Row {
            part 'Удалить выбранную операцию?' : Label;
            part 'Удалить окончательно' : Button;
            part 'Отмена' : Button;
        }
        part 'Состояния' : Column {
            part 'Пустая выборка: операций не найдено' : Banner;
            part 'Ошибка суммы: введите значение больше нуля' : Banner;
            part 'Ошибка сохранения: значения остаются в форме' : Banner;
        }
    }
}
