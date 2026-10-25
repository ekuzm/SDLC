package UI_05_Payments {
    private import UiKit::*;
    part <payments> 'Плановые платежи' : Screen {
        doc /* REQ-FUNC-012, REQ-FUNC-013, REQ-FUNC-018 */
        part 'Навигация' : NavBar {
            part 'Обзор' : Link;
            part 'Операции' : Link;
            part 'Бюджет' : Link;
            part 'Платежи' : Link;
            part 'Прогноз' : Link;
            part 'Выйти' : Link;
        }
        part 'Регулярный платёж' : Form {
            part 'Название [Связь]' : Input;
            part 'Тип [расход]; категория [Жильё и связь]' : Select;
            part 'Сумма [100,00 BYN]' : Input;
            part 'День месяца [28]' : Input;
            part 'Дата начала; дата окончания [необязательно]' : Input;
            part 'Сохранить правило' : Button;
        }
        part 'Расписание' : ListView {
            part 'Ожидается' : ListItem {
                part '28.10 · Связь · 100,00 BYN' : Label;
                part 'Подтвердить' : Button;
            }
            part 'Исполнен' : ListItem {
                part '20.10 · Жильё и связь · 100,00 BYN' : Label;
                part 'Открыть операцию' : Link;
            }
        }
        part 'Подтверждение' : Form {
            part 'Создать операцию / связать существующую' : Select;
            part 'Фактическая сумма и дата' : Input;
            part 'Сохранить исполнение' : Button;
            part 'Одна операция связана только с одним вхождением' : Label;
        }
        part 'Просроченный платёж' : Row {
            part 'Остаётся в прогнозе до решения' : Banner;
            part 'Подтвердить оплату' : Button;
            part 'Пропустить' : Button;
            part 'Перенести дату' : Button;
        }
        part 'Для короткого месяца — последний день; план сам не меняет остаток' : Label;
    }
}
