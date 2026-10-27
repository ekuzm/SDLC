package UI_04_Budget {
    private import UiKit::*;
    part <budget> 'Категории и лимиты' : Screen {
        doc /* REQ-FUNC-004, REQ-FUNC-011 */
        part 'Навигация' : NavBar {
            part 'Обзор' : Link;
            part 'Операции' : Link;
            part 'Бюджет' : Link;
            part 'Платежи' : Link;
            part 'Прогноз' : Link;
            part 'Выйти' : Link;
        }
        part 'Период [октябрь 2026]' : Select;
        part 'Лимиты расходов' : ListView {
            part 'Продукты' : ListItem {
                part '450,00 из 500,00 BYN' : Label;
                part 'Изменить лимит' : Button;
            }
            part 'Транспорт' : ListItem {
                part '150,00 из 120,00 BYN' : Label;
                part 'Превышение на 30,00 BYN' : Banner;
            }
            part 'Жильё и связь' : ListItem {
                part '400,00 из 450,00 BYN' : Label;
                part 'Изменить лимит' : Button;
            }
        }
        part 'Настройка лимита' : Form {
            part 'Месяц [октябрь 2026]' : Select;
            part 'Категория расходов [Транспорт]' : Select;
            part 'Лимит [120,00 BYN]' : Input;
            part 'Сохранить лимит' : Button;
        }
        part 'Категории' : Row {
            part 'Название [Продукты]' : Input;
            part 'Тип [расход]' : Select;
            part 'Создать' : Button;
            part 'Переименовать' : Button;
            part 'Архивировать' : Button;
        }
        part 'Лимит достигнут / превышен — новые операции разрешены' : Banner;
        part 'Архивная категория остаётся в истории; недоступна для новых записей' : Label;
    }
}
