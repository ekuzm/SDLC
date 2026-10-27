package UI_02_Overview {
    private import UiKit::*;
    part <overview> 'Обзор бюджета' : Screen {
        doc /* REQ-FUNC-003, REQ-FUNC-009 */
        part 'Навигация' : NavBar {
            part 'Обзор' : Link;
            part 'Операции' : Link;
            part 'Бюджет' : Link;
            part 'Платежи' : Link;
            part 'Прогноз' : Link;
            part 'Выйти' : Link;
        }
        part 'Начало учёта' : Form {
            part 'Дата начала [01.07.2026]' : Input;
            part 'Начальный остаток [500,00 BYN]' : Input;
            part 'Сохранить настройки' : Button;
            part 'После первой операции дата фиксируется' : Label;
        }
        part 'Октябрь 2026 — по 26 октября' : Row {
            part 'Месяц [октябрь 2026]' : Select;
            part 'Добавить операцию' : Button;
        }
        part 'Показатели' : Row {
            part 'Остаток [1000,00 BYN]' : Label;
            part 'Доходы [1500,00 BYN]' : Label;
            part 'Расходы [1000,00 BYN]' : Label;
        }
        part 'Ближайший платёж' : Row {
            part 'Связь · 28 октября · 100,00 BYN' : Label;
            part 'Открыть плановые платежи' : Link;
        }
        part 'Предупреждения' : Column {
            part 'Отрицательный остаток — предупреждение о дефиците' : Banner;
            part 'Операций пока нет — добавьте первую запись' : Banner;
        }
        part 'Изменение начального остатка требует подтверждения' : Label;
    }
}
